package goeco

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"net/http"
	"strconv"
	"time"
)

func New(httpClient *http.Client, opts ...Opt) (*Client, error) {
	opts = append(opts, withRetry(
		RetryConfig{
			RetryMax:     defaultRetryMax,
			RetryWaitMin: ptr(float64(defaultRetryWaitMin)),
			RetryWaitMax: ptr(float64(defaultRetryWaitMax)),
		},
	))
	c := newClient(httpClient)
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	if c.retryConfig.RetryMax > 0 {
		retryableClient := retryablehttp.NewClient()
		retryableClient.RetryMax = c.retryConfig.RetryMax

		if c.retryConfig.RetryWaitMin != nil {
			retryableClient.RetryWaitMin = time.Duration(*c.retryConfig.RetryWaitMin * float64(time.Second))
		}
		if c.retryConfig.RetryWaitMax != nil {
			retryableClient.RetryWaitMax = time.Duration(*c.retryConfig.RetryWaitMax * float64(time.Second))
		}

		// By default, this is nil and does not log.
		retryableClient.Logger = c.retryConfig.Logger

		// if timeout is set, it is maintained before overwriting client with StandardClient()
		retryableClient.HTTPClient.Timeout = c.httpClient.Timeout

		retryableClient.ErrorHandler = func(resp *http.Response, err error, numTries int) (*http.Response, error) {
			if resp != nil {
				resp.Header.Add(internalHeaderRetryAttempts, strconv.Itoa(numTries))

				return resp, err
			}
			return resp, err
		}
		c.httpClient = retryableClient.StandardClient()
	}
	return c, nil
}
func newClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	c := &Client{
		httpClient: httpClient,
		userAgent:  userAgent,
		headers:    make(map[string]string),
	}
	c.Clusters = &ClusterService{client: c}
	c.Pools = &PoolService{client: c}
	c.Users = &UserService{client: c}
	c.Roles = &RoleService{client: c}

	return c
}
func ptr[T any](v T) *T {
	return &v
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(_ context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	var req *http.Request
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		req, err = http.NewRequest(method, u.String(), nil)
		if err != nil {
			return nil, err
		}
	default:
		buf := new(bytes.Buffer)
		if body != nil {
			err = json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
		}

		req, err = http.NewRequest(method, u.String(), buf)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", mediaType)
	}

	for k, v := range c.headers {
		req.Header.Add(k, v)
	}
	req.Header.Set("Accept", mediaType)
	req.Header.Set("User-Agent", c.userAgent)
	return req, nil
}

// Do send an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return &Response{
			Response: &http.Response{
				Status:     http.StatusText(http.StatusInternalServerError),
				StatusCode: http.StatusInternalServerError,
			},
		}, err
	}

	defer func() {
		// Ensure the response body is fully read and closed
		// before we reconnect, so that we reuse the same TCPConnection.
		// Close the previous response's body. But read at least some of
		// the body so if it's small the underlying TCP connection will be
		// re-used. No need to check for errors: if it fails, the Transport
		// won't reuse it anyway.
		const maxBodySlurpSize = 2 << 10
		if resp.ContentLength == -1 || resp.ContentLength <= maxBodySlurpSize {
			_, _ = io.CopyN(io.Discard, resp.Body, maxBodySlurpSize)
		}

		if rErr := resp.Body.Close(); err == nil {
			err = rErr
		}
	}()

	response := &Response{Response: resp}

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if resp.StatusCode != http.StatusNoContent && v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
		if err != nil {
			return &Response{
				Response: &http.Response{
					Status:     http.StatusText(http.StatusInternalServerError),
					StatusCode: http.StatusInternalServerError,
				},
			}, err
		}
	}

	return response, err
}

// CheckResponse checks the API response for errors, and returns them if present. A response is considered an
// error if it has a status code outside the 200 range. API error responses are expected to have either no response
// body, or a JSON response body that maps to ResponseError. Any other response body will be silently ignored.
// If the API error response does not include the request ID in its body, the one from its header will be used.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ResponseError{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.Message = string(data)
		}
	}

	attempts, strconvErr := strconv.Atoi(r.Header.Get(internalHeaderRetryAttempts))
	if strconvErr == nil {
		errorResponse.Attempts = attempts
	}

	return errorResponse
}
