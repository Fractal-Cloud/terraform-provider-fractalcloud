package fractalCloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetPersonalResourceGroup - Returns specific personal resource group (bounded context).
func (c *Client) GetPersonalResourceGroup(ctx context.Context, resourceGroupId ResourceGroupId) (*PersonalResourceGroup, error) {
	path := fmt.Sprintf("%s/accounts/me/resourcegroups/%s",
		c.HostURL, resourceGroupId.ShortName)

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("building GET request for personal bounded context %q: %w", resourceGroupId.ShortName, err)
	}

	resCode, body, err := c.doRequest(ctx, req, []int{200, 404})
	if err != nil {
		return nil, fmt.Errorf("fetching personal bounded context %q: %w", resourceGroupId.ShortName, err)
	}

	if resCode == 404 {
		c.logDebug(fmt.Sprintf("personal bounded context %q not found", resourceGroupId.ShortName))
		return nil, nil
	}

	resourceGroup := PersonalResourceGroup{}
	if err := json.Unmarshal(body, &resourceGroup); err != nil {
		return nil, fmt.Errorf("decoding personal bounded context %q response: %w", resourceGroupId.ShortName, err)
	}

	return &resourceGroup, nil
}

type UpsertPersonalResourceGroupRequestBody struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func (c *Client) UpsertPersonalResourceGroup(ctx context.Context, resourceGroup PersonalResourceGroup) error {
	resourceGroupId := resourceGroup.Id

	requestBody := UpsertPersonalResourceGroupRequestBody{
		DisplayName: resourceGroup.DisplayName,
		Description: resourceGroup.Description,
	}
	rb, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("encoding personal bounded context %q request: %w", resourceGroupId.ShortName, err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/accounts/me/resourcegroups/%s",
		c.HostURL, resourceGroupId.ShortName), strings.NewReader(string(rb)))
	if err != nil {
		return fmt.Errorf("building POST request for personal bounded context %q: %w", resourceGroupId.ShortName, err)
	}

	_, _, err = c.doRequest(ctx, req, []int{200, 202})
	if err != nil {
		return fmt.Errorf("upserting personal bounded context %q: %w", resourceGroupId.ShortName, err)
	}
	return nil
}

func (c *Client) DeletePersonalResourceGroup(ctx context.Context, resourceGroupId ResourceGroupId) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/accounts/me/resourcegroups/%s",
		c.HostURL, resourceGroupId.ShortName), nil)
	if err != nil {
		return fmt.Errorf("building DELETE request for personal bounded context %q: %w", resourceGroupId.ShortName, err)
	}

	_, _, err = c.doRequest(ctx, req, []int{200, 404})
	if err != nil {
		return fmt.Errorf("deleting personal bounded context %q: %w", resourceGroupId.ShortName, err)
	}
	return nil
}
