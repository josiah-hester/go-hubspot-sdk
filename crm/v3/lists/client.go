// Package lists specifies the client methods for the HubSpot CRM Lists API
package lists

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aacc-dev/go-hubspot-sdk/client"
)

type Client struct {
	apiClient *client.Client
}

// NewClient creates a new lists client
func NewClient(apiClient *client.Client) *Client {
	return &Client{
		apiClient: apiClient,
	}
}

func (c *Client) GetListByID(ctx context.Context, listID string, opts ...GetListOption) (*List, error) {
	req := client.NewRequest("GET", fmt.Sprintf("/crm/v3/lists/%s", listID))
	req.WithContext(ctx)
	req.WithResourceType("lists")

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseListError(err, listID)
	}

	var list ListResponse
	if err := json.Unmarshal(resp.Body, &list); err != nil {
		return nil, fmt.Errorf("failed to unmarshal list response: %w", err)
	}

	return &list.List, nil
}

func (c *Client) GetListByName(ctx context.Context, ObjectTypeID, listName string, opts ...GetListOption) (*List, error) {
	req := client.NewRequest("GET", fmt.Sprintf("/crm/v3/lists/object-type-id/%s/name/%s", ObjectTypeID, listName))
	req.WithContext(ctx)
	req.WithResourceType("lists")

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseListError(err, listName)
	}

	var list ListResponse
	if err := json.Unmarshal(resp.Body, &list); err != nil {
		return nil, fmt.Errorf("failed to unmarshal list response: %w", err)
	}

	return &list.List, nil
}

func (c *Client) CreateList(ctx context.Context, input *ListCreateRequest) (*List, error) {
	req := client.NewRequest("POST", "/crm/v3/lists")
	req.WithContext(ctx)
	req.WithResourceType("lists")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseListError(err, "")
	}

	var listResp ListResponse
	if err := json.Unmarshal(resp.Body, &listResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal list response: %w", err)
	}

	return &listResp.List, nil
}

func (c *Client) GetListsByIDs(ctx context.Context, listIDs []string, opts ...GetListOption) ([]List, error) {
	req := client.NewRequest("GET", "/crm/v3/lists")
	req.WithContext(ctx)
	req.WithResourceType("lists")

	// Add listIds as query parameters
	for _, id := range listIDs {
		req.AddQueryParam("listIds", id)
	}

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseListError(err, "")
	}

	var listsResp ListsByIDResponse
	if err := json.Unmarshal(resp.Body, &listsResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal lists response: %w", err)
	}

	return listsResp.Lists, nil
}

func (c *Client) SearchLists(ctx context.Context, input *ListSearchRequest) (*ListSearchResponse, error) {
	req := client.NewRequest("POST", "/crm/v3/lists/search")
	req.WithContext(ctx)
	req.WithResourceType("lists")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseListError(err, "")
	}

	var searchResp ListSearchResponse
	if err := json.Unmarshal(resp.Body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %w", err)
	}

	return &searchResp, nil
}

func (c *Client) UpdateListName(ctx context.Context, listID, listName string, includeFilters bool) (*List, error) {
	req := client.NewRequest("PUT", fmt.Sprintf("/crm/v3/lists/%s/update-list-name", listID))
	req.WithContext(ctx)
	req.WithResourceType("lists")
	req.AddQueryParam("listName", listName)

	if includeFilters {
		req.AddQueryParam("includeFilters", "true")
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseListError(err, listID)
	}

	var listResp ListResponse
	if err := json.Unmarshal(resp.Body, &listResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal list response: %w", err)
	}

	return &listResp.List, nil
}

func (c *Client) UpdateListFilters(ctx context.Context, listID string, filterBranch FilterBranch, includeFilters bool) (*List, error) {
	req := client.NewRequest("PUT", fmt.Sprintf("/crm/v3/lists/%s/update-list-filters", listID))
	req.WithContext(ctx)
	req.WithResourceType("lists")

	body := map[string]any{
		"filterBranch": filterBranch,
	}
	if includeFilters {
		body["includeFilters"] = true
	}
	req.WithBody(body)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseListError(err, listID)
	}

	var listResp ListResponse
	if err := json.Unmarshal(resp.Body, &listResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal list response: %w", err)
	}

	return &listResp.List, nil
}

func (c *Client) DeleteList(ctx context.Context, listID string) error {
	req := client.NewRequest("DELETE", fmt.Sprintf("/crm/v3/lists/%s", listID))
	req.WithContext(ctx)
	req.WithResourceType("lists")

	_, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return ParseListError(err, listID)
	}

	return nil
}

func (c *Client) RestoreList(ctx context.Context, listID string) error {
	req := client.NewRequest("PUT", fmt.Sprintf("/crm/v3/lists/%s/restore", listID))
	req.WithContext(ctx)
	req.WithResourceType("lists")

	_, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return ParseListError(err, listID)
	}

	return nil
}

