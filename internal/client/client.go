package fractal_cloud

import (
	"net/http"
	"time"
)

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
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
func NewClient(host *string, serviceAccountId *string, serviceAccountSecret *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		// Default Fractal Cloud URL
		HostURL: *host,
	}

	c.Auth = AuthStruct{
		ServiceAccountId:     *serviceAccountId,
		ServiceAccountSecret: *serviceAccountSecret,
	}

	ar, err := c.SignIn()
	if err != nil {
		return nil, err
	}

	c.Token = ar.Token

	return &c, nil
}

func (c *Client) doRequest(req *http.Request, authToken *string) ([]byte, error) {
	return []byte("{}"), nil
}
