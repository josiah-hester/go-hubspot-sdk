// Package deals provides client methods for the HubSpot CRM Deals API
package deals

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/josiah-hester/go-hubspot-sdk/client"
)

// Client represents the Deals API client
type Client struct {
	apiClient *client.Client
}

// NewClient creates a new deals client
func NewClient(apiClient *client.Client) *Client {
	return &Client{
		apiClient: apiClient,
	}
}

// CreateDeal creates a new deal
func (c *Client) CreateDeal(ctx context.Context, input *CreateDealInput) (*Deal, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/deals")
	req.WithContext(ctx)
	req.WithResourceType("deals")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var deal Deal
	if err := json.Unmarshal(resp.Body, &deal); err != nil {
		return nil, fmt.Errorf("failed to unmarshal deal response: %w", err)
	}

	return &deal, nil
}

// GetDeal retrieves a deal by ID
func (c *Client) GetDeal(ctx context.Context, dealID string, opts ...DealOption) (*Deal, error) {
	req := client.NewRequest("GET", fmt.Sprintf("/crm/v3/objects/deals/%s", dealID))
	req.WithContext(ctx)
	req.WithResourceType("deals")

	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var deal Deal
	if err := json.Unmarshal(resp.Body, &deal); err != nil {
		return nil, fmt.Errorf("failed to unmarshal deal response: %w", err)
	}

	return &deal, nil
}

// UpdateDeal updates a deal
func (c *Client) UpdateDeal(ctx context.Context, dealID string, input *UpdateDealInput) (*Deal, error) {
	req := client.NewRequest("PATCH", fmt.Sprintf("/crm/v3/objects/deals/%s", dealID))
	req.WithContext(ctx)
	req.WithResourceType("deals")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var deal Deal
	if err := json.Unmarshal(resp.Body, &deal); err != nil {
		return nil, fmt.Errorf("failed to unmarshal deal response: %w", err)
	}

	return &deal, nil
}

// ArchiveDeal archives (deletes) a deal
func (c *Client) ArchiveDeal(ctx context.Context, dealID string) error {
	req := client.NewRequest("DELETE", fmt.Sprintf("/crm/v3/objects/deals/%s", dealID))
	req.WithContext(ctx)
	req.WithResourceType("deals")

	_, err := c.apiClient.Do(ctx, req)
	return err
}

// ListDeals lists deals with optional filters
func (c *Client) ListDeals(ctx context.Context, opts ...DealOption) (*ListDealsResponse, error) {
	req := client.NewRequest("GET", "/crm/v3/objects/deals")
	req.WithContext(ctx)
	req.WithResourceType("deals")

	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var listResp ListDealsResponse
	if err := json.Unmarshal(resp.Body, &listResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal deals list response: %w", err)
	}

	return &listResp, nil
}

// BatchReadDeals retrieves multiple deals by ID
func (c *Client) BatchReadDeals(ctx context.Context, input *BatchReadDealsInput) (*BatchDealsResponse, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/deals/batch/read")
	req.WithContext(ctx)
	req.WithResourceType("deals")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var batchResp BatchDealsResponse
	if err := json.Unmarshal(resp.Body, &batchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal batch response: %w", err)
	}

	return &batchResp, nil
}

// BatchCreateDeals creates multiple deals
func (c *Client) BatchCreateDeals(ctx context.Context, input *BatchCreateDealsInput) (*BatchDealsResponse, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/deals/batch/create")
	req.WithContext(ctx)
	req.WithResourceType("deals")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var batchResp BatchDealsResponse
	if err := json.Unmarshal(resp.Body, &batchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal batch response: %w", err)
	}

	return &batchResp, nil
}

// BatchUpdateDeals updates multiple deals
func (c *Client) BatchUpdateDeals(ctx context.Context, input *BatchUpdateDealsInput) (*BatchDealsResponse, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/deals/batch/update")
	req.WithContext(ctx)
	req.WithResourceType("deals")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var batchResp BatchDealsResponse
	if err := json.Unmarshal(resp.Body, &batchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal batch response: %w", err)
	}

	return &batchResp, nil
}

// BatchArchiveDeals archives multiple deals
func (c *Client) BatchArchiveDeals(ctx context.Context, input *BatchArchiveDealsInput) error {
	req := client.NewRequest("POST", "/crm/v3/objects/deals/batch/archive")
	req.WithContext(ctx)
	req.WithResourceType("deals")
	req.WithBody(input)

	_, err := c.apiClient.Do(ctx, req)
	return err
}

// SearchDeals searches for deals
func (c *Client) SearchDeals(ctx context.Context, input *SearchDealsInput) (*SearchDealsResponse, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/deals/search")
	req.WithContext(ctx)
	req.WithResourceType("deals")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var searchResp SearchDealsResponse
	if err := json.Unmarshal(resp.Body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %w", err)
	}

	return &searchResp, nil
}
