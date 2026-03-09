package fractalCloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetOrganization - Returns specific organization.
func (c *Client) GetOrganization(ctx context.Context, organizationId string) (*Organization, error) {
	path := fmt.Sprintf("%s/organizations/%s", c.HostURL, organizationId)

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("building GET request for organization %q: %w", organizationId, err)
	}

	resCode, body, err := c.doRequest(ctx, req, []int{200, 404})
	if err != nil {
		return nil, fmt.Errorf("fetching organization %q: %w", organizationId, err)
	}

	if resCode == 404 {
		c.logDebug(fmt.Sprintf("organization %q not found", organizationId))
		return nil, nil
	}

	organization := Organization{}
	if err := json.Unmarshal(body, &organization); err != nil {
		return nil, fmt.Errorf("decoding organization %q response: %w", organizationId, err)
	}

	return &organization, nil
}