func (c *Client) GetRecordMemberships(ctx context.Context, objectTypeID, recordID string) (*RecordMembershipsResponse, error) {
	req := client.NewRequest("GET", fmt.Sprintf("/crm/v3/lists/records/%s/%s/memberships", objectTypeID, recordID))
	req.WithContext(ctx)
	req.WithResourceType("lists")

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseRecordError(err, recordID, "")
	}

	var memberships RecordMembershipsResponse
	if err := json.Unmarshal(resp.Body, &memberships); err != nil {
		return nil, fmt.Errorf("failed to unmarshal memberships response: %w", err)
	}

	return &memberships, nil
}

func (c *Client) BatchGetRecordMemberships(ctx context.Context, inputs []MembershipRecordIdentifier) (*BatchReadMembershipsResponse, error) {
	req := client.NewRequest("POST", "/crm/v3/lists/records/memberships/batch/read")
	req.WithContext(ctx)
	req.WithResourceType("lists")
	req.WithBody(BatchReadMembershipsRequest{Inputs: inputs})

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseListError(err, "")
	}

	var batchResp BatchReadMembershipsResponse
	if err := json.Unmarshal(resp.Body, &batchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal batch memberships response: %w", err)
	}

	return &batchResp, nil
}

func (c *Client) AddRecordsToList(ctx context.Context, listID string, recordIDs []string) (*MembershipChangeResponse, error) {
	req := client.NewRequest("PUT", fmt.Sprintf("/crm/v3/lists/%s/memberships/add", listID))
	req.WithContext(ctx)
	req.WithResourceType("lists")
	req.WithBody(recordIDs)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseListError(err, listID)
	}

	var changeResp MembershipChangeResponse
	if err := json.Unmarshal(resp.Body, &changeResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal membership change response: %w", err)
	}

	return &changeResp, nil
}

func (c *Client) AddFromSourceList(ctx context.Context, listID, sourceListID string) (*MembershipChangeResponse, error) {
	req := client.NewRequest("PUT", fmt.Sprintf("/crm/v3/lists/%s/memberships/add-from/%s", listID, sourceListID))
	req.WithContext(ctx)
	req.WithResourceType("lists")

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseListError(err, listID)
	}

	var changeResp MembershipChangeResponse
	if err := json.Unmarshal(resp.Body, &changeResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal membership change response: %w", err)
	}

	return &changeResp, nil
}

func (c *Client) GetListMemberships(ctx context.Context, listID string, opts ...ListMembershipsOption) (*ListMembershipsResponse, error) {
	req := client.NewRequest("GET", fmt.Sprintf("/crm/v3/lists/%s/memberships", listID))
	req.WithContext(ctx)
	req.WithResourceType("lists")

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseListError(err, listID)
	}

	var memberships ListMembershipsResponse
	if err := json.Unmarshal(resp.Body, &memberships); err != nil {
		return nil, fmt.Errorf("failed to unmarshal list memberships response: %w", err)
	}

	return &memberships, nil
}

func (c *Client) RemoveAllRecords(ctx context.Context, listID string) error {
	req := client.NewRequest("DELETE", fmt.Sprintf("/crm/v3/lists/%s/memberships", listID))
	req.WithContext(ctx)
	req.WithResourceType("lists")

	_, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return ParseListError(err, listID)
	}

	return nil
}

func (c *Client) RemoveRecordsFromList(ctx context.Context, listID string, recordIDs []string) (*MembershipChangeResponse, error) {
	req := client.NewRequest("PUT", fmt.Sprintf("/crm/v3/lists/%s/memberships/remove", listID))
	req.WithContext(ctx)
	req.WithResourceType("lists")
	req.WithBody(recordIDs)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseListError(err, listID)
	}

	var changeResp MembershipChangeResponse
	if err := json.Unmarshal(resp.Body, &changeResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal membership change response: %w", err)
	}

	return &changeResp, nil
}

func (c *Client) ScheduleConversion(ctx context.Context, listID string, conversionReq *ScheduleConversionRequest) (*ScheduleConversionResponse, error) {
	req := client.NewRequest("PUT", fmt.Sprintf("/crm/v3/lists/%s/schedule-conversion", listID))
	req.WithContext(ctx)
	req.WithResourceType("lists")
	req.WithBody(conversionReq)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseListError(err, listID)
	}

	var conversionResp ScheduleConversionResponse
	if err := json.Unmarshal(resp.Body, &conversionResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal conversion schedule response: %w", err)
	}

	return &conversionResp, nil
}

func (c *Client) GetConversionSchedule(ctx context.Context, listID string) (*ScheduleConversionResponse, error) {
	req := client.NewRequest("GET", fmt.Sprintf("/crm/v3/lists/%s/schedule-conversion", listID))
	req.WithContext(ctx)
	req.WithResourceType("lists")

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseListError(err, listID)
	}

	var conversionResp ScheduleConversionResponse
	if err := json.Unmarshal(resp.Body, &conversionResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal conversion schedule response: %w", err)
	}

	return &conversionResp, nil
}

func (c *Client) DeleteConversionSchedule(ctx context.Context, listID string) error {
	req := client.NewRequest("DELETE", fmt.Sprintf("/crm/v3/lists/%s/schedule-conversion", listID))
	req.WithContext(ctx)
	req.WithResourceType("lists")

	_, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return ParseListError(err, listID)
	}

	return nil
}
