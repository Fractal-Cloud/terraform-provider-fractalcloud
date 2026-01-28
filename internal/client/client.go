package fractalCloud

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"time"
)

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Auth       AuthStruct
}

// AuthStruct -
type AuthStruct struct {
	ServiceAccountId     string `json:"ServiceAccountId"`
	ServiceAccountSecret string `json:"ServiceAccountSecret"`
}

// AuthResponse -
type AuthResponse struct {
	UserID           int    `json:"user_id"`
	ServiceAccountId string `json:"ServiceAccountId"`
	Token            string `json:"token"`
}

// NewClient -
func NewClient(host *string, serviceAccountId *string, serviceAccountSecret *string) *Client {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		// Default Fractal Cloud URL
		HostURL: *host,
	}

	c.Auth = AuthStruct{
		ServiceAccountId:     *serviceAccountId,
		ServiceAccountSecret: *serviceAccountSecret,
	}

	return &c
}

func (c *Client) doRequest(req *http.Request, acceptedResponseCodes []int) ([]byte, error) {
	req.Header.Add("X-ClientID", c.Auth.ServiceAccountId)
	req.Header.Add("X-ClientSecret", c.Auth.ServiceAccountSecret)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Print(err)
		}
	}(res.Body)

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[ClientId: %s] Response Code: %d. Impossible to read response body", c.Auth.ServiceAccountId, res.StatusCode))
	}

	if slices.Contains(acceptedResponseCodes, res.StatusCode) {
		return bodyBytes, nil
	}

	return nil, errors.New(fmt.Sprintf("[ClientId: %s] received an unexpected response code: %d. Body: %s", c.Auth.ServiceAccountId, res.StatusCode, string(bodyBytes)))
}
