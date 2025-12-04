package schemas

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/josiah-hester/go-hubspot-sdk/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupMockServer creates a test server with custom handler
func setupMockServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client) {
	server := httptest.NewServer(handler)

	apiClient, err := client.NewClient(
		client.WithTimeout(5*time.Second),
		client.WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	return server, NewClient(apiClient)
}

// respondJSON writes a JSON string response
func respondJSON(w http.ResponseWriter, statusCode int, jsonString string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(jsonString))
}

// TestNewClient tests client creation
func TestNewClient(t *testing.T) {
	apiClient, err := client.NewClient(
		client.WithTimeout(5 * time.Second),
	)
	require.NoError(t, err)

	schemasClient := NewClient(apiClient)
	assert.NotNil(t, schemasClient)
	assert.NotNil(t, schemasClient.apiClient)
}

// TestGetAllSchemas_Success tests successful schemas retrieval
func TestGetAllSchemas_Success(t *testing.T) {
	schemasJSON := `{
		"results": [
			{
				"id": "2-123456",
				"name": "custom_object",
				"labels": {
					"singular": "Custom Object",
					"plural": "Custom Objects"
				},
				"requiredProperties": ["name"],
				"properties": [
					{
						"name": "name",
						"label": "Name",
						"type": "string",
						"fieldType": "text",
						"description": "Object name",
						"groupName": "customobjectinformation",
						"options": []
					}
				],
				"associations": [],
				"archived": false,
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z"
			}
		]
	}`

	server, schemasClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/crm-object-schemas/v3/schemas", r.URL.Path)
		respondJSON(w, http.StatusOK, schemasJSON)
	})
	defer server.Close()

	result, err := schemasClient.GetAllSchemas(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Results, 1)
	assert.Equal(t, "2-123456", result.Results[0].ID)
	assert.Equal(t, "custom_object", result.Results[0].Name)
}

// TestGetAllSchemas_WithArchived tests GetAllSchemas with archived option
func TestGetAllSchemas_WithArchived(t *testing.T) {
	schemasJSON := `{
		"results": [
			{
				"id": "2-123456",
				"name": "custom_object",
				"labels": {
					"singular": "Custom Object",
					"plural": "Custom Objects"
				},
				"requiredProperties": [],
				"properties": [],
				"associations": [],
				"archived": true,
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z"
			}
		]
	}`

	server, schemasClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "true", r.URL.Query().Get("archived"))
		respondJSON(w, http.StatusOK, schemasJSON)
	})
	defer server.Close()

	result, err := schemasClient.GetAllSchemas(context.Background(), WithArchived())

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Results, 1)
	assert.True(t, result.Results[0].Archived)
}

// TestGetAllSchemas_NoResults tests when no schemas found
func TestGetAllSchemas_NoResults(t *testing.T) {
	schemasJSON := `{
		"results": []
	}`

	server, schemasClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, schemasJSON)
	})
	defer server.Close()

	result, err := schemasClient.GetAllSchemas(context.Background())

	require.Error(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, err.Error(), "no schemas found")
}

// TestGetAllSchemas_InvalidJSON tests invalid JSON response
func TestGetAllSchemas_InvalidJSON(t *testing.T) {
	server, schemasClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	result, err := schemasClient.GetAllSchemas(context.Background())

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestGetExistingSchema_Success tests successful schema retrieval
func TestGetExistingSchema_Success(t *testing.T) {
	schemaJSON := `{
		"id": "2-123456",
		"name": "custom_object",
		"labels": {
			"singular": "Custom Object",
			"plural": "Custom Objects"
		},
		"requiredProperties": ["name"],
		"properties": [
			{
				"name": "name",
				"label": "Name",
				"type": "string",
				"fieldType": "text",
				"description": "Object name",
				"groupName": "customobjectinformation",
				"options": []
			}
		],
		"associations": [],
		"archived": false,
		"createdAt": "2024-01-01T00:00:00.000Z",
		"updatedAt": "2024-01-01T00:00:00.000Z",
		"primaryDisplayProperty": "name"
	}`

	server, schemasClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/crm-object-schemas/v3/schemas/custom_object", r.URL.Path)
		respondJSON(w, http.StatusOK, schemaJSON)
	})
	defer server.Close()

	schema, err := schemasClient.GetExistingSchema(context.Background(), "custom_object")

	require.NoError(t, err)
	assert.NotNil(t, schema)
	assert.Equal(t, "2-123456", schema.ID)
	assert.Equal(t, "custom_object", schema.Name)
	assert.Equal(t, "name", schema.PrimaryDisplayProperty)
}

