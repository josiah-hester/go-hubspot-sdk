package deals

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

	dealsClient := NewClient(apiClient)
	return server, dealsClient
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

	dealsClient := NewClient(apiClient)
	assert.NotNil(t, dealsClient)
	assert.NotNil(t, dealsClient.apiClient)
}

// TestCreateDeal tests deal creation
func TestCreateDeal_Success(t *testing.T) {
	responseJSON := `{
		"id": "123456",
		"properties": {
			"dealname": "Q4 Enterprise Deal",
			"amount": "50000",
			"dealstage": "qualifiedtobuy"
		},
		"createdAt": "2024-01-01T00:00:00.000Z",
		"updatedAt": "2024-01-01T00:00:00.000Z",
		"archived": false
	}`

	server, dealsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/deals", r.URL.Path)
		respondJSON(w, http.StatusCreated, responseJSON)
	})
	defer server.Close()

	input := &CreateDealInput{
		Properties: map[string]string{
			"dealname":  "Q4 Enterprise Deal",
			"amount":    "50000",
			"dealstage": "qualifiedtobuy",
		},
	}

	deal, err := dealsClient.CreateDeal(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, deal)
	assert.Equal(t, "123456", deal.ID)
	assert.Equal(t, "Q4 Enterprise Deal", deal.Properties["dealname"])
	assert.Equal(t, "50000", deal.Properties["amount"])
}

func TestCreateDeal_Error(t *testing.T) {
	server, dealsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusBadRequest, `{"status": "error", "message": "Invalid input"}`)
	})
	defer server.Close()

	input := &CreateDealInput{Properties: map[string]string{}}
	_, err := dealsClient.CreateDeal(context.Background(), input)

	require.Error(t, err)
}

// TestGetDeal tests retrieving a deal
func TestGetDeal_Success(t *testing.T) {
	responseJSON := `{
		"id": "123456",
		"properties": {
			"dealname": "Enterprise Deal",
			"amount": "75000"
		},
		"createdAt": "2024-01-01T00:00:00.000Z",
		"updatedAt": "2024-01-01T00:00:00.000Z",
		"archived": false
	}`

	server, dealsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/crm/v3/objects/deals/123456", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	deal, err := dealsClient.GetDeal(context.Background(), "123456")

	require.NoError(t, err)
	assert.NotNil(t, deal)
	assert.Equal(t, "123456", deal.ID)
	assert.Equal(t, "Enterprise Deal", deal.Properties["dealname"])
}

func TestGetDeal_WithOptions(t *testing.T) {
	responseJSON := `{
		"id": "123456",
		"properties": {
			"dealname": "Test Deal"
		},
		"createdAt": "2024-01-01T00:00:00.000Z",
		"updatedAt": "2024-01-01T00:00:00.000Z",
		"archived": false
	}`

	server, dealsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "dealname,amount,dealstage", r.URL.Query().Get("properties"))
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	deal, err := dealsClient.GetDeal(context.Background(), "123456",
		WithProperties([]string{"dealname", "amount", "dealstage"}))

	require.NoError(t, err)
	assert.NotNil(t, deal)
}

func TestGetDeal_NotFound(t *testing.T) {
	server, dealsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, `{"status": "error", "message": "Not found"}`)
	})
	defer server.Close()

	_, err := dealsClient.GetDeal(context.Background(), "999999")
	require.Error(t, err)
}

// TestUpdateDeal tests updating a deal
func TestUpdateDeal_Success(t *testing.T) {
	responseJSON := `{
		"id": "123456",
		"properties": {
			"dealname": "Updated Deal Name",
			"amount": "100000",
			"dealstage": "closedwon"
		},
		"createdAt": "2024-01-01T00:00:00.000Z",
		"updatedAt": "2024-01-02T00:00:00.000Z",
		"archived": false
	}`

	server, dealsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method)
		assert.Equal(t, "/crm/v3/objects/deals/123456", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &UpdateDealInput{
		Properties: map[string]string{
			"dealname":  "Updated Deal Name",
			"dealstage": "closedwon",
		},
	}

	deal, err := dealsClient.UpdateDeal(context.Background(), "123456", input)

	require.NoError(t, err)
	assert.NotNil(t, deal)
	assert.Equal(t, "Updated Deal Name", deal.Properties["dealname"])
	assert.Equal(t, "closedwon", deal.Properties["dealstage"])
}

// TestArchiveDeal tests archiving a deal
func TestArchiveDeal_Success(t *testing.T) {
	server, dealsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/crm/v3/objects/deals/123456", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := dealsClient.ArchiveDeal(context.Background(), "123456")
	require.NoError(t, err)
}

// TestListDeals tests listing deals
func TestListDeals_Success(t *testing.T) {
	responseJSON := `{
		"results": [
			{
				"id": "1",
				"properties": {"dealname": "Deal A", "amount": "10000"},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			},
			{
				"id": "2",
				"properties": {"dealname": "Deal B", "amount": "25000"},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}
		],
		"paging": {
			"next": {
				"after": "abc123",
				"link": "?after=abc123"
			}
		}
	}`

	server, dealsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/crm/v3/objects/deals", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	resp, err := dealsClient.ListDeals(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Results, 2)
	assert.Equal(t, "abc123", resp.Paging.Next.After)
}

