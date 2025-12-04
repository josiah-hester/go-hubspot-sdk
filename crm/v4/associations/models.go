package associations

// Association category constants
const (
	// AssociationCategoryHubSpotDefined represents HubSpot's predefined associations
	AssociationCategoryHubSpotDefined = "HUBSPOT_DEFINED"

	// AssociationCategoryIntegratorDefined represents integrator-defined associations
	AssociationCategoryIntegratorDefined = "INTEGRATOR_DEFINED"

	// AssociationCategoryUserDefined represents user-defined associations
	AssociationCategoryUserDefined = "USER_DEFINED"
)

// AssociationSpec defines an association between two objects
type AssociationSpec struct {
	AssociationCategory string `json:"associationCategory"`
	AssociationTypeID   int    `json:"associationTypeId"`
}

// AssociationLabel represents an association label/type
type AssociationLabel struct {
	Category string `json:"category"`
	TypeID   int    `json:"typeId"`
	Label    string `json:"label"`
}

// CreateAssociationInput represents input for creating an association
type CreateAssociationInput struct {
	Inputs []AssociationInput `json:"inputs"`
}

// AssociationInput represents a single association to create
type AssociationInput struct {
	From  AssociationEndpoint `json:"from"`
	To    AssociationEndpoint `json:"to"`
	Types []AssociationSpec   `json:"types"`
}

// AssociationEndpoint represents an object in an association
type AssociationEndpoint struct {
	ID string `json:"id"`
}

// BatchAssociationInput represents input for batch association operations
type BatchAssociationInput struct {
	Inputs []struct {
		From struct {
			ID string `json:"id"`
		} `json:"from"`
		To []struct {
			ID    string            `json:"id"`
			Types []AssociationSpec `json:"types"`
		} `json:"to"`
	} `json:"inputs"`
}

// AssociationResponse represents response from creating/updating association operations
type AssociationResponse struct {
	FromObjectTypeID string   `json:"fromObjectTypeId"`
	FromObjectID     int      `json:"fromObjectId"`
	ToObjectTypeID   string   `json:"toObjectTypeId"`
	ToObjectID       int      `json:"toObjectId"`
	Labels           []string `json:"labels"`
}

// ListAssociationsResponse represents response from listing associations
type ListAssociationsResponse struct {
	Results []AssociatedObject `json:"results"`
	Paging  *Paging            `json:"paging"`
}

// AssociatedObject represents an associated object
type AssociatedObject struct {
	ToObjectID       string            `json:"toObjectId"`
	AssociationTypes []AssociationSpec `json:"associationTypes"`
}

// Paging represents pagination information
type Paging struct {
	Next *PagingLink `json:"next"`
	Prev *PagingLink `json:"prev"`
}

// PagingLink represents a pagination link
type PagingLink struct {
	After string `json:"after"`
	Link  string `json:"link"`
}

// BatchAssociationResponse represents response from batch operations
type BatchAssociationResponse struct {
	Status      string   `json:"status"`
	Results     []string `json:"results"`
	StartedAt   string   `json:"startedAt"`
	CompletedAt string   `json:"completedAt"`
}

// GetAssociationLabelsResponse represents response from getting association labels
type GetAssociationLabelsResponse struct {
	Results []AssociationLabel `json:"results"`
}
