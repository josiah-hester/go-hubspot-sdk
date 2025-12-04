package companies

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

	companiesClient := NewClient(apiClient)
	return server, companiesClient
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

	companiesClient := NewClient(apiClient)
	assert.NotNil(t, companiesClient)
	assert.NotNil(t, companiesClient.apiClient)
}

// TestCreateCompany tests company creation
func TestCreateCompany_Success(t *testing.T) {
	responseJSON := `{
		"id": "123456",
		"properties": {
			"name": "Acme Corp",
			"domain": "acme.com"
		},
		"createdAt": "2024-01-01T00:00:00.000Z",
		"updatedAt": "2024-01-01T00:00:00.000Z",
		"archived": false
	}`

	server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/companies", r.URL.Path)
		respondJSON(w, http.StatusCreated, responseJSON)
	})
	defer server.Close()

	input := &CreateCompanyInput{
		Properties: map[string]string{
			"name":   "Acme Corp",
			"domain": "acme.com",
		},
	}

	company, err := companiesClient.CreateCompany(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, company)
	assert.Equal(t, "123456", company.ID)
	assert.Equal(t, "Acme Corp", company.Properties["name"])
}

func TestCreateCompany_Error(t *testing.T) {
	server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusBadRequest, `{"status": "error", "message": "Invalid input"}`)
	})
	defer server.Close()

	input := &CreateCompanyInput{Properties: map[string]string{}}
	_, err := companiesClient.CreateCompany(context.Background(), input)

	require.Error(t, err)
}

// TestGetCompany tests retrieving a company
func TestGetCompany_Success(t *testing.T) {
	responseJSON := `{
		"id": "123456",
		"properties": {
			"name": "Acme Corp",
			"domain": "acme.com"
		},
		"createdAt": "2024-01-01T00:00:00.000Z",
		"updatedAt": "2024-01-01T00:00:00.000Z",
		"archived": false
	}`

	server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/crm/v3/objects/companies/123456", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	company, err := companiesClient.GetCompany(context.Background(), "123456")

	require.NoError(t, err)
	assert.NotNil(t, company)
	assert.Equal(t, "123456", company.ID)
}

func TestGetCompany_WithOptions(t *testing.T) {
	responseJSON := `{
		"id": "123456",
		"properties": {
			"name": "Acme Corp"
		},
		"createdAt": "2024-01-01T00:00:00.000Z",
		"updatedAt": "2024-01-01T00:00:00.000Z",
		"archived": false
	}`

	server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "name,domain", r.URL.Query().Get("properties"))
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	company, err := companiesClient.GetCompany(context.Background(), "123456",
		WithProperties([]string{"name", "domain"}))

	require.NoError(t, err)
	assert.NotNil(t, company)
}

func TestGetCompany_NotFound(t *testing.T) {
	server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, `{"status": "error", "message": "Not found"}`)
	})
	defer server.Close()

	_, err := companiesClient.GetCompany(context.Background(), "999999")
	require.Error(t, err)
}

// TestUpdateCompany tests updating a company
func TestUpdateCompany_Success(t *testing.T) {
	responseJSON := `{
		"id": "123456",
		"properties": {
			"name": "Acme Corporation",
			"domain": "acme.com"
		},
		"createdAt": "2024-01-01T00:00:00.000Z",
		"updatedAt": "2024-01-02T00:00:00.000Z",
		"archived": false
	}`

	server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method)
		assert.Equal(t, "/crm/v3/objects/companies/123456", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &UpdateCompanyInput{
		Properties: map[string]string{
			"name": "Acme Corporation",
		},
	}

	company, err := companiesClient.UpdateCompany(context.Background(), "123456", input)

	require.NoError(t, err)
	assert.NotNil(t, company)
	assert.Equal(t, "Acme Corporation", company.Properties["name"])
}

// TestArchiveCompany tests archiving a company
func TestArchiveCompany_Success(t *testing.T) {
	server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/crm/v3/objects/companies/123456", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := companiesClient.ArchiveCompany(context.Background(), "123456")
	require.NoError(t, err)
}

