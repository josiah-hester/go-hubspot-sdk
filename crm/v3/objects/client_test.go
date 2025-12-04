package objects

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

var TestObj = Object{
	ID:                    "1234567890",
	Properties:            make(map[string]string),
	CreatedAt:             "2023-11-07T05:31:56Z",
	UpdatedAt:             "2023-11-07T05:31:56Z",
	ArchivedAt:            "2023-11-07T05:31:56Z",
	Archived:              true,
	Associations:          make(map[string]AssociationResponse),
	PropertiesWithHistory: make(map[string][]PropertyWithHistory),
	ObjectWriteTraceID:    "",
}

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

// respondJSON writes a JSON string resposne
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

	objectsClient := NewClient(apiClient)
	assert.NotNil(t, objectsClient)
	assert.NotNil(t, objectsClient.apiClient)
}

// TestListObjects_Success tests successful objects retrieval
func TestListObjects_Success(t *testing.T) {
	objectJSON := `{
    "results": [
      {
        "archived": true,
        "createdAt": "2023-11-07T05:31:56Z",
        "id": "1234567890",
        "properties": {},
        "updatedAt": "2023-11-07T05:31:56Z",
        "archivedAt": "2023-11-07T05:31:56Z",
        "associations": {},
        "objectWriteTraceId": "",
        "propertiesWithHistory": {},
        "url": "test.com"
      }
    ],
    "paging": {
      "next": {
        "after": "NTI1Cg%3D%3D",
        "link": "?after=NTI1Cg%3D%3D"
      },
      "prev": {
        "before": "NTI1Cg%3D%3D",
        "link": "?before=NTI1Cg%3D%3D"
      }
    }
  }`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/crm/v3/objects/contacts", r.URL.Path)
		respondJSON(w, http.StatusOK, objectJSON)
	})
	defer server.Close()

	objects, paging, err := objectClient.ListObjects(context.Background(), "contacts")

	require.NoError(t, err)
	assert.NotNil(t, objects)
	assert.NotNil(t, paging)
	assert.Contains(t, objects, TestObj)
}

// TestListObjects_WithOptions tests list with various options
func TestListObjects_WithOptions(t *testing.T) {
	objectJSON := `{
		"results": [
			{
				"archived": false,
				"createdAt": "2023-11-07T05:31:56Z",
				"id": "1234567890",
				"properties": {
					"email": "test@example.com",
					"firstname": "John"
				},
				"updatedAt": "2023-11-07T05:31:56Z"
			}
		]
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "10", r.URL.Query().Get("limit"))
		assert.Equal(t, "cursor123", r.URL.Query().Get("after"))
		assert.Equal(t, "email,firstname", r.URL.Query().Get("properties"))
		respondJSON(w, http.StatusOK, objectJSON)
	})
	defer server.Close()

	objects, _, err := objectClient.ListObjects(
		context.Background(),
		"contacts",
		WithLimit(10),
		WithAfter("cursor123"),
		WithProperties([]string{"email", "firstname"}),
	)

	require.NoError(t, err)
	assert.Len(t, objects, 1)
}

// TestListObjects_WithPropertiesWithHistory tests list with properties history
func TestListObjects_WithPropertiesWithHistory(t *testing.T) {
	objectJSON := `{
		"results": [
			{
				"archived": false,
				"createdAt": "2023-11-07T05:31:56Z",
				"id": "1234567890",
				"properties": {
					"email": "test@example.com"
				},
				"updatedAt": "2023-11-07T05:31:56Z",
				"propertiesWithHistory": {
					"email": [
						{
							"value": "test@example.com",
							"timestamp": "2023-11-07T05:31:56Z",
							"sourceType": "API"
						}
					]
				}
			}
		]
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "email", r.URL.Query().Get("propertiesWithHistory"))
		respondJSON(w, http.StatusOK, objectJSON)
	})
	defer server.Close()

	objects, _, err := objectClient.ListObjects(
		context.Background(),
		"contacts",
		WithPropertiesWithHistory([]string{"email"}),
	)

	require.NoError(t, err)
	assert.Len(t, objects, 1)
}

