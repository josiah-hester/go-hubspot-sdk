package companies

// Company represents a HubSpot company object
type Company struct {
	ID                    string                           `json:"id"`
	Properties            map[string]string                `json:"properties"`
	PropertiesWithHistory map[string][]PropertyWithHistory `json:"propertiesWithHistory"`
	CreatedAt             string                           `json:"createdAt"`
	UpdatedAt             string                           `json:"updatedAt"`
	Archived              bool                             `json:"archived"`
	ArchivedAt            string                           `json:"archivedAt"`
}

// PropertyWithHistory represents a property with its historical values
type PropertyWithHistory struct {
	Value           string `json:"value"`
	Timestamp       string `json:"timestamp"`
	SourceType      string `json:"sourceType"`
	SourceID        string `json:"sourceId"`
	SourceLabel     string `json:"sourceLabel"`
	UpdatedByUserID int    `json:"updatedByUserId"`
}

// CreateCompanyInput represents the input for creating a company
type CreateCompanyInput struct {
	Properties map[string]string `json:"properties"`
}

// UpdateCompanyInput represents the input for updating a company
type UpdateCompanyInput struct {
	Properties map[string]string `json:"properties"`
}

// ListCompaniesResponse represents the response from listing companies
type ListCompaniesResponse struct {
	Results []Company `json:"results"`
	Paging  *Paging   `json:"paging"`
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

// BatchReadCompaniesInput represents input for batch read
type BatchReadCompaniesInput struct {
	Properties            []string `json:"properties"`
	PropertiesWithHistory []string `json:"propertiesWithHistory"`
	IDProperty            string   `json:"idProperty"`
	Inputs                []struct {
		ID string `json:"id"`
	} `json:"inputs"`
}

// BatchCreateCompaniesInput represents input for batch create
type BatchCreateCompaniesInput struct {
	Inputs []CreateCompanyInput `json:"inputs"`
}

// BatchUpdateCompaniesInput represents input for batch update
type BatchUpdateCompaniesInput struct {
	Inputs []struct {
		ID         string            `json:"id"`
		Properties map[string]string `json:"properties"`
	} `json:"inputs"`
}

// BatchArchiveCompaniesInput represents input for batch archive
type BatchArchiveCompaniesInput struct {
	Inputs []struct {
		ID string `json:"id"`
	} `json:"inputs"`
}

// BatchCompaniesResponse represents response from batch operations
type BatchCompaniesResponse struct {
	Status      string    `json:"status"`
	Results     []Company `json:"results"`
	StartedAt   string    `json:"startedAt"`
	CompletedAt string    `json:"completedAt"`
}

// SearchCompaniesInput represents input for searching companies
type SearchCompaniesInput struct {
	FilterGroups []FilterGroup `json:"filterGroups"`
	Sorts        []string      `json:"sorts"`
	Query        string        `json:"query"`
	Properties   []string      `json:"properties"`
	Limit        int           `json:"limit"`
	After        string        `json:"after"`
}

// FilterGroup represents a group of filters
type FilterGroup struct {
	Filters []Filter `json:"filters"`
}

// Filter represents a single filter
type Filter struct {
	PropertyName string      `json:"propertyName"`
	Operator     string      `json:"operator"`
	Value        interface{} `json:"value"`
}

// SearchCompaniesResponse represents response from search
type SearchCompaniesResponse struct {
	Total   int       `json:"total"`
	Results []Company `json:"results"`
	Paging  *Paging   `json:"paging"`
}
