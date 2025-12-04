// Package companies provides client methods for the HubSpot CRM Companies API
package companies

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/josiah-hester/go-hubspot-sdk/client"
)

// Client represents the Companies API client
type Client struct {
	apiClient *client.Client
}

// NewClient creates a new companies client
func NewClient(apiClient *client.Client) *Client {
	return &Client{
		apiClient: apiClient,
	}
}

// CreateCompany creates a new company
func (c *Client) CreateCompany(ctx context.Context, input *CreateCompanyInput) (*Company, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/companies")
	req.WithContext(ctx)
	req.WithResourceType("companies")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var company Company
	if err := json.Unmarshal(resp.Body, &company); err != nil {
		return nil, fmt.Errorf("failed to unmarshal company response: %w", err)
	}

	return &company, nil
}

// GetCompany retrieves a company by ID
func (c *Client) GetCompany(ctx context.Context, companyID string, opts ...CompanyOption) (*Company, error) {
	req := client.NewRequest("GET", fmt.Sprintf("/crm/v3/objects/companies/%s", companyID))
	req.WithContext(ctx)
	req.WithResourceType("companies")

	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var company Company
	if err := json.Unmarshal(resp.Body, &company); err != nil {
		return nil, fmt.Errorf("failed to unmarshal company response: %w", err)
	}

	return &company, nil
}

// UpdateCompany updates a company
func (c *Client) UpdateCompany(ctx context.Context, companyID string, input *UpdateCompanyInput) (*Company, error) {
	req := client.NewRequest("PATCH", fmt.Sprintf("/crm/v3/objects/companies/%s", companyID))
	req.WithContext(ctx)
	req.WithResourceType("companies")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var company Company
	if err := json.Unmarshal(resp.Body, &company); err != nil {
		return nil, fmt.Errorf("failed to unmarshal company response: %w", err)
	}

	return &company, nil
}

// ArchiveCompany archives (deletes) a company
func (c *Client) ArchiveCompany(ctx context.Context, companyID string) error {
	req := client.NewRequest("DELETE", fmt.Sprintf("/crm/v3/objects/companies/%s", companyID))
	req.WithContext(ctx)
	req.WithResourceType("companies")

	_, err := c.apiClient.Do(ctx, req)
	return err
}

// ListCompanies lists companies with optional filters
func (c *Client) ListCompanies(ctx context.Context, opts ...CompanyOption) (*ListCompaniesResponse, error) {
	req := client.NewRequest("GET", "/crm/v3/objects/companies")
	req.WithContext(ctx)
	req.WithResourceType("companies")

	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var listResp ListCompaniesResponse
	if err := json.Unmarshal(resp.Body, &listResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal companies list response: %w", err)
	}

	return &listResp, nil
}

// BatchReadCompanies retrieves multiple companies by ID
func (c *Client) BatchReadCompanies(ctx context.Context, input *BatchReadCompaniesInput) (*BatchCompaniesResponse, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/companies/batch/read")
	req.WithContext(ctx)
	req.WithResourceType("companies")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var batchResp BatchCompaniesResponse
	if err := json.Unmarshal(resp.Body, &batchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal batch response: %w", err)
	}

	return &batchResp, nil
}

// BatchCreateCompanies creates multiple companies
func (c *Client) BatchCreateCompanies(ctx context.Context, input *BatchCreateCompaniesInput) (*BatchCompaniesResponse, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/companies/batch/create")
	req.WithContext(ctx)
	req.WithResourceType("companies")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var batchResp BatchCompaniesResponse
	if err := json.Unmarshal(resp.Body, &batchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal batch response: %w", err)
	}

	return &batchResp, nil
}

// BatchUpdateCompanies updates multiple companies
func (c *Client) BatchUpdateCompanies(ctx context.Context, input *BatchUpdateCompaniesInput) (*BatchCompaniesResponse, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/companies/batch/update")
	req.WithContext(ctx)
	req.WithResourceType("companies")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var batchResp BatchCompaniesResponse
	if err := json.Unmarshal(resp.Body, &batchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal batch response: %w", err)
	}

	return &batchResp, nil
}

// BatchArchiveCompanies archives multiple companies
func (c *Client) BatchArchiveCompanies(ctx context.Context, input *BatchArchiveCompaniesInput) error {
	req := client.NewRequest("POST", "/crm/v3/objects/companies/batch/archive")
	req.WithContext(ctx)
	req.WithResourceType("companies")
	req.WithBody(input)

	_, err := c.apiClient.Do(ctx, req)
	return err
}

// SearchCompanies searches for companies
func (c *Client) SearchCompanies(ctx context.Context, input *SearchCompaniesInput) (*SearchCompaniesResponse, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/companies/search")
	req.WithContext(ctx)
	req.WithResourceType("companies")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var searchResp SearchCompaniesResponse
	if err := json.Unmarshal(resp.Body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %w", err)
	}

	return &searchResp, nil
}
