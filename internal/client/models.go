package fractalCloud

// ResourceGroup -
type ResourceGroup struct {
	ID          ResourceGroupId `json:"id"`
	DisplayName string          `json:"displayName"`
	Description string          `json:"description"`
	CreatedAt   string          `json:"createdAt"`
	UpdatedAt   string          `json:"updatedAt"`
}

type ResourceGroupId struct {
	Type      string `json:"resourceGroupType"`
	OwnerId   string `json:"ownerId"`
	ShortName string `json:"shortName"`
}
