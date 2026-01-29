package fractalCloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetResourceGroup - Returns specific resource group
func (c *Client) GetResourceGroup(resourceGroupID ResourceGroupId) (*ResourceGroup, error) {
	relativePath := ""
	if resourceGroupID.Type == "Personal" {
		relativePath = fmt.Sprintf("accounts/me/resourcegroups/%s", resourceGroupID.ShortName)
	} else {
		relativePath = fmt.Sprintf("organizations/%s/resourcegroups/%s", resourceGroupID.OwnerId, resourceGroupID.ShortName)
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.HostURL, relativePath), nil)
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

	resourceGroup := ResourceGroup{}
	err = json.Unmarshal(body, &resourceGroup)
	if err != nil {
		return nil, err
	}

	return &resourceGroup, nil
}

type UpsertResourceGroupRequestBody struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

func (c *Client) UpsertResourceGroup(resourceGroup ResourceGroup) error {
	relativePath := ""
	resourceGroupID := resourceGroup.ID
	if resourceGroupID.Type == "Personal" {
		relativePath = fmt.Sprintf("accounts/me/resourcegroups/%s", resourceGroupID.ShortName)
	} else {
		relativePath = fmt.Sprintf("organizations/%s/resourcegroups/%s", resourceGroupID.OwnerId, resourceGroupID.ShortName)
	}

	requestBody := UpsertResourceGroupRequestBody{
		DisplayName: resourceGroup.DisplayName,
		Description: resourceGroup.Description,
	}
	rb, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.HostURL, relativePath), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, _, err = c.doRequest(req, []int{200})
	return err
}

func (c *Client) DeleteResourceGroup(resourceGroupID ResourceGroupId) error {
	relativePath := ""
	if resourceGroupID.Type == "Personal" {
		relativePath = fmt.Sprintf("accounts/me/resourcegroups/%s", resourceGroupID.ShortName)
	} else {
		relativePath = fmt.Sprintf("organizations/%s/resourcegroups/%s", resourceGroupID.OwnerId, resourceGroupID.ShortName)
	}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", c.HostURL, relativePath), nil)
	if err != nil {
		return err
	}

	_, _, err = c.doRequest(req, []int{200, 404})
	if err != nil {
		return err
	}

	return nil
}
