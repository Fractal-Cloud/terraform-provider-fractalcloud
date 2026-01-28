package fractalCloud

// ResourceGroup -
type ResourceGroup struct {
	ID          ResourceGroupId `json:"id"`
	DisplayName string          `json:"displayName"`
	Description string          `json:"description"`
	Icon        string          `json:"icon"`
	CreatedAt   string          `json:"createdAt"`
	CreatedBy   string          `json:"createdBy"`
	UpdatedAt   string          `json:"updatedAt"`
	UpdatedBy   string          `json:"updatedBy"`
}

type ResourceGroupId struct {
	Type      string `json:"resourceGroupType"`
	OwnerId   string `json:"ownerId"`
	ShortName string `json:"shortName"`
}
