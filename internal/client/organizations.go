package fractalCloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// GetOrganization - Returns specific resource group
func (c *Client) GetOrganization(organizationId string) (*Organization, error) {
	path := fmt.Sprintf("%s/organizations/%s", c.HostURL, organizationId)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	c.logDebug("Calling GET " + path)

	resCode, body, err := c.doRequest(req, []int{200, 404})
	if err != nil {
		return nil, err
	}

	c.logDebug("Response code: " + strconv.Itoa(resCode))

	if resCode == 404 {
		return nil, nil
	}

	c.logDebug(string(body))

	organization := Organization{}
	err = json.Unmarshal(body, &organization)
	if err != nil {
		return nil, err
	}

	return &organization, nil
}
