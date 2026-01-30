package fractalCloud

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetOrganization - Returns specific resource group
func (c *Client) GetOrganization(organizationID string) (*Organization, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/organizations/%s", c.HostURL, organizationID), nil)
	if err != nil {
		return nil, err
	}

	resCode, body, err := c.doRequest(req, []int{200, 404})
	if err != nil {
		return nil, err
	}

	if resCode == 404 {
		return nil, nil
	}

	organization := Organization{}
	err = json.Unmarshal(body, &organization)
	if err != nil {
		return nil, err
	}

	return &organization, nil
}