// TestListCompanies tests listing companies
func TestListCompanies_Success(t *testing.T) {
	responseJSON := `{
		"results": [
			{
				"id": "1",
				"properties": {"name": "Company A"},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			},
			{
				"id": "2",
				"properties": {"name": "Company B"},
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

	server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/crm/v3/objects/companies", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	resp, err := companiesClient.ListCompanies(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Results, 2)
	assert.Equal(t, "abc123", resp.Paging.Next.After)
}

func TestListCompanies_WithOptions(t *testing.T) {
	responseJSON := `{"results": [], "paging": null}`

	server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "10", r.URL.Query().Get("limit"))
		assert.Equal(t, "xyz789", r.URL.Query().Get("after"))
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	_, err := companiesClient.ListCompanies(context.Background(),
		WithLimit(10),
		WithAfter("xyz789"))

	require.NoError(t, err)
}

// TestBatchReadCompanies tests batch read
func TestBatchReadCompanies_Success(t *testing.T) {
	responseJSON := `{
		"status": "COMPLETE",
		"results": [
			{
				"id": "1",
				"properties": {"name": "Company A"},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}
		],
		"startedAt": "2024-01-01T00:00:00.000Z",
		"completedAt": "2024-01-01T00:00:05.000Z"
	}`

	server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/companies/batch/read", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchReadCompaniesInput{
		Properties: []string{"name", "domain"},
		Inputs: []struct {
			ID string `json:"id"`
		}{
			{ID: "1"},
		},
	}

	resp, err := companiesClient.BatchReadCompanies(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "COMPLETE", resp.Status)
	assert.Len(t, resp.Results, 1)
}

// TestBatchCreateCompanies tests batch create
func TestBatchCreateCompanies_Success(t *testing.T) {
	responseJSON := `{
		"status": "COMPLETE",
		"results": [
			{
				"id": "1",
				"properties": {"name": "New Company"},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}
		],
		"startedAt": "2024-01-01T00:00:00.000Z",
		"completedAt": "2024-01-01T00:00:05.000Z"
	}`

	server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/companies/batch/create", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchCreateCompaniesInput{
		Inputs: []CreateCompanyInput{
			{Properties: map[string]string{"name": "New Company"}},
		},
	}

	resp, err := companiesClient.BatchCreateCompanies(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "COMPLETE", resp.Status)
}

// TestBatchUpdateCompanies tests batch update
func TestBatchUpdateCompanies_Success(t *testing.T) {
	responseJSON := `{
		"status": "COMPLETE",
		"results": [],
		"startedAt": "2024-01-01T00:00:00.000Z",
		"completedAt": "2024-01-01T00:00:05.000Z"
	}`

	server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/companies/batch/update", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchUpdateCompaniesInput{}
	resp, err := companiesClient.BatchUpdateCompanies(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

// TestBatchArchiveCompanies tests batch archive
func TestBatchArchiveCompanies_Success(t *testing.T) {
	server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/companies/batch/archive", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	input := &BatchArchiveCompaniesInput{}
	err := companiesClient.BatchArchiveCompanies(context.Background(), input)

	require.NoError(t, err)
}

// TestSearchCompanies tests search functionality
func TestSearchCompanies_Success(t *testing.T) {
	responseJSON := `{
		"total": 1,
		"results": [
			{
				"id": "1",
				"properties": {"name": "Acme Corp"},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}
		],
		"paging": null
	}`

	server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/companies/search", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &SearchCompaniesInput{
		FilterGroups: []FilterGroup{
			{
				Filters: []Filter{
					{
						PropertyName: "name",
						Operator:     "EQ",
						Value:        "Acme Corp",
					},
				},
			},
		},
		Properties: []string{"name", "domain"},
		Limit:      10,
	}

	resp, err := companiesClient.SearchCompanies(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.Total)
	assert.Len(t, resp.Results, 1)
}

// TestOptions tests all option functions
func TestOptions(t *testing.T) {
	t.Run("WithPropertiesWithHistory", func(t *testing.T) {
		server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "name,industry", r.URL.Query().Get("propertiesWithHistory"))
			respondJSON(w, http.StatusOK, `{
				"id": "123",
				"properties": {"name": "Test"},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}`)
		})
		defer server.Close()

		_, err := companiesClient.GetCompany(context.Background(), "123",
			WithPropertiesWithHistory([]string{"name", "industry"}))
		require.NoError(t, err)
	})

	t.Run("WithAssociations", func(t *testing.T) {
		server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "contacts,deals", r.URL.Query().Get("associations"))
			respondJSON(w, http.StatusOK, `{
				"id": "123",
				"properties": {"name": "Test"},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}`)
		})
		defer server.Close()

		_, err := companiesClient.GetCompany(context.Background(), "123",
			WithAssociations([]string{"contacts", "deals"}))
		require.NoError(t, err)
	})

	t.Run("WithArchived", func(t *testing.T) {
		server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "true", r.URL.Query().Get("archived"))
			respondJSON(w, http.StatusOK, `{"results": [], "paging": null}`)
		})
		defer server.Close()

		_, err := companiesClient.ListCompanies(context.Background(), WithArchived())
		require.NoError(t, err)
	})

	t.Run("WithIDProperty", func(t *testing.T) {
		server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "domain", r.URL.Query().Get("idProperty"))
			respondJSON(w, http.StatusOK, `{
				"id": "123",
				"properties": {"name": "Test", "domain": "test.com"},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}`)
		})
		defer server.Close()

		_, err := companiesClient.GetCompany(context.Background(), "test.com",
			WithIDProperty("domain"))
		require.NoError(t, err)
	})
}

// TestErrorHandling tests various error scenarios
func TestErrorHandling(t *testing.T) {
	t.Run("UpdateCompany_Error", func(t *testing.T) {
		server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			respondJSON(w, http.StatusBadRequest, `{"status": "error", "message": "Invalid properties"}`)
		})
		defer server.Close()

		input := &UpdateCompanyInput{Properties: map[string]string{"invalid": "data"}}
		_, err := companiesClient.UpdateCompany(context.Background(), "123", input)
		require.Error(t, err)
	})

	t.Run("ListCompanies_Error", func(t *testing.T) {
		server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			respondJSON(w, http.StatusInternalServerError, `{"status": "error", "message": "Server error"}`)
		})
		defer server.Close()

		_, err := companiesClient.ListCompanies(context.Background())
		require.Error(t, err)
	})

	t.Run("BatchReadCompanies_Error", func(t *testing.T) {
		server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			respondJSON(w, http.StatusBadRequest, `{"status": "error", "message": "Invalid batch input"}`)
		})
		defer server.Close()

		input := &BatchReadCompaniesInput{}
		_, err := companiesClient.BatchReadCompanies(context.Background(), input)
		require.Error(t, err)
	})

	t.Run("BatchCreateCompanies_Error", func(t *testing.T) {
		server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			respondJSON(w, http.StatusBadRequest, `{"status": "error", "message": "Invalid batch create"}`)
		})
		defer server.Close()

		input := &BatchCreateCompaniesInput{}
		_, err := companiesClient.BatchCreateCompanies(context.Background(), input)
		require.Error(t, err)
	})

	t.Run("BatchUpdateCompanies_Error", func(t *testing.T) {
		server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			respondJSON(w, http.StatusBadRequest, `{"status": "error", "message": "Invalid batch update"}`)
		})
		defer server.Close()

		input := &BatchUpdateCompaniesInput{}
		_, err := companiesClient.BatchUpdateCompanies(context.Background(), input)
		require.Error(t, err)
	})

	t.Run("BatchArchiveCompanies_Error", func(t *testing.T) {
		server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			respondJSON(w, http.StatusBadRequest, `{"status": "error", "message": "Invalid batch archive"}`)
		})
		defer server.Close()

		input := &BatchArchiveCompaniesInput{}
		err := companiesClient.BatchArchiveCompanies(context.Background(), input)
		require.Error(t, err)
	})

	t.Run("SearchCompanies_Error", func(t *testing.T) {
		server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			respondJSON(w, http.StatusBadRequest, `{"status": "error", "message": "Invalid search query"}`)
		})
		defer server.Close()

		input := &SearchCompaniesInput{}
		_, err := companiesClient.SearchCompanies(context.Background(), input)
		require.Error(t, err)
	})

	t.Run("ArchiveCompany_Error", func(t *testing.T) {
		server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			respondJSON(w, http.StatusNotFound, `{"status": "error", "message": "Company not found"}`)
		})
		defer server.Close()

		err := companiesClient.ArchiveCompany(context.Background(), "999999")
		require.Error(t, err)
	})
}

// TestInvalidJSON tests handling of invalid JSON responses
func TestInvalidJSON(t *testing.T) {
	t.Run("CreateCompany_InvalidJSON", func(t *testing.T) {
		server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("invalid json"))
		})
		defer server.Close()

		input := &CreateCompanyInput{Properties: map[string]string{"name": "Test"}}
		_, err := companiesClient.CreateCompany(context.Background(), input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unmarshal")
	})

	t.Run("ListCompanies_InvalidJSON", func(t *testing.T) {
		server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("not json"))
		})
		defer server.Close()

		_, err := companiesClient.ListCompanies(context.Background())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unmarshal")
	})

	t.Run("SearchCompanies_InvalidJSON", func(t *testing.T) {
		server, companiesClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{invalid}"))
		})
		defer server.Close()

		input := &SearchCompaniesInput{}
		_, err := companiesClient.SearchCompanies(context.Background(), input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unmarshal")
	})
}