// TestListObjects_WithArchived tests list with archived option
func TestListObjects_WithArchived(t *testing.T) {
	objectJSON := `{
		"results": [
			{
				"archived": true,
				"createdAt": "2023-11-07T05:31:56Z",
				"id": "1234567890",
				"properties": {},
				"updatedAt": "2023-11-07T05:31:56Z",
				"archivedAt": "2023-11-07T05:31:56Z"
			}
		]
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "true", r.URL.Query().Get("archived"))
		respondJSON(w, http.StatusOK, objectJSON)
	})
	defer server.Close()

	objects, _, err := objectClient.ListObjects(
		context.Background(),
		"contacts",
		WithArchived(),
	)

	require.NoError(t, err)
	assert.Len(t, objects, 1)
	assert.True(t, objects[0].Archived)
}

// TestListObjects_NoResults tests when no objects found
func TestListObjects_NoResults(t *testing.T) {
	objectJSON := `{
		"results": []
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, objectJSON)
	})
	defer server.Close()

	objects, _, err := objectClient.ListObjects(context.Background(), "contacts")

	require.Error(t, err)
	assert.Nil(t, objects)
	assert.Contains(t, err.Error(), "no objects found")
}

// TestListObjects_InvalidJSON tests invalid JSON response
func TestListObjects_InvalidJSON(t *testing.T) {
	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	objects, _, err := objectClient.ListObjects(context.Background(), "contacts")

	require.Error(t, err)
	assert.Nil(t, objects)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestCreateObject_Success tests successful object creation
func TestCreateObject_Success(t *testing.T) {
	responseJSON := `{
		"createResourceId": "resource-123",
		"entity": {
			"id": "1234567890",
			"properties": {
				"email": "test@example.com",
				"firstname": "John",
				"lastname": "Doe"
			},
			"createdAt": "2024-01-01T00:00:00.000Z",
			"updatedAt": "2024-01-01T00:00:00.000Z",
			"archived": false
		}
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/contacts", r.URL.Path)
		respondJSON(w, http.StatusCreated, responseJSON)
	})
	defer server.Close()

	input := &CreateObjectInput{
		Properties: map[string]string{
			"email":     "test@example.com",
			"firstname": "John",
			"lastname":  "Doe",
		},
	}

	object, err := objectClient.CreateObject(context.Background(), input, "contacts")

	require.NoError(t, err)
	assert.NotNil(t, object)
	assert.Equal(t, "1234567890", object.ID)
	assert.Equal(t, "test@example.com", object.Properties["email"])
}

// TestCreateObject_ValidationError tests validation error
func TestCreateObject_ValidationError(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Invalid email",
		"category": "VALIDATION_ERROR"
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusBadRequest, errorJSON)
	})
	defer server.Close()

	input := &CreateObjectInput{
		Properties: map[string]string{
			"email": "invalid-email",
		},
	}

	object, err := objectClient.CreateObject(context.Background(), input, "contacts")

	require.Error(t, err)
	assert.Nil(t, object)
}

// TestCreateObject_InvalidJSON tests invalid JSON response
func TestCreateObject_InvalidJSON(t *testing.T) {
	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	input := &CreateObjectInput{
		Properties: map[string]string{
			"email": "test@example.com",
		},
	}

	object, err := objectClient.CreateObject(context.Background(), input, "contacts")

	require.Error(t, err)
	assert.Nil(t, object)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestReadObject_Success tests successful object retrieval
func TestReadObject_Success(t *testing.T) {
	objectJSON := `{
		"id": "1234567890",
		"properties": {
			"email": "test@example.com",
			"firstname": "John",
			"lastname": "Doe"
		},
		"createdAt": "2024-01-01T00:00:00.000Z",
		"updatedAt": "2024-01-01T00:00:00.000Z",
		"archived": false
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/crm/v3/objects/contacts/1234567890", r.URL.Path)
		respondJSON(w, http.StatusOK, objectJSON)
	})
	defer server.Close()

	object, err := objectClient.ReadObject(context.Background(), "contacts", "1234567890")

	require.NoError(t, err)
	assert.NotNil(t, object)
	assert.Equal(t, "1234567890", object.ID)
	assert.Equal(t, "test@example.com", object.Properties["email"])
}