// TestGetExistingSchema_NotFound tests 404 error handling
func TestGetExistingSchema_NotFound(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Schema not found",
		"category": "OBJECT_NOT_FOUND"
	}`

	server, schemasClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, errorJSON)
	})
	defer server.Close()

	schema, err := schemasClient.GetExistingSchema(context.Background(), "nonexistent")

	require.Error(t, err)
	assert.Nil(t, schema)
}

// TestGetExistingSchema_InvalidJSON tests invalid JSON response
func TestGetExistingSchema_InvalidJSON(t *testing.T) {
	server, schemasClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	schema, err := schemasClient.GetExistingSchema(context.Background(), "custom_object")

	require.Error(t, err)
	assert.Nil(t, schema)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestCreateNewSchema_Success tests successful schema creation
func TestCreateNewSchema_Success(t *testing.T) {
	responseJSON := `{
		"id": "2-123456",
		"name": "custom_object",
		"labels": {
			"singular": "Custom Object",
			"plural": "Custom Objects"
		},
		"requiredProperties": ["name"],
		"properties": [
			{
				"name": "name",
				"label": "Name",
				"type": "string",
				"fieldType": "text",
				"description": "Object name",
				"groupName": "customobjectinformation",
				"options": []
			}
		],
		"associations": [],
		"archived": false,
		"createdAt": "2024-01-01T00:00:00.000Z",
		"updatedAt": "2024-01-01T00:00:00.000Z"
	}`

	server, schemasClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm-object-schemas/v3/schemas", r.URL.Path)
		respondJSON(w, http.StatusCreated, responseJSON)
	})
	defer server.Close()

	input := &CreateNewSchemaInput{
		Name:               "custom_object",
		RequiredProperties: []string{"name"},
		AssociatedObjects:  []string{"contacts"},
		Properties: []Property{
			{
				Name:        "name",
				Label:       "Name",
				Type:        "string",
				FieldType:   "text",
				Description: "Object name",
				GroupName:   "customobjectinformation",
				Options:     []Option{},
			},
		},
	}

	schema, err := schemasClient.CreateNewSchema(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, schema)
	assert.Equal(t, "2-123456", schema.ID)
	assert.Equal(t, "custom_object", schema.Name)
}

// TestCreateNewSchema_ValidationError tests validation error
func TestCreateNewSchema_ValidationError(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Invalid schema name",
		"category": "VALIDATION_ERROR"
	}`

	server, schemasClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusBadRequest, errorJSON)
	})
	defer server.Close()

	input := &CreateNewSchemaInput{
		Name:               "",
		RequiredProperties: []string{},
		AssociatedObjects:  []string{},
		Properties:         []Property{},
	}

	schema, err := schemasClient.CreateNewSchema(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, schema)
}

// TestCreateNewSchema_InvalidJSON tests invalid JSON response
func TestCreateNewSchema_InvalidJSON(t *testing.T) {
	server, schemasClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	input := &CreateNewSchemaInput{
		Name:               "custom_object",
		RequiredProperties: []string{},
		AssociatedObjects:  []string{},
		Properties:         []Property{},
	}

	schema, err := schemasClient.CreateNewSchema(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, schema)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestCreateNewAssociationSchema_Success tests successful association creation
func TestCreateNewAssociationSchema_Success(t *testing.T) {
	responseJSON := `{
		"id": "123",
		"fromObjectTypeId": "2-123456",
		"toObjectTypeId": "0-1",
		"name": "custom_to_contact",
		"createdAt": "2024-01-01T00:00:00.000Z",
		"updatedAt": "2024-01-01T00:00:00.000Z"
	}`

	server, schemasClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm-object-schemas/v3/schemas/custom_object/associations", r.URL.Path)
		respondJSON(w, http.StatusCreated, responseJSON)
	})
	defer server.Close()

	input := &CreateNewAssociationSchemaInput{
		FromObjectTypeID: "2-123456",
		ToObjectTypeID:   "0-1",
		Name:             "custom_to_contact",
	}

	assoc, err := schemasClient.CreateNewAssociationSchema(context.Background(), "custom_object", input)

	require.NoError(t, err)
	assert.NotNil(t, assoc)
	assert.Equal(t, "123", assoc.ID)
	assert.Equal(t, "2-123456", assoc.FromObjectTypeID)
	assert.Equal(t, "0-1", assoc.ToObjectTypeID)
}

