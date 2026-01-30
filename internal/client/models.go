package fractalCloud

type ResourceGroupId struct {
	Type      string `json:"resourceGroupType"`
	OwnerID   string `json:"ownerId"`
	ShortName string `json:"shortName"`
}

type PersonalResourceGroup struct {
	ID             ResourceGroupId `json:"id"`
	DisplayName    string          `json:"displayName"`
	Status         string          `json:"status"`
	Description    string          `json:"description"`
	Icon           string          `json:"icon"`
	LiveSystemsIds []string        `json:"liveSystems"`
	FractalsIds    []string        `json:"fractals"`
	CreatedAt      string          `json:"createdAt"`
	UpdatedAt      string          `json:"updatedAt"`
}

type OrganizationalResourceGroup struct {
	ID             ResourceGroupId `json:"id"`
	DisplayName    string          `json:"displayName"`
	Description    string          `json:"description"`
	Icon           string          `json:"icon"`
	Status         string          `json:"status"`
	MembersIds     []string        `json:"membersIds"`
	TeamsIds       []string        `json:"teamsIds"`
	ManagersIds    []string        `json:"managersIds"`
	LiveSystemsIds []string        `json:"liveSystems"`
	FractalsIds    []string        `json:"fractals"`
	CreatedAt      string          `json:"createdAt"`
	CreatedBy      string          `json:"createdBy"`
	UpdatedAt      string          `json:"updatedAt"`
	UpdatedBy      string          `json:"updatedBy"`
}

type Organization struct {
	ID                string   `json:"id"`
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
