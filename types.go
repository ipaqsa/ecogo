package ecogo

import (
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	userAgent  string
	headers    map[string]string

	regionID  string
	projectID string
	apiKey    string

	Clusters ClusterServiceI
	Pools    PoolServiceI
	Users    UserServiceI

	retryConfig RetryConfig
}

// Opt are options for New.
type Opt func(*Client) error

type RetryConfig struct {
	RetryMax     int
	RetryWaitMin *float64    // Minimum time to wait
	RetryWaitMax *float64    // Maximum time to wait
	Logger       interface{} // Customer logger instance. Must implement either go-retryablehttp.Logger or go-retryablehttp.LeveledLogger
}

type Response struct {
	*http.Response
}

// An ResponseError reports the error caused by an API request.
type ResponseError struct {
	// HTTP response that caused this error
	Response *http.Response

	// Error message
	Message string `json:"message"`

	// Attempts is the number of times the request was attempted when retries are enabled.
	Attempts int
}

func (r *ResponseError) Error() string {
	var attempted string
	if r.Attempts > 0 {
		attempted = fmt.Sprintf("; giving up after %d attempt(s)", r.Attempts)
	}

	return fmt.Sprintf("%v %v: %d %v%s",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message, attempted)
}