// TestReadObject_WithOptions tests ReadObject with options
func TestReadObject_WithOptions(t *testing.T) {
	objectJSON := `{
		"id": "1234567890",
		"properties": {
			"email": "test@example.com"
		},
		"archived": false
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "email,firstname", r.URL.Query().Get("properties"))
		assert.Equal(t, "companies", r.URL.Query().Get("associations"))
		assert.Equal(t, "email", r.URL.Query().Get("idProperty"))
		respondJSON(w, http.StatusOK, objectJSON)
	})
	defer server.Close()

	object, err := objectClient.ReadObject(
		context.Background(),
		"contacts",
		"test@example.com",
		WithProperties([]string{"email", "firstname"}),
		WithAssociations([]string{"companies"}),
		WithIDProperty("email"),
	)

	require.NoError(t, err)
	assert.NotNil(t, object)
}

// TestReadObject_NotFound tests 404 error handling
func TestReadObject_NotFound(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Object not found",
		"category": "OBJECT_NOT_FOUND"
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, errorJSON)
	})
	defer server.Close()

	object, err := objectClient.ReadObject(context.Background(), "contacts", "99999")

	require.Error(t, err)
	assert.Nil(t, object)
}

// TestReadObject_InvalidJSON tests invalid JSON response
func TestReadObject_InvalidJSON(t *testing.T) {
	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	object, err := objectClient.ReadObject(context.Background(), "contacts", "1234567890")

	require.Error(t, err)
	assert.Nil(t, object)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestUpdateObject_Success tests successful object update
func TestUpdateObject_Success(t *testing.T) {
	responseJSON := `{
		"id": "1234567890",
		"properties": {
			"email": "test@example.com",
			"firstname": "Jane",
			"lastname": "Doe"
		},
		"updatedAt": "2024-01-02T00:00:00.000Z",
		"archived": false
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method)
		assert.Equal(t, "/crm/v3/objects/contacts/1234567890", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &UpdateObjectInput{
		Properties: map[string]string{
			"firstname": "Jane",
		},
	}

	object, err := objectClient.UpdateObject(context.Background(), "contacts", "1234567890", input)

	require.NoError(t, err)
	assert.NotNil(t, object)
	assert.Equal(t, "1234567890", object.ID)
	assert.Equal(t, "Jane", object.Properties["firstname"])
}

// TestUpdateObject_WithIDProperty tests update with custom ID property
func TestUpdateObject_WithIDProperty(t *testing.T) {
	responseJSON := `{
		"id": "1234567890",
		"properties": {
			"email": "test@example.com",
			"firstname": "Jane"
		},
		"archived": false
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "email", r.URL.Query().Get("idProperty"))
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &UpdateObjectInput{
		Properties: map[string]string{
			"firstname": "Jane",
		},
	}

	object, err := objectClient.UpdateObject(
		context.Background(),
		"contacts",
		"test@example.com",
		input,
		WithIDProperty("email"),
	)

	require.NoError(t, err)
	assert.NotNil(t, object)
}

// TestUpdateObject_NotFound tests 404 error on update
func TestUpdateObject_NotFound(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Object not found",
		"category": "OBJECT_NOT_FOUND"
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, errorJSON)
	})
	defer server.Close()

	input := &UpdateObjectInput{
		Properties: map[string]string{
			"firstname": "Jane",
		},
	}

	object, err := objectClient.UpdateObject(context.Background(), "contacts", "99999", input)

	require.Error(t, err)
	assert.Nil(t, object)
}

// TestUpdateObject_ValidationError tests validation error on update
func TestUpdateObject_ValidationError(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Invalid property",
		"category": "VALIDATION_ERROR"
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusBadRequest, errorJSON)
	})
	defer server.Close()

	input := &UpdateObjectInput{
		Properties: map[string]string{
			"invalid_field": "value",
		},
	}

	object, err := objectClient.UpdateObject(context.Background(), "contacts", "1234567890", input)

	require.Error(t, err)
	assert.Nil(t, object)
}

// TestUpdateObject_InvalidJSON tests invalid JSON response
func TestUpdateObject_InvalidJSON(t *testing.T) {
	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	input := &UpdateObjectInput{
		Properties: map[string]string{
			"firstname": "Jane",
		},
	}

	object, err := objectClient.UpdateObject(context.Background(), "contacts", "1234567890", input)

	require.Error(t, err)
	assert.Nil(t, object)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestArchiveObject_Success tests successful object archival
func TestArchiveObject_Success(t *testing.T) {
	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/crm/v3/objects/contacts/1234567890", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := objectClient.ArchiveObject(context.Background(), "contacts", "1234567890")

	assert.NoError(t, err)
}

// TestArchiveObject_NotFound tests 404 error on archive
func TestArchiveObject_NotFound(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Object not found",
		"category": "OBJECT_NOT_FOUND"
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, errorJSON)
	})
	defer server.Close()

	err := objectClient.ArchiveObject(context.Background(), "contacts", "99999")

	require.Error(t, err)
}

// TestMergeObjects_Success tests successful object merge
func TestMergeObjects_Success(t *testing.T) {
	responseJSON := `{
		"id": "1234567890",
		"properties": {
			"email": "test@example.com",
			"firstname": "John",
			"lastname": "Doe"
		},
		"createdAt": "2024-01-01T00:00:00.000Z",
		"updatedAt": "2024-01-02T00:00:00.000Z",
		"archived": false
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/contacts/merge", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &MergeObjectsInput{
		PrimaryObjectID: "1234567890",
		ObjectIDToMerge: "9876543210",
	}

	object, err := objectClient.MergeObjects(context.Background(), "contacts", input)

	require.NoError(t, err)
	assert.NotNil(t, object)
	assert.Equal(t, "1234567890", object.ID)
}

// TestMergeObjects_ValidationError tests validation error on merge
func TestMergeObjects_ValidationError(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Invalid object IDs",
		"category": "VALIDATION_ERROR"
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusBadRequest, errorJSON)
	})
	defer server.Close()

	input := &MergeObjectsInput{
		PrimaryObjectID: "",
		ObjectIDToMerge: "",
	}

	object, err := objectClient.MergeObjects(context.Background(), "contacts", input)

	require.Error(t, err)
	assert.Nil(t, object)
}

