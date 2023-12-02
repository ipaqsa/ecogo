package goeco

import (
	"fmt"
	"net/url"
	"strings"
)

func SetURL(bu string) Opt {
	return func(c *Client) error {
		u, err := url.Parse(bu)
		if err != nil {
			return err
		}

		c.baseURL = u

		return nil
	}
}

func SetProjectID(projectID string) Opt {
	return func(c *Client) error {
		c.projectID = projectID
		return nil
	}
}
func SetRegionID(regionID string) Opt {
	return func(c *Client) error {
		c.regionID = regionID
		return nil
	}
}
func SetAPIKey(apiKey string) Opt {
	return func(c *Client) error {
		tokenPartsCount := 2
		parts := strings.SplitN(apiKey, " ", tokenPartsCount)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "apikey" {
			apiKey = parts[1]
		}
		c.apiKey = apiKey
		c.headers["Authorization"] = fmt.Sprintf("apikey %s", c.apiKey)

		return nil
	}
}

func SetUserAgent(ua string) Opt {
	return func(c *Client) error {
		c.userAgent = fmt.Sprintf("%s %s", ua, c.userAgent)
		return nil
	}
}

func SetRequestHeaders(headers map[string]string) Opt {
	return func(c *Client) error {
		for k, v := range headers {
			c.headers[k] = v
		}
		return nil
	}
}

func withRetry(retryConfig RetryConfig) Opt {
	return func(c *Client) error {
		c.retryConfig.RetryMax = retryConfig.RetryMax
		c.retryConfig.RetryWaitMax = retryConfig.RetryWaitMax
		c.retryConfig.RetryWaitMin = retryConfig.RetryWaitMin
		c.retryConfig.Logger = retryConfig.Logger
		return nil
	}
}
