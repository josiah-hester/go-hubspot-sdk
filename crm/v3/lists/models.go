package lists

import "time"

// ListProcessingType represents the processing type of a list
type ListProcessingType string

const (
	Manual   ListProcessingType = "MANUAL"
	Dynamic  ListProcessingType = "DYNAMIC"
	Snapshot ListProcessingType = "SNAPSHOT"
)

// ListProcessingStatus represents the processing status of a list
type ListProcessingStatus string

const (
	Processing ListProcessingStatus = "PROCESSING"
	Complete   ListProcessingStatus = "COMPLETE"
	Error      ListProcessingStatus = "ERROR"
)

// FilterType represents the type of filter
type FilterType string

const (
	Property                 FilterType = "PROPERTY"
	Association              FilterType = "ASSOCIATION"
	PageView                 FilterType = "PAGE_VIEW"
	CTA                      FilterType = "CTA"
	Event                    FilterType = "EVENT"
	FormSubmission           FilterType = "FORM_SUBMISSION"
	FormSubmissionOnPage     FilterType = "FORM_SUBMISSION_ON_PAGE"
	IntegrationEvent         FilterType = "INTEGRATION_EVENT"
	EmailSubscription        FilterType = "EMAIL_SUBSCRIPTION"
	CommunicationSubscription FilterType = "COMMUNICATION_SUBSCRIPTION"
	InList                   FilterType = "IN_LIST"
	NumAssociations          FilterType = "NUM_ASSOCIATIONS"
	UnifiedEvents            FilterType = "UNIFIED_EVENTS"
	PropertyAssociation      FilterType = "PROPERTY_ASSOCIATION"
	Webinar                  FilterType = "WEBINAR"
	EmailEvent               FilterType = "EMAIL_EVENT"
	Privacy                  FilterType = "PRIVACY"
	AdsSearch                FilterType = "ADS_SEARCH"
	AdsTime                  FilterType = "ADS_TIME"
	SurveyMonkey             FilterType = "SURVEY_MONKEY"
	SurveyMonkeyValue        FilterType = "SURVEY_MONKEY_VALUE"
	CampaignInfluenced       FilterType = "CAMPAIGN_INFLUENCED"
	Constant                 FilterType = "CONSTANT"
)

// FilterBranchType represents the type of filter branch
type FilterBranchType string

const (
	Or                        FilterBranchType = "OR"
	And                       FilterBranchType = "AND"
	NotAll                    FilterBranchType = "NOT_ALL"
	NotAny                    FilterBranchType = "NOT_ANY"
	Restricted                FilterBranchType = "RESTRICTED"
	UnifiedEventsBranch       FilterBranchType = "UNIFIED_EVENTS"
	PropertyAssociationBranch FilterBranchType = "PROPERTY_ASSOCIATION"
	AssociationBranch         FilterBranchType = "ASSOCIATION"
)

// AssociationCategory represents the category of an association
type AssociationCategory string

const (
	HubspotDefined    AssociationCategory = "HUBSPOT_DEFINED"
	UserDefined       AssociationCategory = "USER_DEFINED"
	IntegratorDefined AssociationCategory = "INTEGRATOR_DEFINED"
)

// FilterBranch represents a filter branch with nested logic
type FilterBranch struct {
	FilterBranchType     FilterBranchType `json:"filterBranchType"`
	FilterBranchOperator string           `json:"filterBranchOperator"`
	FilterBranches       []FilterBranch   `json:"filterBranches,omitempty"`
	Filters              []Filter         `json:"filters,omitempty"`

	// Association-specific fields
	AssociationTypeID   *int                 `json:"associationTypeId,omitempty"`
	AssociationCategory *AssociationCategory `json:"associationCategory,omitempty"`
	ObjectTypeID        *string              `json:"objectTypeId,omitempty"`
	Operator            *string              `json:"operator,omitempty"`

	// Property association-specific fields
	PropertyWithObjectID *string `json:"propertyWithObjectId,omitempty"`

	// Unified events-specific fields
	EventTypeID        *string        `json:"eventTypeId,omitempty"`
	CoalescingRefineBy map[string]any `json:"coalescingRefineBy,omitempty"`
}