// TestMergeObjects_InvalidJSON tests invalid JSON response
func TestMergeObjects_InvalidJSON(t *testing.T) {
	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	input := &MergeObjectsInput{
		PrimaryObjectID: "1234567890",
		ObjectIDToMerge: "9876543210",
	}

	object, err := objectClient.MergeObjects(context.Background(), "contacts", input)

	require.Error(t, err)
	assert.Nil(t, object)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestBatchReadObjects_Success tests successful batch read
func TestBatchReadObjects_Success(t *testing.T) {
	responseJSON := `{
		"completedAt": "2024-01-01T00:00:05.000Z",
		"startedAt": "2024-01-01T00:00:00.000Z",
		"status": "COMPLETE",
		"results": [
			{
				"id": "1",
				"properties": {
					"email": "test1@example.com",
					"firstname": "John"
				},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			},
			{
				"id": "2",
				"properties": {
					"email": "test2@example.com",
					"firstname": "Jane"
				},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}
		]
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/contacts/batch/read", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchReadObjectsInput{
		Inputs: []struct {
			ID string `json:"id" required:"yes"`
		}{
			{ID: "1"},
			{ID: "2"},
		},
		Properties:            []string{"email", "firstname"},
		PropertiesWithHistory: []string{},
	}

	result, err := objectClient.BatchReadObjects(context.Background(), "contacts", input)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, Complete, result.Status)
	assert.Len(t, result.Results, 2)
}

// TestBatchReadObjects_WithErrors tests batch read with errors
func TestBatchReadObjects_WithErrors(t *testing.T) {
	responseJSON := `{
		"completedAt": "2024-01-01T00:00:05.000Z",
		"startedAt": "2024-01-01T00:00:00.000Z",
		"status": "COMPLETE",
		"results": [
			{
				"id": "1",
				"properties": {
					"email": "test1@example.com"
				},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}
		],
		"numErrors": 1,
		"errors": [
			{
				"status": "error",
				"category": "OBJECT_NOT_FOUND",
				"message": "Object not found",
				"context": {},
				"links": {},
				"errors": [
					{
						"message": "Object with ID 99999 not found"
					}
				]
			}
		]
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchReadObjectsInput{
		Inputs: []struct {
			ID string `json:"id" required:"yes"`
		}{
			{ID: "1"},
			{ID: "99999"},
		},
		Properties:            []string{"email"},
		PropertiesWithHistory: []string{},
	}

	result, err := objectClient.BatchReadObjects(context.Background(), "contacts", input)

	require.Error(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.NumErrors)
	assert.Len(t, result.Errors, 1)
}

// TestBatchReadObjects_InvalidJSON tests invalid JSON response
func TestBatchReadObjects_InvalidJSON(t *testing.T) {
	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	input := &BatchReadObjectsInput{
		Inputs: []struct {
			ID string `json:"id" required:"yes"`
		}{{ID: "1"}},
		Properties:            []string{"email"},
		PropertiesWithHistory: []string{},
	}

	result, err := objectClient.BatchReadObjects(context.Background(), "contacts", input)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestBatchCreateObjects_Success tests successful batch create
func TestBatchCreateObjects_Success(t *testing.T) {
	responseJSON := `{
		"completedAt": "2024-01-01T00:00:05.000Z",
		"startedAt": "2024-01-01T00:00:00.000Z",
		"status": "COMPLETE",
		"results": [
			{
				"id": "1",
				"properties": {
					"email": "test1@example.com",
					"firstname": "John"
				},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			},
			{
				"id": "2",
				"properties": {
					"email": "test2@example.com",
					"firstname": "Jane"
				},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}
		]
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/contacts/batch/create", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchCreateObjectsInput{
		Inputs: []struct {
			Associations       []Association     `json:"associations" required:"yes"`
			Properties         map[string]string `json:"properties" required:"yes"`
			ObjectWriteTraceID string            `json:"objectWriteTraceId"`
		}{
			{
				Properties: map[string]string{
					"email":     "test1@example.com",
					"firstname": "John",
				},
				Associations: []Association{},
			},
			{
				Properties: map[string]string{
					"email":     "test2@example.com",
					"firstname": "Jane",
				},
				Associations: []Association{},
			},
		},
	}

	result, err := objectClient.BatchCreateObjects(context.Background(), "contacts", input)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, Complete, result.Status)
	assert.Len(t, result.Results, 2)
}

// TestBatchCreateObjects_WithErrors tests batch create with errors
func TestBatchCreateObjects_WithErrors(t *testing.T) {
	responseJSON := `{
		"completedAt": "2024-01-01T00:00:05.000Z",
		"startedAt": "2024-01-01T00:00:00.000Z",
		"status": "COMPLETE",
		"results": [
			{
				"id": "1",
				"properties": {
					"email": "test1@example.com"
				},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}
		],
		"numErrors": 1,
		"errors": [
			{
				"status": "error",
				"category": "VALIDATION_ERROR",
				"message": "Invalid email",
				"context": {},
				"links": {},
				"errors": [
					{
						"message": "Email format is invalid"
					}
				]
			}
		]
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchCreateObjectsInput{
		Inputs: []struct {
			Associations       []Association     `json:"associations" required:"yes"`
			Properties         map[string]string `json:"properties" required:"yes"`
			ObjectWriteTraceID string            `json:"objectWriteTraceId"`
		}{
			{
				Properties:   map[string]string{"email": "test@example.com"},
				Associations: []Association{},
			},
			{
				Properties:   map[string]string{"email": "invalid"},
				Associations: []Association{},
			},
		},
	}

	result, err := objectClient.BatchCreateObjects(context.Background(), "contacts", input)

	require.Error(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.NumErrors)
	assert.Len(t, result.Errors, 1)
}

// TestBatchCreateObjects_InvalidJSON tests invalid JSON response
func TestBatchCreateObjects_InvalidJSON(t *testing.T) {
	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	input := &BatchCreateObjectsInput{
		Inputs: []struct {
			Associations       []Association     `json:"associations" required:"yes"`
			Properties         map[string]string `json:"properties" required:"yes"`
			ObjectWriteTraceID string            `json:"objectWriteTraceId"`
		}{
			{
				Properties:   map[string]string{"email": "test@example.com"},
				Associations: []Association{},
			},
		},
	}

	result, err := objectClient.BatchCreateObjects(context.Background(), "contacts", input)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to ubmarshal")
}

// TestBatchUpdateObjects_Success tests successful batch update
func TestBatchUpdateObjects_Success(t *testing.T) {
	responseJSON := `{
		"completedAt": "2024-01-01T00:00:05.000Z",
		"startedAt": "2024-01-01T00:00:00.000Z",
		"status": "COMPLETE",
		"results": [
			{
				"id": "1",
				"properties": {
					"email": "test1@example.com",
					"firstname": "John Updated"
				},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-02T00:00:00.000Z",
				"archived": false
			},
			{
				"id": "2",
				"properties": {
					"email": "test2@example.com",
					"firstname": "Jane Updated"
				},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-02T00:00:00.000Z",
				"archived": false
			}
		]
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/contacts/batch/update", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchUpdateObjectsInput{
		Inputs: []struct {
			ID                 string            `json:"id" required:"yes"`
			Properties         map[string]string `json:"properties" required:"yes"`
			IDProperty         string            `json:"idProperty"`
			ObjectWriteTraceID string            `json:"objectWriteTraceId"`
		}{
			{
				ID:         "1",
				Properties: map[string]string{"firstname": "John Updated"},
			},
			{
				ID:         "2",
				Properties: map[string]string{"firstname": "Jane Updated"},
			},
		},
	}

	result, err := objectClient.BatchUpdateObjects(context.Background(), "contacts", input)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, Complete, result.Status)
	assert.Len(t, result.Results, 2)
}

// TestBatchUpdateObjects_WithErrors tests batch update with errors
func TestBatchUpdateObjects_WithErrors(t *testing.T) {
	responseJSON := `{
		"completedAt": "2024-01-01T00:00:05.000Z",
		"startedAt": "2024-01-01T00:00:00.000Z",
		"status": "COMPLETE",
		"results": [
			{
				"id": "1",
				"properties": {
					"firstname": "John Updated"
				},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-02T00:00:00.000Z",
				"archived": false
			}
		],
		"numErrors": 1,
		"errors": [
			{
				"status": "error",
				"category": "OBJECT_NOT_FOUND",
				"message": "Object not found",
				"context": {},
				"links": {},
				"errors": [
					{
						"message": "Object with ID 99999 not found"
					}
				]
			}
		]
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchUpdateObjectsInput{
		Inputs: []struct {
			ID                 string            `json:"id" required:"yes"`
			Properties         map[string]string `json:"properties" required:"yes"`
			IDProperty         string            `json:"idProperty"`
			ObjectWriteTraceID string            `json:"objectWriteTraceId"`
		}{
			{
				ID:         "1",
				Properties: map[string]string{"firstname": "John Updated"},
			},
			{
				ID:         "99999",
				Properties: map[string]string{"firstname": "Jane Updated"},
			},
		},
	}

	result, err := objectClient.BatchUpdateObjects(context.Background(), "contacts", input)

	require.Error(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.NumErrors)
	assert.Len(t, result.Errors, 1)
}

// TestBatchUpdateObjects_InvalidJSON tests invalid JSON response
func TestBatchUpdateObjects_InvalidJSON(t *testing.T) {
	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	input := &BatchUpdateObjectsInput{
		Inputs: []struct {
			ID                 string            `json:"id" required:"yes"`
			Properties         map[string]string `json:"properties" required:"yes"`
			IDProperty         string            `json:"idProperty"`
			ObjectWriteTraceID string            `json:"objectWriteTraceId"`
		}{
			{
				ID:         "1",
				Properties: map[string]string{"firstname": "John"},
			},
		},
	}

	result, err := objectClient.BatchUpdateObjects(context.Background(), "contacts", input)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to ubmarshal")
}

// TestBatchCreateOrUpdateObjects_Success tests successful batch upsert
func TestBatchCreateOrUpdateObjects_Success(t *testing.T) {
	responseJSON := `{
		"completedAt": "2024-01-01T00:00:05.000Z",
		"startedAt": "2024-01-01T00:00:00.000Z",
		"status": "COMPLETE",
		"results": [
			{
				"id": "1",
				"properties": {
					"email": "existing@example.com",
					"firstname": "John Updated"
				},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-02T00:00:00.000Z",
				"archived": false
			},
			{
				"id": "2",
				"properties": {
					"email": "new@example.com",
					"firstname": "Jane"
				},
				"createdAt": "2024-01-02T00:00:00.000Z",
				"updatedAt": "2024-01-02T00:00:00.000Z",
				"archived": false
			}
		]
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/contacts/batch/upsert", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchCreateOrUpdateObjectsInput{
		Inputs: []struct {
			ID                 string            `json:"id" required:"yes"`
			Properties         map[string]string `json:"properties" required:"yes"`
			IDProperty         string            `json:"idProperty"`
			ObjectWriteTraceID string            `json:"objectWriteTraceId"`
		}{
			{
				ID:         "existing@example.com",
				Properties: map[string]string{"firstname": "John Updated"},
				IDProperty: "email",
			},
			{
				ID:         "new@example.com",
				Properties: map[string]string{"firstname": "Jane"},
				IDProperty: "email",
			},
		},
	}

	result, err := objectClient.BatchCreateOrUpdateObjects(context.Background(), "contacts", input)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, Complete, result.Status)
	assert.Len(t, result.Results, 2)
}

// TestBatchCreateOrUpdateObjects_WithErrors tests batch upsert with errors
func TestBatchCreateOrUpdateObjects_WithErrors(t *testing.T) {
	responseJSON := `{
		"completedAt": "2024-01-01T00:00:05.000Z",
		"startedAt": "2024-01-01T00:00:00.000Z",
		"status": "COMPLETE",
		"results": [
			{
				"id": "1",
				"properties": {
					"email": "test@example.com"
				},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-02T00:00:00.000Z",
				"archived": false
			}
		],
		"numErrors": 1,
		"errors": [
			{
				"status": "error",
				"category": "VALIDATION_ERROR",
				"message": "Validation failed",
				"context": {},
				"links": {},
				"errors": [
					{
						"message": "Invalid property value"
					}
				]
			}
		]
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchCreateOrUpdateObjectsInput{
		Inputs: []struct {
			ID                 string            `json:"id" required:"yes"`
			Properties         map[string]string `json:"properties" required:"yes"`
			IDProperty         string            `json:"idProperty"`
			ObjectWriteTraceID string            `json:"objectWriteTraceId"`
		}{
			{
				ID:         "test@example.com",
				Properties: map[string]string{"email": "test@example.com"},
			},
			{
				ID:         "invalid",
				Properties: map[string]string{"email": "invalid"},
			},
		},
	}

	result, err := objectClient.BatchCreateOrUpdateObjects(context.Background(), "contacts", input)

	require.Error(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.NumErrors)
	assert.Len(t, result.Errors, 1)
}

// TestBatchCreateOrUpdateObjects_InvalidJSON tests invalid JSON response
func TestBatchCreateOrUpdateObjects_InvalidJSON(t *testing.T) {
	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	input := &BatchCreateOrUpdateObjectsInput{
		Inputs: []struct {
			ID                 string            `json:"id" required:"yes"`
			Properties         map[string]string `json:"properties" required:"yes"`
			IDProperty         string            `json:"idProperty"`
			ObjectWriteTraceID string            `json:"objectWriteTraceId"`
		}{
			{
				ID:         "test@example.com",
				Properties: map[string]string{"firstname": "John"},
			},
		},
	}

	result, err := objectClient.BatchCreateOrUpdateObjects(context.Background(), "contacts", input)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to ubmarshal")
}

// TestBatchArchiveObjects_Success tests successful batch archive
func TestBatchArchiveObjects_Success(t *testing.T) {
	responseJSON := `{
		"completedAt": "2024-01-01T00:00:05.000Z",
		"startedAt": "2024-01-01T00:00:00.000Z",
		"status": "COMPLETE",
		"results": []
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/contacts/batch/archive", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchArchiveObjectsInput{
		Inputs: []struct {
			ID string `json:"id" required:"yes"`
		}{
			{ID: "1"},
			{ID: "2"},
		},
	}

	result, err := objectClient.BatchArchiveObjects(context.Background(), "contacts", input)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, Complete, result.Status)
}

// TestBatchArchiveObjects_WithErrors tests batch archive with errors
func TestBatchArchiveObjects_WithErrors(t *testing.T) {
	responseJSON := `{
		"completedAt": "2024-01-01T00:00:05.000Z",
		"startedAt": "2024-01-01T00:00:00.000Z",
		"status": "COMPLETE",
		"results": [],
		"numErrors": 1,
		"errors": [
			{
				"status": "error",
				"category": "OBJECT_NOT_FOUND",
				"message": "Object not found",
				"context": {},
				"links": {},
				"errors": [
					{
						"message": "Object with ID 99999 not found"
					}
				]
			}
		]
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchArchiveObjectsInput{
		Inputs: []struct {
			ID string `json:"id" required:"yes"`
		}{
			{ID: "1"},
			{ID: "99999"},
		},
	}

	result, err := objectClient.BatchArchiveObjects(context.Background(), "contacts", input)

	require.Error(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.NumErrors)
	assert.Len(t, result.Errors, 1)
}

// TestBatchArchiveObjects_InvalidJSON tests invalid JSON response
func TestBatchArchiveObjects_InvalidJSON(t *testing.T) {
	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	input := &BatchArchiveObjectsInput{
		Inputs: []struct {
			ID string `json:"id" required:"yes"`
		}{
			{ID: "1"},
		},
	}

	result, err := objectClient.BatchArchiveObjects(context.Background(), "contacts", input)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to ubmarshal")
}

// TestSearchObjects_Success tests successful object search
func TestSearchObjects_Success(t *testing.T) {
	responseJSON := `{
		"total": 2,
		"results": [
			{
				"id": "1",
				"properties": {
					"email": "test1@example.com",
					"firstname": "John"
				},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			},
			{
				"id": "2",
				"properties": {
					"email": "test2@example.com",
					"firstname": "Jane"
				},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}
		],
		"paging": {
			"next": {
				"after": "cursor123",
				"link": "?after=cursor123"
			}
		}
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/contacts/search", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &SearchObjectsInput{
		Limit:      10,
		After:      "",
		Sorts:      []string{},
		Properties: []string{"email", "firstname"},
		FilterGroups: []struct {
			Filters []struct {
				PropertyName string         `json:"propertyName" required:"yes"`
				Operator     FilterOperator `json:"operator" required:"yes"`
				HighValue    string         `json:"highValue"`
				Values       []string       `json:"values"`
				Value        string         `json:"value"`
			} `json:"filters" required:"yes"`
		}{
			{
				Filters: []struct {
					PropertyName string         `json:"propertyName" required:"yes"`
					Operator     FilterOperator `json:"operator" required:"yes"`
					HighValue    string         `json:"highValue"`
					Values       []string       `json:"values"`
					Value        string         `json:"value"`
				}{
					{
						PropertyName: "email",
						Operator:     ContainsToken,
						Value:        "test",
					},
				},
			},
		},
	}

	result, err := objectClient.SearchObjects(context.Background(), "contacts", input)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.Total)
	assert.Len(t, result.Results, 2)
}

// TestSearchObjects_WithFilters tests search with complex filters
func TestSearchObjects_WithFilters(t *testing.T) {
	responseJSON := `{
		"total": 1,
		"results": [
			{
				"id": "1",
				"properties": {
					"email": "test@example.com",
					"age": "30"
				},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}
		]
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &SearchObjectsInput{
		Limit:      10,
		After:      "",
		Sorts:      []string{},
		Properties: []string{"email", "age"},
		FilterGroups: []struct {
			Filters []struct {
				PropertyName string         `json:"propertyName" required:"yes"`
				Operator     FilterOperator `json:"operator" required:"yes"`
				HighValue    string         `json:"highValue"`
				Values       []string       `json:"values"`
				Value        string         `json:"value"`
			} `json:"filters" required:"yes"`
		}{
			{
				Filters: []struct {
					PropertyName string         `json:"propertyName" required:"yes"`
					Operator     FilterOperator `json:"operator" required:"yes"`
					HighValue    string         `json:"highValue"`
					Values       []string       `json:"values"`
					Value        string         `json:"value"`
				}{
					{
						PropertyName: "age",
						Operator:     GT,
						Value:        "25",
					},
					{
						PropertyName: "email",
						Operator:     HasProperty,
					},
				},
			},
		},
	}

	result, err := objectClient.SearchObjects(context.Background(), "contacts", input)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.Total)
}

// TestSearchObjects_NoResults tests search with no results
func TestSearchObjects_NoResults(t *testing.T) {
	responseJSON := `{
		"total": 0,
		"results": []
	}`

	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &SearchObjectsInput{
		Limit:      10,
		After:      "",
		Sorts:      []string{},
		Properties: []string{"email"},
		FilterGroups: []struct {
			Filters []struct {
				PropertyName string         `json:"propertyName" required:"yes"`
				Operator     FilterOperator `json:"operator" required:"yes"`
				HighValue    string         `json:"highValue"`
				Values       []string       `json:"values"`
				Value        string         `json:"value"`
			} `json:"filters" required:"yes"`
		}{},
	}

	result, err := objectClient.SearchObjects(context.Background(), "contacts", input)

	require.Error(t, err)
	assert.Nil(t, result)
}

// TestSearchObjects_InvalidJSON tests invalid JSON response
func TestSearchObjects_InvalidJSON(t *testing.T) {
	server, objectClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	input := &SearchObjectsInput{
		Limit:      10,
		After:      "",
		Sorts:      []string{},
		Properties: []string{"email"},
		FilterGroups: []struct {
			Filters []struct {
				PropertyName string         `json:"propertyName" required:"yes"`
				Operator     FilterOperator `json:"operator" required:"yes"`
				HighValue    string         `json:"highValue"`
				Values       []string       `json:"values"`
				Value        string         `json:"value"`
			} `json:"filters" required:"yes"`
		}{},
	}

	result, err := objectClient.SearchObjects(context.Background(), "contacts", input)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}