// TestCreateNewAssociationSchema_ValidationError tests validation error
func TestCreateNewAssociationSchema_ValidationError(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Invalid object type",
		"category": "VALIDATION_ERROR"
	}`

	server, schemasClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusBadRequest, errorJSON)
	})
	defer server.Close()

	input := &CreateNewAssociationSchemaInput{
		FromObjectTypeID: "",
		ToObjectTypeID:   "",
	}

	assoc, err := schemasClient.CreateNewAssociationSchema(context.Background(), "custom_object", input)

	require.Error(t, err)
	assert.Nil(t, assoc)
}

// TestCreateNewAssociationSchema_InvalidJSON tests invalid JSON response
func TestCreateNewAssociationSchema_InvalidJSON(t *testing.T) {
	server, schemasClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	input := &CreateNewAssociationSchemaInput{
		FromObjectTypeID: "2-123456",
		ToObjectTypeID:   "0-1",
	}

	assoc, err := schemasClient.CreateNewAssociationSchema(context.Background(), "custom_object", input)

	require.Error(t, err)
	assert.Nil(t, assoc)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestUpdateSchema_Success tests successful schema update
func TestUpdateSchema_Success(t *testing.T) {
	responseJSON := `{
		"id": "2-123456",
		"name": "custom_object",
		"labels": {
			"singular": "Updated Object",
			"plural": "Updated Objects"
		},
		"requiredProperties": ["name"],
		"properties": [],
		"associations": [],
		"archived": false,
		"updatedAt": "2024-01-02T00:00:00.000Z"
	}`

	server, schemasClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method)
		assert.Equal(t, "/crm-object-schemas/v3/schemas/custom_object", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &UpdateSchemaInput{
		Labels: struct {
			Singular string `json:"singular"`
			Plural   string `json:"plural"`
		}{
			Singular: "Updated Object",
			Plural:   "Updated Objects",
		},
	}

	schema, err := schemasClient.UpdateSchema(context.Background(), "custom_object", input)

	require.NoError(t, err)
	assert.NotNil(t, schema)
	assert.Equal(t, "2-123456", schema.ID)
	assert.Equal(t, "Updated Object", schema.Labels.Singular)
}

// TestUpdateSchema_NotFound tests 404 error on update
func TestUpdateSchema_NotFound(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Schema not found",
		"category": "OBJECT_NOT_FOUND"
	}`

	server, schemasClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, errorJSON)
	})
	defer server.Close()

	input := &UpdateSchemaInput{
		Description: "Updated description",
	}

	schema, err := schemasClient.UpdateSchema(context.Background(), "nonexistent", input)

	require.Error(t, err)
	assert.Nil(t, schema)
}

// TestUpdateSchema_InvalidJSON tests invalid JSON response
func TestUpdateSchema_InvalidJSON(t *testing.T) {
	server, schemasClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	input := &UpdateSchemaInput{
		Description: "Updated description",
	}

	schema, err := schemasClient.UpdateSchema(context.Background(), "custom_object", input)

	require.Error(t, err)
	assert.Nil(t, schema)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestDeleteSchema_Success tests successful schema deletion
func TestDeleteSchema_Success(t *testing.T) {
	server, schemasClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/crm-object-schemas/v3/schemas/custom_object", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := schemasClient.DeleteSchema(context.Background(), "custom_object")

	assert.NoError(t, err)
}

// TestDeleteSchema_WithArchived tests deletion with archived option
func TestDeleteSchema_WithArchived(t *testing.T) {
	server, schemasClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "true", r.URL.Query().Get("archived"))
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := schemasClient.DeleteSchema(context.Background(), "custom_object", WithArchived())

	assert.NoError(t, err)
}

// TestDeleteSchema_NotFound tests 404 error on delete
func TestDeleteSchema_NotFound(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Schema not found",
		"category": "OBJECT_NOT_FOUND"
	}`

	server, schemasClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, errorJSON)
	})
	defer server.Close()

	err := schemasClient.DeleteSchema(context.Background(), "nonexistent")

	require.Error(t, err)
}

// TestRemoveAssociationSchema_Success tests successful association removal
func TestRemoveAssociationSchema_Success(t *testing.T) {
	server, schemasClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/crm-object-schemas/v3/schemas/custom_object/associations/123", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := schemasClient.RemoveAssociationSchema(context.Background(), "custom_object", "123")

	assert.NoError(t, err)
}

// TestRemoveAssociationSchema_NotFound tests 404 error on remove
func TestRemoveAssociationSchema_NotFound(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Association not found",
		"category": "OBJECT_NOT_FOUND"
	}`

	server, schemasClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, errorJSON)
	})
	defer server.Close()

	err := schemasClient.RemoveAssociationSchema(context.Background(), "custom_object", "999")

	require.Error(t, err)
}
