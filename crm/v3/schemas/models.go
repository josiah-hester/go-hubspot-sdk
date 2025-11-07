package schemas

type DataSensitivity string

const (
	NonSensitive    DataSensitivity = "non_sensitive"
	Sensitive       DataSensitivity = "sensitive"
	HighlySensitive DataSensitivity = "highly_sensitive"
)

type NumberDisplayHint string

const (
	Unformatted NumberDisplayHint = "unformatted"
	Formatted   NumberDisplayHint = "formatted"
	Currency    NumberDisplayHint = "currency"
	Percentage  NumberDisplayHint = "percentage"
	Duration    NumberDisplayHint = "duration"
	Probability NumberDisplayHint = "probability"
)

type Schema struct {
	Associations []Association `json:"associations" required:"yes"`
	Labels       struct {
		Plural   string `json:"plural"`
		Singular string `json:"singular"`
	} `json:"labels" required:"yes"`
	RequiredProperties         []string   `json:"requiredProperties" required:"yes"`
	Name                       string     `json:"name" required:"yes"`
	ID                         string     `json:"id" required:"yes"`
	Properties                 []Property `json:"properties" required:"yes"`
	SecondaryDisplayProperties []string   `json:"secondaryDisplayProperties"`
	CreatedByUserID            int        `json:"createdByUserId"`
	ObjectTypeID               string     `json:"objectTypeID"`
	Description                string     `json:"description"`
	UpdatedByUserID            int        `json:"updatedByUserId"`
	FullyQualifiedName         string     `json:"fullyQualifiedName"`
	Archived                   bool       `json:"archived"`
	CreatedAt                  string     `json:"createdAt"`
	SearchableProperties       []string   `json:"searchableProperties"`
	PortalID                   int        `json:"portalId"`
	PrimaryDisplayProperty     string     `json:"primaryDisplayProperty"`
	UpdatedAt                  string     `json:"updatedAt"`
}

type Association struct {
	FromObjectTypeID string `json:"fromObjectTypeId" required:"yes"`
	ID               string `json:"id" required:"yes"`
	ToObjectTypeID   string `json:"toObjectTypeId" required:"yes"`
	CreatedAt        string `json:"createdAt"`
	Name             string `json:"name"`
	UpdatedAt        string `json:"updatedAt"`
}

type Property struct {
	Description          string   `json:"description" required:"yes"`
	Type                 string   `json:"type" required:"yes"`
	Options              []Option `json:"options" required:"yes"`
	Label                string   `json:"label" required:"yes"`
	GroupName            string   `json:"groupName" required:"yes"`
	Name                 string   `json:"name" required:"yes"`
	FieldType            string   `json:"fieldType" required:"yes"`
	Hidden               bool     `json:"hidden"`
	DisplayOrder         int      `json:"displayOrder"`
	ShowCurrencySymbol   bool     `json:"showCurrencySymbol"`
	HubspotDefined       bool     `json:"hubspotDefined"`
	CreatedAt            string   `json:"createdAt"`
	Archived             bool     `json:"archived"`
	HasUniqueValue       bool     `json:"hasUniqueValue"`
	Calculated           bool     `json:"calculated"`
	ExternalOptions      bool     `json:"externalOptions"`
	UdpatedAt            string   `json:"updatedAt"`
	CreatedUserID        string   `json:"createdUserId"`
	ModificationMetadata struct {
		ReadOnlyValue      bool `json:"readOnlyValue" required:"yes"`
		ReadOnlyDefinition bool `json:"readOnlyDefinition" required:"yes"`
		Archivable         bool `json:"archivable" required:"yes"`
		ReadOnlyOptions    bool `json:"readOnlyOptions"`
	} `json:"ModificationMetadata"`
	SensitiveDataCategories  []string           `json:"sensitiveDataCategories"`
	SearchableInGlobalSearch bool               `json:"searchableInGlobalSearch"`
	NumberDisplayHint        *NumberDisplayHint `json:"numberDisplayHint"`
	FormField                bool               `json:"formField"`
	DataSensitivity          *DataSensitivity   `json:"dataSensitivity"`
	ArchivedAt               string             `json:"archivedAt"`
	ReferencedObjectType     string             `json:"referencedObjectType"`
	CalculationFormula       string             `json:"CalculationFormula"`
	UpdatedUserID            string             `json:"updatedUserID"`
}

type Option struct {
	Hidden       bool   `json:"hidden" required:"yes"`
	Label        string `json:"label" required:"yes"`
	Value        string `json:"value" required:"yes"`
	DisplayOrder int    `json:"displayOrder"`
	Description  string `json:"description" `
}

type GetAllSchemasResponse struct {
	Results []Schema `json:"results" required:"yes"`
}

type CreateNewSchemaInput struct {
	RequiredProperties []string   `json:"requiredProperties" required:"yes"`
	Name               string     `json:"name" required:"yes"`
	AssociatedObjects  []string   `json:"associatedObjects" required:"yes"`
	Properties         []Property `json:"properties" required:"yes"`
	Labels             struct {
		Singular string `json:"singular"`
		Plural   string `json:"plural"`
	} `json:"labels"`
	SecondaryDisplayProperties []string `json:"secondaryDisplayProperties"`
	SearchableProperties       []string `json:"searchableProperties"`
	PrimaryDisplayProperty     string   `json:"primaryDisplayProperty"`
	Description                string   `json:"description"`
}

type CreateNewAssociationSchemaInput struct {
	FromObjectTypeID string `json:"fromObjectTypeId" required:"yes"`
	ToObjectTypeID   string `json:"toObjectTypeId" required:"yes"`
	Name             string `json:"name"`
}

type CreateNewAssociationSchemaResponse struct {
	FromObjectTypeID string `json:"fromObjectTypeId" required:"yes"`
	ID               string `json:"id" required:"yes"`
	ToObjectTypeID   string `json:"toObjectTypeId" required:"yes"`
	CreatedAt        string `json:"createdAt"`
	Name             string `json:"name"`
	UpdatedAt        string `json:"updatedAt"`
}

type UpdateSchemaInput struct {
	SecondaryDisplayProperties []string `json:"secondaryDisplayProperties"`
	RequiredProperties         []string `json:"requiredProperties"`
	SearchableProperties       []string `json:"searchableProperties"`
	ClearDescription           bool     `json:"clearDescription"`
	PrimaryDisplayProperty     string   `json:"primaryDisplayProperty"`
	Description                string   `json:"description"`
	Restorable                 bool     `json:"restorable"`
	Labels                     struct {
		Singular string `json:"singular"`
		Plural   string `json:"plural"`
	}
}
