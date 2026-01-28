package fractalCloud

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetResourceGroup - Returns specific resource group
func (c *Client) GetResourceGroup(resourceGroupID ResourceGroupId) (ResourceGroup, error) {
	relativePath := ""
	if resourceGroupID.Type == "Personal" {
		relativePath = fmt.Sprintf("accounts/me/resourcegroups/%s", resourceGroupID.ShortName)
	} else {
		relativePath = fmt.Sprintf("organizations/%s/resourcegroups/%s", resourceGroupID.OwnerId, resourceGroupID.ShortName)
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.HostURL, relativePath), nil)
	if err != nil {
		return ResourceGroup{}, err
	}

	body, err := c.doRequest(req, []int{200})
	if err != nil {
		return ResourceGroup{}, err
	}

	resourceGroup := ResourceGroup{}
	err = json.Unmarshal(body, &resourceGroup)
	if err != nil {
		return ResourceGroup{}, err
	}

	return resourceGroup, nil
}
