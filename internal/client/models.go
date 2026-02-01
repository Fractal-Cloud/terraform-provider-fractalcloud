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
	ComponentId string                 `json:"componentId"`
	Settings    map[string]interface{} `json:"settings"`
}

type Component struct {
	DisplayName       string                 `json:"displayName"`
	Description       string                 `json:"description"`
	Type              string                 `json:"type"`
	Id                string                 `json:"id"`
	Version           string                 `json:"version"`
	IsLocked          bool                   `json:"locked"`
	RecreateOnFailure bool                   `json:"recreateOnFailure"`
	Parameters        map[string]interface{} `json:"parameters"`
	Dependencies      []string               `json:"dependencies"`
	Links             []ComponentLink        `json:"links"`
	OutputFields      []string               `json:"outputFields"`
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