// Filter represents a list filter
type Filter struct {
	FilterType FilterType     `json:"filterType"`
	Property   *string        `json:"property,omitempty"`
	Operation  map[string]any `json:"operation,omitempty"`
	Operator   *string        `json:"operator,omitempty"`

	// IN_LIST specific fields
	ListID   *string        `json:"listId,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`

	// Association specific fields
	AssociationTypeID   *int                 `json:"associationTypeId,omitempty"`
	AssociationCategory *AssociationCategory `json:"associationCategory,omitempty"`
	ToObjectType        *string              `json:"toObjectType,omitempty"`
	ToObjectTypeID      *string              `json:"toObjectTypeId,omitempty"`
	CoalescingRefineBy  map[string]any       `json:"coalescingRefineBy,omitempty"`

	// Form submission specific fields
	FormID          *string        `json:"formId,omitempty"`
	PruningRefineBy map[string]any `json:"pruningRefineBy,omitempty"`

	// Page view specific fields
	PageURL        *string `json:"pageUrl,omitempty"`
	EnableTracking *bool   `json:"enableTracking,omitempty"`

	// CTA specific fields
	CTAName *string `json:"ctaName,omitempty"`

	// Event specific fields
	EventID *string `json:"eventId,omitempty"`

	// Email subscription specific fields
	SubscriptionType *string  `json:"subscriptionType,omitempty"`
	SubscriptionIDs  []string `json:"subscriptionIds,omitempty"`
	AcceptedStatuses []string `json:"acceptedStatuses,omitempty"`

	// Communication subscription specific fields
	Channel           *string  `json:"channel,omitempty"`
	AcceptedOptStates []string `json:"acceptedOptStates,omitempty"`
	BusinessUnitID    *string  `json:"businessUnitId,omitempty"`

	// Webinar specific fields
	WebinarID *string `json:"webinarId,omitempty"`

	// Email event specific fields
	AppID    *string `json:"appId,omitempty"`
	EmailID  *string `json:"emailId,omitempty"`
	Level    *string `json:"level,omitempty"`
	ClickURL *string `json:"clickUrl,omitempty"`

	// Privacy specific fields
	PrivacyName *string `json:"privacyName,omitempty"`

	// Ads search specific fields
	SearchTerms    []string `json:"searchTerms,omitempty"`
	EntityType     *string  `json:"entityType,omitempty"`
	AdNetwork      *string  `json:"adNetwork,omitempty"`
	SearchTermType *string  `json:"searchTermType,omitempty"`

	// Survey Monkey specific fields
	SurveyID          *string        `json:"surveyId,omitempty"`
	SurveyQuestion    *string        `json:"surveyQuestion,omitempty"`
	SurveyAnswerRowID *string        `json:"surveyAnswerRowId,omitempty"`
	SurveyAnswerColID *string        `json:"surveyAnswerColId,omitempty"`
	ValueComparison   map[string]any `json:"valueComparison,omitempty"`

	// Campaign influenced specific fields
	CampaignID *string `json:"campaignId,omitempty"`

	// Constant filter specific fields
	ShouldAccept *bool   `json:"shouldAccept,omitempty"`
	Source       *string `json:"source,omitempty"`

	// Integration event specific fields
	EventTypeIDInt *int             `json:"eventTypeId,omitempty"`
	FilterLines    []map[string]any `json:"filterLines,omitempty"`

	// Property association specific fields
	PropertyWithObjectID *string `json:"propertyWithObjectId,omitempty"`
}

// MembershipSettings represents list membership settings
type MembershipSettings struct {
	MembershipTeamID  *int  `json:"membershipTeamId,omitempty"`
	IncludeUnassigned *bool `json:"includeUnassigned,omitempty"`
}

// ListPermissions represents list access permissions
type ListPermissions struct {
	TeamsWithEditAccess []int `json:"teamsWithEditAccess,omitempty"`
	UsersWithEditAccess []int `json:"usersWithEditAccess,omitempty"`
}

