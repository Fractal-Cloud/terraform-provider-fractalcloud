package fractalCloud

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"
)

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Auth       AuthStruct
	Logger     *ClientLogger
}

type ClientLogger struct {
	Debug       func(string)
	Information func(string)
	Warning     func(string)
	Error       func(string)
}

// AuthStruct -
type AuthStruct struct {
	ServiceAccountId     string `json:"ServiceAccountId"`
	ServiceAccountSecret string `json:"ServiceAccountSecret"`
}

// NewClient -
func NewClient(logger *ClientLogger, host *string, serviceAccountId *string, serviceAccountSecret *string) *Client {
	c := Client{
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		HostURL:    *host,
		Logger:     logger,
	}

	c.Auth = AuthStruct{
		ServiceAccountId:     *serviceAccountId,
		ServiceAccountSecret: *serviceAccountSecret,
	}

	return &c
}

func (c *Client) doRequest(ctx context.Context, req *http.Request, acceptedResponseCodes []int) (int, []byte, error) {
	req = req.WithContext(ctx)
	req.Header.Set("X-ClientID", c.Auth.ServiceAccountId)
	req.Header.Set("X-ClientSecret", c.Auth.ServiceAccountSecret)
	req.Header.Set("Content-Type", "application/json")

	c.logDebug(fmt.Sprintf("%s %s", req.Method, req.URL.String()))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("HTTP request failed (%s %s): %w", req.Method, req.URL.Path, err)
	}

	defer func(Body io.ReadCloser) {
		if closeErr := Body.Close(); closeErr != nil {
			c.logWarning(fmt.Sprintf("failed to close response body: %s", closeErr))
		}
	}(res.Body)

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, nil, fmt.Errorf(
			"failed to read response body (%s %s, status %d): %w",
			req.Method, req.URL.Path, res.StatusCode, err,
		)
	}

	c.logDebug(fmt.Sprintf("%s %s -> %d (%d bytes)", req.Method, req.URL.Path, res.StatusCode, len(bodyBytes)))

	if slices.Contains(acceptedResponseCodes, res.StatusCode) {
		return res.StatusCode, bodyBytes, nil
	}

	return res.StatusCode, bodyBytes, fmt.Errorf(
		"unexpected status %d from %s %s: %s",
		res.StatusCode, req.Method, req.URL.Path, truncateBody(bodyBytes, 500),
	)
}

// truncateBody returns the body as a string, truncated to maxLen for error messages.
func truncateBody(body []byte, maxLen int) string {
	if len(body) <= maxLen {
		return string(body)
	}
	return string(body[:maxLen]) + "... (truncated)"
}
