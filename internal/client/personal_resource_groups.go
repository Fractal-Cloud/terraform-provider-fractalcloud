package fractalCloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetPersonalResourceGroup - Returns specific resource group
func (c *Client) GetPersonalResourceGroup(resourceGroupID ResourceGroupId) (*PersonalResourceGroup, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/accounts/me/resourcegroups/%s",
		c.HostURL, resourceGroupID.ShortName), nil)
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

	resourceGroup := PersonalResourceGroup{}
	err = json.Unmarshal(body, &resourceGroup)
	if err != nil {
		return nil, err
	}

	return &resourceGroup, nil
}

type UpsertPersonalResourceGroupRequestBody struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func (c *Client) UpsertPersonalResourceGroup(resourceGroup PersonalResourceGroup) error {
	resourceGroupID := resourceGroup.ID

	requestBody := UpsertPersonalResourceGroupRequestBody{
		DisplayName: resourceGroup.DisplayName,
		Description: resourceGroup.Description,
	}
	rb, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/accounts/me/resourcegroups/%s",
		c.HostURL, resourceGroupID.ShortName), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, _, err = c.doRequest(req, []int{200})
	return err
}

func (c *Client) DeletePersonalResourceGroup(resourceGroupID ResourceGroupId) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/accounts/me/resourcegroups/%s",
		c.HostURL, resourceGroupID.ShortName), nil)
	if err != nil {
		return err
	}

	_, _, err = c.doRequest(req, []int{200, 404})
	if err != nil {
		return err
	}

	return nil
}
