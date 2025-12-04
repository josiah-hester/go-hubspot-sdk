package associations

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/josiah-hester/go-hubspot-sdk/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper functions
func setupMockServer(t *testing.T, handler func(w http.ResponseWriter, r *http.Request)) (*httptest.Server, *Client) {
	server := httptest.NewServer(http.HandlerFunc(handler))

	apiClient, err := client.NewClient(
		client.WithBaseURL(server.URL),
		client.WithAccessToken("test-token"),
		client.WithRateLimitEnabled(false),
		client.WithRetryEnabled(false),
	)
	require.NoError(t, err)

	assocClient := NewClient(apiClient)
	return server, assocClient
}

func respondJSON(w http.ResponseWriter, statusCode int, body string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(body))
}

// TestNewClient tests client creation
func TestNewClient(t *testing.T) {
	apiClient, err := client.NewClient()
	require.NoError(t, err)

	assocClient := NewClient(apiClient)
	assert.NotNil(t, assocClient)
	assert.NotNil(t, assocClient.apiClient)
}

// TestCreateAssociation tests creating an association
func TestCreateAssociation_Success(t *testing.T) {
	server, assocClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/crm/v4/objects/contacts/123/associations/companies/456", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "success"}`))
	})
	defer server.Close()

	err := assocClient.CreateAssociation(context.Background(),
		"contacts", "123",
		"companies", "456",
		1)

	require.NoError(t, err)
}

func TestCreateAssociation_Error(t *testing.T) {
	server, assocClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusBadRequest, `{"status": "error", "message": "Invalid association"}`)
	})
	defer server.Close()

	err := assocClient.CreateAssociation(context.Background(),
		"contacts", "123",
		"companies", "999",
		1)

	require.Error(t, err)
}

// TestDeleteAssociation tests deleting an association
func TestDeleteAssociation_Success(t *testing.T) {
	server, assocClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/crm/v4/objects/contacts/123/associations/companies/456", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := assocClient.DeleteAssociation(context.Background(),
		"contacts", "123",
		"companies", "456",
		1)

	require.NoError(t, err)
}

// TestListAssociations tests listing associations
func TestListAssociations_Success(t *testing.T) {
	responseJSON := `{
		"results": [
			{
				"toObjectId": "456",
				"associationTypes": [
					{
						"associationCategory": "HUBSPOT_DEFINED",
						"associationTypeId": 1
					}
				]
			},
			{
				"toObjectId": "789",
				"associationTypes": [
					{
						"associationCategory": "HUBSPOT_DEFINED",
						"associationTypeId": 1
					}
				]
			}
		],
		"paging": {
			"next": {
				"after": "abc123",
				"link": "?after=abc123"
			}
		}
	}`

	server, assocClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/crm/v4/objects/contacts/123/associations/companies", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	resp, err := assocClient.ListAssociations(context.Background(),
		"contacts", "123",
		"companies")

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Results, 2)
	assert.Equal(t, "456", resp.Results[0].ToObjectID)
	assert.Equal(t, "abc123", resp.Paging.Next.After)
}

func TestListAssociations_WithOptions(t *testing.T) {
	responseJSON := `{"results": [], "paging": null}`

	server, assocClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "10", r.URL.Query().Get("limit"))
		assert.Equal(t, "xyz789", r.URL.Query().Get("after"))
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	_, err := assocClient.ListAssociations(context.Background(),
		"contacts", "123",
		"companies",
		WithLimit(10),
		WithAfter("xyz789"))

	require.NoError(t, err)
}

func TestListAssociations_Empty(t *testing.T) {
	responseJSON := `{"results": [], "paging": null}`

	server, assocClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	resp, err := assocClient.ListAssociations(context.Background(),
		"contacts", "999",
		"companies")

	require.NoError(t, err)
	assert.Len(t, resp.Results, 0)
}

// TestBatchCreateAssociations tests batch create
func TestBatchCreateAssociations_Success(t *testing.T) {
	server, assocClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v4/associations/contacts/companies/batch/create", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "success"}`))
	})
	defer server.Close()

	input := &BatchAssociationInput{
		Inputs: []struct {
			From struct {
				ID string `json:"id"`
			} `json:"from"`
			To []struct {
				ID    string            `json:"id"`
				Types []AssociationSpec `json:"types"`
			} `json:"to"`
		}{
			{
				From: struct {
					ID string `json:"id"`
				}{ID: "123"},
				To: []struct {
					ID    string            `json:"id"`
					Types []AssociationSpec `json:"types"`
				}{
					{
						ID: "456",
						Types: []AssociationSpec{
							{
								AssociationCategory: "HUBSPOT_DEFINED",
								AssociationTypeID:   1,
							},
						},
					},
				},
			},
		},
	}

	err := assocClient.BatchCreateAssociations(context.Background(),
		"contacts", "companies",
		input)

	require.NoError(t, err)
}

// TestBatchDeleteAssociations tests batch delete
func TestBatchDeleteAssociations_Success(t *testing.T) {
	server, assocClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v4/associations/contacts/companies/batch/archive", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	input := &BatchAssociationInput{}
	err := assocClient.BatchDeleteAssociations(context.Background(),
		"contacts", "companies",
		input)

	require.NoError(t, err)
}

// TestGetAssociationLabels tests getting association labels
func TestGetAssociationLabels_Success(t *testing.T) {
	responseJSON := `{
		"results": [
			{
				"category": "HUBSPOT_DEFINED",
				"typeId": 1,
				"label": "Primary"
			},
			{
				"category": "HUBSPOT_DEFINED",
				"typeId": 2,
				"label": "Secondary"
			}
		]
	}`

	server, assocClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/crm/v4/associations/contacts/companies/labels", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	resp, err := assocClient.GetAssociationLabels(context.Background(),
		"contacts", "companies")

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Results, 2)
	assert.Equal(t, "Primary", resp.Results[0].Label)
	assert.Equal(t, 1, resp.Results[0].TypeID)
}

func TestGetAssociationLabels_Empty(t *testing.T) {
	responseJSON := `{"results": []}`

	server, assocClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	resp, err := assocClient.GetAssociationLabels(context.Background(),
		"customobject", "contacts")

	require.NoError(t, err)
	assert.Len(t, resp.Results, 0)
}

// TestAssociations_MultipleObjectTypes tests various object type combinations
func TestAssociations_MultipleObjectTypes(t *testing.T) {
	testCases := []struct {
		name           string
		fromObjectType string
		toObjectType   string
	}{
		{"Contact to Company", "contacts", "companies"},
		{"Contact to Deal", "contacts", "deals"},
		{"Company to Deal", "companies", "deals"},
		{"Deal to Order", "deals", "orders"},
		{"Contact to Order", "contacts", "orders"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server, assocClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"status": "success"}`))
			})
			defer server.Close()

			err := assocClient.CreateAssociation(context.Background(),
				tc.fromObjectType, "123",
				tc.toObjectType, "456",
				1)

			require.NoError(t, err)
		})
	}
}