// List represents a HubSpot list object
type List struct {
	ListID             string               `json:"listId"`
	Name               string               `json:"name"`
	ObjectTypeID       string               `json:"objectTypeId"`
	ProcessingType     ListProcessingType   `json:"processingType"`
	ProcessingStatus   ListProcessingStatus `json:"processingStatus"`
	ListVersion        int                  `json:"listVersion"`
	Size               *int64               `json:"size,omitempty"`
	CreatedAt          time.Time            `json:"createdAt"`
	UpdatedAt          time.Time            `json:"updatedAt"`
	DeletedAt          *time.Time           `json:"deletedAt,omitempty"`
	CreatedByID        string               `json:"createdById"`
	UpdatedByID        string               `json:"updatedById"`
	FiltersUpdatedAt   *time.Time           `json:"filtersUpdatedAt,omitempty"`
	MembershipSettings *MembershipSettings  `json:"membershipSettings,omitempty"`
	ListPermissions    *ListPermissions     `json:"listPermissions,omitempty"`
	FilterBranch       *FilterBranch        `json:"filterBranch,omitempty"`
}

// ListCreateRequest represents the request to create a list
type ListCreateRequest struct {
	Name               string              `json:"name"`
	ObjectTypeID       string              `json:"objectTypeId"`
	ProcessingType     ListProcessingType  `json:"processingType"`
	FilterBranch       *FilterBranch       `json:"filterBranch,omitempty"`
	MembershipSettings *MembershipSettings `json:"membershipSettings,omitempty"`
	ListPermissions    *ListPermissions    `json:"listPermissions,omitempty"`
	ListFolderID       *int                `json:"listFolderId,omitempty"`
	CustomProperties   map[string]string   `json:"customProperties,omitempty"`
}

// ListResponse represents the response from creating a list
type ListResponse struct {
	List List `json:"list"`
}

// ListUpdateNameRequest represents the request to update a list name
type ListUpdateNameRequest struct {
	ListName       string `json:"listName"`
	IncludeFilters *bool  `json:"includeFilters,omitempty"`
}

// ListUpdateFiltersRequest represents the request to update list filters
type ListUpdateFiltersRequest struct {
	FilterBranch   FilterBranch `json:"filterBranch"`
	IncludeFilters *bool        `json:"includeFilters,omitempty"`
}

// ListSearchRequest represents a request to search for lists
type ListSearchRequest struct {
	Query                *string              `json:"query,omitempty"`
	ProcessingTypes      []ListProcessingType `json:"processingTypes,omitempty"`
	AdditionalProperties []string             `json:"additionalProperties,omitempty"`
	Count                *int                 `json:"count,omitempty"`
	Offset               *int                 `json:"offset,omitempty"`
}

// ListSearchResponse represents the response from searching lists
type ListSearchResponse struct {
	Lists   []List `json:"lists"`
	Total   int    `json:"total"`
	HasMore bool   `json:"hasMore"`
	Offset  int    `json:"offset"`
}

// ListsByIDResponse represents the response from fetching multiple lists
type ListsByIDResponse struct {
	Lists []List `json:"lists"`
}

// RecordListMembership represents a record's membership in a list
type RecordListMembership struct {
	ListID              string    `json:"listId"`
	ListVersion         int       `json:"listVersion"`
	FirstAddedTimestamp time.Time `json:"firstAddedTimestamp"`
	LastAddedTimestamp  time.Time `json:"lastAddedTimestamp"`
	IsPublicList        *bool     `json:"isPublicList,omitempty"`
}

// RecordMembershipsResponse represents the response with list memberships
type RecordMembershipsResponse struct {
	Results []RecordListMembership `json:"results"`
	Total   *int64                 `json:"total,omitempty"`
}

// MembershipAddRequest represents the request to add records to a list
type MembershipAddRequest []string

// MembershipRemoveRequest represents the request to remove records from a list
type MembershipRemoveRequest []string

