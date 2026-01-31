package fractalCloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetOrganizationalResourceGroup - Returns specific organizational resource group
func (c *Client) GetOrganizationalResourceGroup(resourceGroupID ResourceGroupId) (*OrganizationalResourceGroup, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/organizations/%s/resourcegroups/%s",
		c.HostURL, resourceGroupID.OwnerID, resourceGroupID.ShortName), nil)
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

	c.logDebug(string(body))

	resourceGroup := OrganizationalResourceGroup{}
	err = json.Unmarshal(body, &resourceGroup)
	if err != nil {
		return nil, err
	}

	return &resourceGroup, nil
}

type UpsertOrganizationalResourceGroupRequestBody struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	IconUrl     string `json:"iconUrl"`
}

func (c *Client) UpsertOrganizationalResourceGroup(resourceGroup OrganizationalResourceGroup) error {
	resourceGroupID := resourceGroup.ID

	requestBody := UpsertOrganizationalResourceGroupRequestBody{
		DisplayName: resourceGroup.DisplayName,
		Description: resourceGroup.Description,
	}
	rb, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/organizations/%s/resourcegroups/%s",
		c.HostURL, resourceGroupID.OwnerID, resourceGroupID.ShortName), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, _, err = c.doRequest(req, []int{200})
	return err
}

func (c *Client) DeleteOrganizationalResourceGroup(resourceGroupID ResourceGroupId) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/organizations/%s/resourcegroups/%s",
		c.HostURL, resourceGroupID.OwnerID, resourceGroupID.ShortName), nil)
	if err != nil {
		return err
	}

	_, _, err = c.doRequest(req, []int{200, 404})
	if err != nil {
		return err
	}

	return nil
}
