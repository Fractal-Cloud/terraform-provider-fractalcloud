package fractalCloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type FractalId struct {
	ResourceGroupId ResourceGroupId
	Name            string
	Version         string
}

func (id *FractalId) ToString() string {
	return id.ResourceGroupId.ToString() + "/" + id.Name + ":" + id.Version
}

func (id *ResourceGroupId) ToString() string {
	return id.Type + "/" + id.OwnerId + "/" + id.ShortName
}

// GetBlueprint - Returns specific organizational resource group
func (c *Client) GetBlueprint(id FractalId) (*Blueprint, error) {
	resourceGroupId := id.ResourceGroupId
	path := fmt.Sprintf("%s/blueprints/%s/%s/%s/%s/%s",
		c.HostURL, resourceGroupId.Type, resourceGroupId.OwnerId,
		resourceGroupId.ShortName, id.Name, id.Version)
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

	blueprint := BlueprintInternal{}
	err = json.Unmarshal(body, &blueprint)
	if err != nil {
		return nil, err
	}

	normalizedComponents := make([]Component, len(blueprint.Components))
	for i, component := range blueprint.Components {
		links := make([]ComponentLink, len(component.Links))
		for j, link := range component.Links {
			links[j] = ComponentLink{
				ComponentId: link.ComponentId,
				Settings:    MapAnyToMapStringJSON(c.Logger, link.Settings),
			}
		}
		normalizedComponents[i] = Component{
			Id:                component.Id,
			Type:              component.Type,
			DisplayName:       component.DisplayName,
			Description:       component.Description,
			Version:           component.Version,
			IsLocked:          component.IsLocked,
			RecreateOnFailure: component.RecreateOnFailure,
			Parameters:        MapAnyToMapStringJSON(c.Logger, component.Parameters),
			DependenciesIds:   component.DependenciesIds,
			Links:             links,
			OutputFields:      component.OutputFields,
		}
	}

	return &Blueprint{
		FractalId:   blueprint.FractalId,
		IsPrivate:   blueprint.IsPrivate,
		Status:      blueprint.Status,
		ReasonCode:  blueprint.ReasonCode,
		Description: blueprint.Description,
		Components:  normalizedComponents,
		CreatedAt:   blueprint.CreatedAt,
	}, nil
}

type ComponentLinkInternal struct {
	ComponentId string                 `json:"componentId"`
	Settings    map[string]interface{} `json:"settings"`
}

type ComponentInternal struct {
	Id                string                  `json:"id"`
	Type              string                  `json:"type"`
	DisplayName       string                  `json:"displayName"`
	Description       string                  `json:"description"`
	Version           string                  `json:"version"`
	IsLocked          bool                    `json:"locked"`
	RecreateOnFailure bool                    `json:"recreateOnFailure"`
	Parameters        map[string]interface{}  `json:"parameters"`
	DependenciesIds   []string                `json:"dependencies"`
	Links             []ComponentLinkInternal `json:"links"`
	OutputFields      []string                `json:"outputFields"`
}

type BlueprintInternal struct {
	FractalId   string              `json:"fractalId"`
	IsPrivate   bool                `json:"isPrivate"`
	Status      string              `json:"status"`
	ReasonCode  string              `json:"reasonCode"`
	Description string              `json:"description"`
	Components  []ComponentInternal `json:"components"`
	CreatedAt   string              `json:"created"`
}

func MapAnyToMapStringJSON(logger *ClientLogger, in map[string]interface{}) map[string]string {
	out := make(map[string]string, len(in))

	for k, v := range in {
		b, err := json.Marshal(v)
		if err == nil {
			out[k] = string(b)
		} else {
			logger.Warning(fmt.Sprintf("marshal key %q: %w", k, err))
		}
	}

	return out
}

type CreateBlueprintCommandRequestBody struct {
	Description string      `json:"description"`
	IsPrivate   bool        `json:"isPrivate"`
	Components  []Component `json:"components"`
}

func (c *Client) CreateBlueprint(id FractalId, description string, isPrivate bool, components []Component) error {
	resourceGroupId := id.ResourceGroupId

	requestBody := CreateBlueprintCommandRequestBody{
		Description: description,
		IsPrivate:   isPrivate,
		Components:  components,
	}

	rb, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/blueprints/%s/%s/%s/%s/%s",
		c.HostURL, resourceGroupId.Type, resourceGroupId.OwnerId,
		resourceGroupId.ShortName, id.Name, id.Version)

	req, err := http.NewRequest("POST", path, strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, _, err = c.doRequest(req, []int{200})
	return err
}

type UpdateBlueprintCommandRequest struct {
	ResourceGroupId ResourceGroupId `json:"resourceGroupId"`
	FractalName     string          `json:"fractalName"`
	FractalVersion  string          `json:"fractalVersion"`
	Description     string          `json:"description"`
	IsPrivate       bool            `json:"isPrivate"`
	Components      []Component     `json:"components"`
}

func (c *Client) UpdateBlueprint(id FractalId, newId FractalId, description string, isPrivate bool, components []Component) error {
	resourceGroupId := id.ResourceGroupId

	requestBody := UpdateBlueprintCommandRequest{
		ResourceGroupId: resourceGroupId,
		FractalName:     newId.Name,
		FractalVersion:  newId.Version,
		Description:     description,
		IsPrivate:       isPrivate,
		Components:      components,
	}

	rb, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/blueprints/%s/%s/%s/%s/%s",
		c.HostURL, resourceGroupId.Type, resourceGroupId.OwnerId,
		resourceGroupId.ShortName, id.Name, id.Version)

	req, err := http.NewRequest("PUT", path, strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, _, err = c.doRequest(req, []int{200})
	return err
}

func (c *Client) DeleteBlueprint(id FractalId) error {
	resourceGroupId := id.ResourceGroupId

	path := fmt.Sprintf("%s/blueprints/%s/%s/%s/%s/%s",
		c.HostURL, resourceGroupId.Type, resourceGroupId.OwnerId,
		resourceGroupId.ShortName, id.Name, id.Version)

	req, err := http.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, _, err = c.doRequest(req, []int{200, 404})
	if err != nil {
		return err
	}

	return nil
}