// MembershipAddAndRemoveRequest represents the request to add and remove records
type MembershipAddAndRemoveRequest struct {
	RecordIDsToAdd    []string `json:"recordIdsToAdd"`
	RecordIDsToRemove []string `json:"recordIdsToRemove"`
}

// MembershipChangeResponse represents the response from adding/removing records
type MembershipChangeResponse struct {
	RecordIDsAdded   []string `json:"recordIdsAdded,omitempty"`
	RecordIDsRemoved []string `json:"recordIdsRemoved,omitempty"`
}

// ListMembershipsResponse represents the response from fetching list members
type ListMembershipsResponse struct {
	Results []string `json:"results"`
	HasMore *bool    `json:"hasMore,omitempty"`
	Offset  *string  `json:"offset,omitempty"`
}

// BatchReadMembershipsRequest represents a batch request to read memberships
type BatchReadMembershipsRequest struct {
	Inputs []MembershipRecordIdentifier `json:"inputs"`
}

// MembershipRecordIdentifier identifies a record for membership operations
type MembershipRecordIdentifier struct {
	ObjectTypeID string `json:"objectTypeId"`
	RecordID     string `json:"recordId"`
}

// BatchReadMembershipsResponse represents the batch response
type BatchReadMembershipsResponse struct {
	Results []RecordMembershipsResponse `json:"results"`
}

// ConversionType represents the type of list conversion
type ConversionType string

const (
	ConversionDate ConversionType = "CONVERSION_DATE"
	Inactivity     ConversionType = "INACTIVITY"
)

// TimeUnit represents a unit of time for inactivity conversion
type TimeUnit string

const (
	Day   TimeUnit = "DAY"
	Week  TimeUnit = "WEEK"
	Month TimeUnit = "MONTH"
)

// ScheduleConversionRequest represents a request to schedule list conversion
type ScheduleConversionRequest struct {
	ConversionType ConversionType `json:"conversionType"`

	// For CONVERSION_DATE
	Year  *int `json:"year,omitempty"`
	Month *int `json:"month,omitempty"`
	Day   *int `json:"day,omitempty"`

	// For INACTIVITY
	TimeUnit *TimeUnit `json:"timeUnit,omitempty"`
	Offset   *int      `json:"offset,omitempty"`
}

// ScheduledConversionTime represents the scheduled conversion time
type ScheduledConversionTime struct {
	ConversionType ConversionType `json:"conversionType"`

	// For CONVERSION_DATE
	Year  *int `json:"year,omitempty"`
	Month *int `json:"month,omitempty"`
	Day   *int `json:"day,omitempty"`

	// For INACTIVITY
	TimeUnit *TimeUnit `json:"timeUnit,omitempty"`
	Offset   *int      `json:"offset,omitempty"`
}

// ScheduleConversionResponse represents the response from scheduling conversion
type ScheduleConversionResponse struct {
	ListID                  string                  `json:"listId"`
	RequestedConversionTime ScheduledConversionTime `json:"requestedConversionTime"`
	ConvertedAt             *time.Time              `json:"convertedAt,omitempty"`
}

// ListIDMappingRequest represents a request to map legacy list IDs
type ListIDMappingRequest struct {
	LegacyListIDs []string `json:"legacyListIds"`
}

// ListIDMapping represents a mapping between legacy and v3 list IDs
type ListIDMapping struct {
	LegacyListID string `json:"legacyListId"`
	ListID       string `json:"listId"`
}

// ListIDMappingResponse represents the response from list ID mapping
type ListIDMappingResponse struct {
	Mappings             []ListIDMapping `json:"mappings"`
	MissingLegacyListIDs []string        `json:"missingLegacyListIds,omitempty"`
}

// FolderCreateRequest represents a request to create a list folder
type FolderCreateRequest struct {
	Name           string `json:"name"`
	ParentFolderID *int   `json:"parentFolderId,omitempty"`
}

// Folder represents a list folder
type Folder struct {
	FolderID       int    `json:"folderId"`
	Name           string `json:"name"`
	ParentFolderID *int   `json:"parentFolderId,omitempty"`
}
