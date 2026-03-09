package fractalCloud

type ResourceGroupId struct {
	Type      string `json:"resourceGroupType"`
	OwnerId   string `json:"ownerId"`
	ShortName string `json:"shortName"`
}

type PersonalResourceGroup struct {
	Id             ResourceGroupId `json:"id"`
	DisplayName    string          `json:"displayName"`
	Status         string          `json:"status"`
	Description    string          `json:"description"`
	Icon           string          `json:"icon"`
	LiveSystemsIds []string        `json:"livesystems"`
	FractalsIds    []string        `json:"fractals"`
	CreatedAt      string          `json:"createdAt"`
	UpdatedAt      string          `json:"updatedAt"`
}

type OrganizationalResourceGroup struct {
	Id             ResourceGroupId `json:"id"`
	DisplayName    string          `json:"displayName"`
	Description    string          `json:"description"`
	Icon           string          `json:"icon"`
	Status         string          `json:"status"`
	MembersIds     []string        `json:"membersIds"`
	TeamsIds       []string        `json:"teamsIds"`
	ManagersIds    []string        `json:"managersIds"`
	LiveSystemsIds []string        `json:"livesystems"`
	FractalsIds    []string        `json:"fractals"`
	CreatedAt      string          `json:"createdAt"`
	CreatedBy      string          `json:"createdBy"`
	UpdatedAt      string          `json:"updatedAt"`
	UpdatedBy      string          `json:"updatedBy"`
}

type Organization struct {
	Id                string   `json:"id"`
	DisplayName       string   `json:"name"`
	Description       string   `json:"description"`
	Icon              string   `json:"icon"`
	AdminsIds         []string `json:"admins"`
	MembersIds        []string `json:"members"`
	ResourceGroupsIds []string `json:"resourceGroups"`
	TeamsIds          []string `json:"teams"`
	SocialLinks       []string `json:"socialLinks"`
	Tags              []string `json:"tags"`
	Status            string   `json:"status"`
	SubscriptionId    string   `json:"subscriptionId"`
	CreatedAt         string   `json:"createdAt"`
	CreatedBy         string   `json:"createdBy"`
	UpdatedAt         string   `json:"updatedAt"`
	UpdatedBy         string   `json:"updatedBy"`
}

type ComponentLink struct {
	ComponentId string            `json:"componentId" tfsdk:"component_id"`
	Settings    map[string]string `json:"settings" tfsdk:"settings"`
}

type Component struct {
	Id                string            `json:"id" tfsdk:"id"`
	Type              string            `json:"type" tfsdk:"type"`
	DisplayName       *string           `json:"displayName,omitempty" tfsdk:"display_name"`
	Description       *string           `json:"description,omitempty" tfsdk:"description"`
	Version           *string           `json:"version,omitempty" tfsdk:"version"`
	IsLocked          *bool             `json:"locked,omitempty" tfsdk:"is_locked"`
	RecreateOnFailure *bool             `json:"recreateOnFailure,omitempty" tfsdk:"recreate_on_failure"`
	Parameters        map[string]string `json:"parameters,omitempty" tfsdk:"parameters"`
	DependenciesIds   []string          `json:"dependencies,omitempty" tfsdk:"dependencies_ids"`
	Links             []ComponentLink   `json:"links,omitempty" tfsdk:"links"`
	OutputFields      []string          `json:"outputFields,omitempty" tfsdk:"output_fields"`
}

type Blueprint struct {
	FractalId   string      `json:"fractalId"`
	IsPrivate   bool        `json:"isPrivate"`
	Status      string      `json:"status"`
	ReasonCode  string      `json:"reasonCode"`
	Description string      `json:"description"`
	Components  []Component `json:"components"`
	CreatedAt   string      `json:"created"`
}
