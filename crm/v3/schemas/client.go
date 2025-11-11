// Package schemas specifies the client methods for the HubSpot CRM Schemas API
package schemas

import (
	"context"
	"fmt"

	"github.com/josiah-hester/go-hubspot-sdk/client"
	"github.com/josiah-hester/go-hubspot-sdk/internal/tools"
)

type Client struct {
	apiClient *client.Client
}

// NewClient creates a new schemas client
func NewClient(apiClient *client.Client) *Client {
	return &Client{
		apiClient: apiClient,
	}
}

// -------- Basic Methods --------

// GetAllSchemas gets all schemas
//
// opts:
// WithArchived
func (c *Client) GetAllSchemas(ctx context.Context, opts ...SchemaOption) (*GetAllSchemasResponse, error) {
	req := client.NewRequest("GET", "/crm-object-schemas/v3/schemas")
	req.WithContext(ctx)
	req.WithResourceType("schemas")

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var schemas GetAllSchemasResponse
	if err := tools.NewRequiredTagStruct(schemas).UnmarhsalJSON(resp.Body); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schemas response: %w", err)
	}

	if len(schemas.Results) == 0 {
		return &schemas, fmt.Errorf("no schemas found")
	}

	return &schemas, nil
}

// GetExistingSchema gets an existing schema by object type
func (c *Client) GetExistingSchema(ctx context.Context, objectType string) (*Schema, error) {
	req := client.NewRequest("GET", fmt.Sprintf("/crm-object-schemas/v3/schemas/%s", objectType))
	req.WithContext(ctx)
	req.WithResourceType("schemas")

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var schema Schema
	if err := tools.NewRequiredTagStruct(schema).UnmarhsalJSON(resp.Body); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema response: %w", err)
	}

	return &schema, nil
}

// CreateNewSchema creates a new object schema
func (c *Client) CreateNewSchema(ctx context.Context, input *CreateNewSchemaInput) (*Schema, error) {
	req := client.NewRequest("POST", "/crm-object-schemas/v3/schemas")
	req.WithContext(ctx)
	req.WithResourceType("schemas")

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var schema Schema
	if err := tools.NewRequiredTagStruct(schema).UnmarhsalJSON(resp.Body); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema response: %w", err)
	}

	return &schema, nil
}

// CreateNewAssociationSchema creates a new object assocation
func (c *Client) CreateNewAssociationSchema(ctx context.Context, objectType string, input *CreateNewAssociationSchemaInput) (*CreateNewAssociationSchemaResponse, error) {
	req := client.NewRequest("POST", fmt.Sprintf("/crm-object-schemas/v3/schemas/%s/associations", objectType))
	req.WithContext(ctx)
	req.WithResourceType("schemas")

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var assocResp CreateNewAssociationSchemaResponse
	if err := tools.NewRequiredTagStruct(assocResp).UnmarhsalJSON(resp.Body); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema response: %w", err)
	}

	return &assocResp, nil
}

// UpdateSchema udpates an existing object schema
func (c *Client) UpdateSchema(ctx context.Context, objectType string, input *UpdateSchemaInput) (*Schema, error) {
	req := client.NewRequest("PATCH", fmt.Sprintf("/crm-object-schemas/v3/schemas/%s", objectType))
	req.WithContext(ctx)
	req.WithResourceType("schemas")

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var schema Schema
	if err := tools.NewRequiredTagStruct(schema).UnmarhsalJSON(resp.Body); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema response: %w", err)
	}

	return &schema, nil
}

// DeleteSchema deletes an existing object schema
func (c *Client) DeleteSchema(ctx context.Context, objectType string, opts ...SchemaOption) error {
	req := client.NewRequest("DELETE", fmt.Sprintf("/crm-object-schemas/v3/schemas/%s", objectType))
	req.WithContext(ctx)
	req.WithResourceType("schemas")

	for _, opt := range opts {
		opt(req)
	}

	_, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

// RemoveAssociationSchema removes an association schema
func (c *Client) RemoveAssociationSchema(ctx context.Context, objectType, associationIdentifier string) error {
	req := client.NewRequest("DELETE", fmt.Sprintf("/crm-object-schemas/v3/schemas/%s/associations/%s", objectType, associationIdentifier))
	req.WithContext(ctx)
	req.WithResourceType("schemas")

	_, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
