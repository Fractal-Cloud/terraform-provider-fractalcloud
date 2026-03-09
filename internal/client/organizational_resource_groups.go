package fractalCloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetOrganizationalResourceGroup - Returns specific organizational resource group (bounded context).
func (c *Client) GetOrganizationalResourceGroup(ctx context.Context, resourceGroupId ResourceGroupId) (*OrganizationalResourceGroup, error) {
	path := fmt.Sprintf("%s/organizations/%s/resourcegroups/%s",
		c.HostURL, resourceGroupId.OwnerId, resourceGroupId.ShortName)

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("building GET request for organizational bounded context %q (org %q): %w",
			resourceGroupId.ShortName, resourceGroupId.OwnerId, err)
	}

	resCode, body, err := c.doRequest(ctx, req, []int{200, 404})
	if err != nil {
		return nil, fmt.Errorf("fetching organizational bounded context %q (org %q): %w",
			resourceGroupId.ShortName, resourceGroupId.OwnerId, err)
	}

	if resCode == 404 {
		c.logDebug(fmt.Sprintf("organizational bounded context %q (org %q) not found",
			resourceGroupId.ShortName, resourceGroupId.OwnerId))
		return nil, nil
	}

	resourceGroup := OrganizationalResourceGroup{}
	if err := json.Unmarshal(body, &resourceGroup); err != nil {
		return nil, fmt.Errorf("decoding organizational bounded context %q (org %q) response: %w",
			resourceGroupId.ShortName, resourceGroupId.OwnerId, err)
	}

	return &resourceGroup, nil
}

type UpsertOrganizationalResourceGroupRequestBody struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	IconUrl     string `json:"iconUrl"`
}

func (c *Client) UpsertOrganizationalResourceGroup(ctx context.Context, resourceGroup OrganizationalResourceGroup) error {
	resourceGroupId := resourceGroup.Id

	requestBody := UpsertOrganizationalResourceGroupRequestBody{
		DisplayName: resourceGroup.DisplayName,
		Description: resourceGroup.Description,
	}
	rb, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("encoding organizational bounded context %q (org %q) request: %w",
			resourceGroupId.ShortName, resourceGroupId.OwnerId, err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/organizations/%s/resourcegroups/%s",
		c.HostURL, resourceGroupId.OwnerId, resourceGroupId.ShortName), strings.NewReader(string(rb)))
	if err != nil {
		return fmt.Errorf("building POST request for organizational bounded context %q (org %q): %w",
			resourceGroupId.ShortName, resourceGroupId.OwnerId, err)
	}

	_, _, err = c.doRequest(ctx, req, []int{200, 202})
	if err != nil {
		return fmt.Errorf("upserting organizational bounded context %q (org %q): %w",
			resourceGroupId.ShortName, resourceGroupId.OwnerId, err)
	}
	return nil
}

func (c *Client) DeleteOrganizationalResourceGroup(ctx context.Context, resourceGroupId ResourceGroupId) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/organizations/%s/resourcegroups/%s",
		c.HostURL, resourceGroupId.OwnerId, resourceGroupId.ShortName), nil)
	if err != nil {
		return fmt.Errorf("building DELETE request for organizational bounded context %q (org %q): %w",
			resourceGroupId.ShortName, resourceGroupId.OwnerId, err)
	}

	_, _, err = c.doRequest(ctx, req, []int{200, 404})
	if err != nil {
		return fmt.Errorf("deleting organizational bounded context %q (org %q): %w",
			resourceGroupId.ShortName, resourceGroupId.OwnerId, err)
	}
	return nil
}