func TestListDeals_WithOptions(t *testing.T) {
	responseJSON := `{"results": [], "paging": null}`

	server, dealsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "20", r.URL.Query().Get("limit"))
		assert.Equal(t, "xyz789", r.URL.Query().Get("after"))
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	_, err := dealsClient.ListDeals(context.Background(),
		WithLimit(20),
		WithAfter("xyz789"))

	require.NoError(t, err)
}

// TestBatchReadDeals tests batch read
func TestBatchReadDeals_Success(t *testing.T) {
	responseJSON := `{
		"status": "COMPLETE",
		"results": [
			{
				"id": "1",
				"properties": {"dealname": "Deal A"},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}
		],
		"startedAt": "2024-01-01T00:00:00.000Z",
		"completedAt": "2024-01-01T00:00:05.000Z"
	}`

	server, dealsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/deals/batch/read", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchReadDealsInput{
		Properties: []string{"dealname", "amount"},
		Inputs: []struct {
			ID string `json:"id"`
		}{
			{ID: "1"},
		},
	}

	resp, err := dealsClient.BatchReadDeals(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "COMPLETE", resp.Status)
	assert.Len(t, resp.Results, 1)
}

// TestBatchCreateDeals tests batch create
func TestBatchCreateDeals_Success(t *testing.T) {
	responseJSON := `{
		"status": "COMPLETE",
		"results": [
			{
				"id": "1",
				"properties": {"dealname": "New Deal"},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}
		],
		"startedAt": "2024-01-01T00:00:00.000Z",
		"completedAt": "2024-01-01T00:00:05.000Z"
	}`

	server, dealsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/deals/batch/create", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchCreateDealsInput{
		Inputs: []CreateDealInput{
			{Properties: map[string]string{"dealname": "New Deal"}},
		},
	}

	resp, err := dealsClient.BatchCreateDeals(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "COMPLETE", resp.Status)
}

// TestBatchUpdateDeals tests batch update
func TestBatchUpdateDeals_Success(t *testing.T) {
	responseJSON := `{
		"status": "COMPLETE",
		"results": [],
		"startedAt": "2024-01-01T00:00:00.000Z",
		"completedAt": "2024-01-01T00:00:05.000Z"
	}`

	server, dealsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/deals/batch/update", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchUpdateDealsInput{}
	resp, err := dealsClient.BatchUpdateDeals(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

// TestBatchArchiveDeals tests batch archive
func TestBatchArchiveDeals_Success(t *testing.T) {
	server, dealsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/deals/batch/archive", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	input := &BatchArchiveDealsInput{}
	err := dealsClient.BatchArchiveDeals(context.Background(), input)

	require.NoError(t, err)
}

// TestSearchDeals tests search functionality
func TestSearchDeals_Success(t *testing.T) {
	responseJSON := `{
		"total": 1,
		"results": [
			{
				"id": "1",
				"properties": {"dealname": "Enterprise Deal", "amount": "100000"},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}
		],
		"paging": null
	}`

	server, dealsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/deals/search", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &SearchDealsInput{
		FilterGroups: []FilterGroup{
			{
				Filters: []Filter{
					{
						PropertyName: "amount",
						Operator:     "GTE",
						Value:        "50000",
					},
				},
			},
		},
		Properties: []string{"dealname", "amount"},
		Limit:      10,
	}

	resp, err := dealsClient.SearchDeals(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.Total)
	assert.Len(t, resp.Results, 1)
}

// TestOptions tests all option functions
func TestOptions(t *testing.T) {
	tests := []struct {
		name     string
		option   DealOption
		expected string
		param    string
	}{
		{"WithPropertiesWithHistory", WithPropertiesWithHistory([]string{"dealname", "amount"}), "dealname,amount", "propertiesWithHistory"},
		{"WithAssociations", WithAssociations([]string{"contacts", "companies"}), "contacts,companies", "associations"},
		{"WithIDProperty", WithIDProperty("dealid"), "dealid", "idProperty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, dealsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tt.expected, r.URL.Query().Get(tt.param))
				respondJSON(w, http.StatusOK, `{"id": "123", "properties": {"dealname": "Test"}, "createdAt": "2024-01-01T00:00:00.000Z", "updatedAt": "2024-01-01T00:00:00.000Z", "archived": false}`)
			})
			defer server.Close()

			_, err := dealsClient.GetDeal(context.Background(), "123", tt.option)
			require.NoError(t, err)
		})
	}

	t.Run("WithArchived", func(t *testing.T) {
		server, dealsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "true", r.URL.Query().Get("archived"))
			respondJSON(w, http.StatusOK, `{"results": [], "paging": null}`)
		})
		defer server.Close()

		_, err := dealsClient.ListDeals(context.Background(), WithArchived())
		require.NoError(t, err)
	})
}
