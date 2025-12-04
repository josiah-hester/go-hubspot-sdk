// Package associations provides client methods for the HubSpot CRM Associations v4 API
package associations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/josiah-hester/go-hubspot-sdk/client"
)

// Client represents the Associations API client
type Client struct {
	apiClient *client.Client
}

// NewClient creates a new associations client
func NewClient(apiClient *client.Client) *Client {
	return &Client{
		apiClient: apiClient,
	}
}

// CreateAssociation creates an association between two objects and returns the association details
func (c *Client) CreateAssociation(ctx context.Context, fromObjectType, fromObjectID, toObjectType, toObjectID string, associationSpecs []AssociationSpec) (*AssociationResponse, error) {
	req := client.NewRequest("PUT", fmt.Sprintf("/crm/v4/objects/%s/%s/associations/%s/%s",
		fromObjectType, fromObjectID, toObjectType, toObjectID))
	req.WithContext(ctx)
	req.WithResourceType("associations")
	req.WithBody(associationSpecs)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var assocResp AssociationResponse
	if err := json.Unmarshal(resp.Body, &assocResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal association response: %w", err)
	}

	return &assocResp, nil
}

// DeleteAssociation removes an association between two objects
func (c *Client) DeleteAssociation(ctx context.Context, fromObjectType, fromObjectID, toObjectType, toObjectID string, associationSpecs []AssociationSpec) error {
	req := client.NewRequest("DELETE", fmt.Sprintf("/crm/v4/objects/%s/%s/associations/%s/%s",
		fromObjectType, fromObjectID, toObjectType, toObjectID))
	req.WithContext(ctx)
	req.WithResourceType("associations")
	req.WithBody(associationSpecs)

	_, err := c.apiClient.Do(ctx, req)
	return err
}

// ListAssociations retrieves all associations for an object
func (c *Client) ListAssociations(ctx context.Context, fromObjectType, fromObjectID, toObjectType string, opts ...AssociationOption) (*ListAssociationsResponse, error) {
	req := client.NewRequest("GET", fmt.Sprintf("/crm/v4/objects/%s/%s/associations/%s",
		fromObjectType, fromObjectID, toObjectType))
	req.WithContext(ctx)
	req.WithResourceType("associations")

	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var listResp ListAssociationsResponse
	if err := json.Unmarshal(resp.Body, &listResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal associations response: %w", err)
	}

	return &listResp, nil
}

// BatchCreateAssociations creates multiple associations
func (c *Client) BatchCreateAssociations(ctx context.Context, fromObjectType, toObjectType string, input *BatchAssociationInput) error {
	req := client.NewRequest("POST", fmt.Sprintf("/crm/v4/associations/%s/%s/batch/create",
		fromObjectType, toObjectType))
	req.WithContext(ctx)
	req.WithResourceType("associations")
	req.WithBody(input)

	_, err := c.apiClient.Do(ctx, req)
	return err
}

// BatchDeleteAssociations removes multiple associations
func (c *Client) BatchDeleteAssociations(ctx context.Context, fromObjectType, toObjectType string, input *BatchAssociationInput) error {
	req := client.NewRequest("POST", fmt.Sprintf("/crm/v4/associations/%s/%s/batch/archive",
		fromObjectType, toObjectType))
	req.WithContext(ctx)
	req.WithResourceType("associations")
	req.WithBody(input)

	_, err := c.apiClient.Do(ctx, req)
	return err
}

// GetAssociationLabels retrieves all association labels between two object types
func (c *Client) GetAssociationLabels(ctx context.Context, fromObjectType, toObjectType string) (*GetAssociationLabelsResponse, error) {
	req := client.NewRequest("GET", fmt.Sprintf("/crm/v4/associations/%s/%s/labels",
		fromObjectType, toObjectType))
	req.WithContext(ctx)
	req.WithResourceType("associations")

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var labelsResp GetAssociationLabelsResponse
	if err := json.Unmarshal(resp.Body, &labelsResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal labels response: %w", err)
	}

	return &labelsResp, nil
}
