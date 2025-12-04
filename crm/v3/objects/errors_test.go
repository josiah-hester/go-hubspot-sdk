package objects

import (
	"errors"
	"testing"

	"github.com/josiah-hester/go-hubspot-sdk/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestObjectNotFoundError_Error tests the Error() method
func TestObjectNotFoundError_Error(t *testing.T) {
	err := &ObjectNotFoundError{
		ObjectType: "contacts",
		Original: &client.HubSpotError{
			Status:  404,
			Message: "Not found",
		},
	}

	expectedMsg := "object contacts not found"
	assert.Equal(t, expectedMsg, err.Error())
}

// TestObjectValidationError_Error tests the Error() method
func TestObjectValidationError_Error(t *testing.T) {
	err := &ObjectValidationError{
		Field:   "email",
		Message: "Invalid email format",
		Original: &client.HubSpotError{
			Status:   400,
			Category: "VALIDATION_ERROR",
		},
	}

	expectedMsg := "validation error on field email: Invalid email format"
	assert.Equal(t, expectedMsg, err.Error())
}

// TestObjectValidationError_Error_EmptyField tests Error() with empty field
func TestObjectValidationError_Error_EmptyField(t *testing.T) {
	err := &ObjectValidationError{
		Field:   "",
		Message: "Validation failed",
		Original: &client.HubSpotError{
			Status:   400,
			Category: "VALIDATION_ERROR",
		},
	}

	expectedMsg := "validation error on field : Validation failed"
	assert.Equal(t, expectedMsg, err.Error())
}

// TestObjectAlreadyExistsError_Error tests the Error() method
func TestObjectAlreadyExistsError_Error(t *testing.T) {
	err := &ObjectAlreadyExistsError{
		ObjectID: "12345",
		Original: &client.HubSpotError{
			Status:  409,
			Message: "Conflict",
		},
	}

	expectedMsg := "object with id 12345 already exists"
	assert.Equal(t, expectedMsg, err.Error())
}

// TestParseObjectError_NotFound tests parsing 404 errors
func TestParseObjectError_NotFound(t *testing.T) {
	hubspotErr := &client.HubSpotError{
		Status:  404,
		Message: "Object not found",
	}

	result := ParseObjectError(hubspotErr, "contacts")

	var notFoundErr *ObjectNotFoundError
	require.ErrorAs(t, result, &notFoundErr)
	assert.Equal(t, "contacts", notFoundErr.ObjectType)
	assert.Equal(t, hubspotErr, notFoundErr.Original)
}

// TestParseObjectError_ValidationError tests parsing 400 validation errors
func TestParseObjectError_ValidationError(t *testing.T) {
	hubspotErr := &client.HubSpotError{
		Status:   400,
		Category: "VALIDATION_ERROR",
		Message:  "Invalid email format",
	}

	result := ParseObjectError(hubspotErr, "contacts")

	var validationErr *ObjectValidationError
	require.ErrorAs(t, result, &validationErr)
	assert.Equal(t, "Invalid email format", validationErr.Field)
	assert.Equal(t, hubspotErr, validationErr.Original)
}

// TestParseObjectError_BadRequestNonValidation tests parsing 400 non-validation errors
func TestParseObjectError_BadRequestNonValidation(t *testing.T) {
	hubspotErr := &client.HubSpotError{
		Status:   400,
		Category: "OTHER_ERROR",
		Message:  "Bad request",
	}

	result := ParseObjectError(hubspotErr, "contacts")

	// Should return the original error unchanged since it's not a validation error
	assert.Equal(t, hubspotErr, result)
}

// TestParseObjectError_Conflict tests parsing 409 conflict errors
func TestParseObjectError_Conflict(t *testing.T) {
	hubspotErr := &client.HubSpotError{
		Status:  409,
		Message: "Object already exists",
	}

	result := ParseObjectError(hubspotErr, "contacts")

	var existsErr *ObjectAlreadyExistsError
	require.ErrorAs(t, result, &existsErr)
	assert.Equal(t, hubspotErr, existsErr.Original)
}

// TestParseObjectError_OtherHubSpotError tests parsing other HubSpot errors
func TestParseObjectError_OtherHubSpotError(t *testing.T) {
	hubspotErr := &client.HubSpotError{
		Status:  500,
		Message: "Internal server error",
	}

	result := ParseObjectError(hubspotErr, "contacts")

	// Should return the original error unchanged for non-mapped status codes
	assert.Equal(t, hubspotErr, result)
}

// TestParseObjectError_NonHubSpotError tests parsing non-HubSpot errors
func TestParseObjectError_NonHubSpotError(t *testing.T) {
	regularErr := errors.New("some other error")

	result := ParseObjectError(regularErr, "contacts")

	// Should return the original error unchanged
	assert.Equal(t, regularErr, result)
}

// TestBatchError_Error_WithMessage tests BatchError.Error() with message only
func TestBatchError_Error_WithMessage(t *testing.T) {
	err := &BatchError{
		Message:  "Batch operation failed",
		Category: "BATCH_ERROR",
		Status:   "error",
		Context:  make(map[string][]string),
		Links:    make(map[string]string),
		Errors:   []ObjectError{},
	}

	expectedMsg := "Batch operation failed"
	assert.Equal(t, expectedMsg, err.Error())
}

// TestBatchError_Error_WithErrors tests BatchError.Error() with errors
func TestBatchError_Error_WithErrors(t *testing.T) {
	err := &BatchError{
		Message:  "Batch operation failed",
		Category: "BATCH_ERROR",
		Status:   "error",
		Context:  make(map[string][]string),
		Links:    make(map[string]string),
		Errors: []ObjectError{
			{
				Message: "Object 1 not found",
			},
			{
				Message: "Object 2 validation error",
			},
		},
	}

	expectedMsg := "Batch operation failed: Object 1 not found; Object 2 validation error"
	assert.Equal(t, expectedMsg, err.Error())
}

// TestBatchError_Error_WithoutMessage tests BatchError.Error() without message
func TestBatchError_Error_WithoutMessage(t *testing.T) {
	err := &BatchError{
		Message:  "",
		Category: "BATCH_ERROR",
		Status:   "error",
		Context:  make(map[string][]string),
		Links:    make(map[string]string),
		Errors: []ObjectError{
			{
				Message: "Object 1 not found",
			},
		},
	}

	expectedMsg := "batch error: Object 1 not found"
	assert.Equal(t, expectedMsg, err.Error())
}

// TestBatchError_Error_WithMultipleErrors tests BatchError.Error() with multiple errors
func TestBatchError_Error_WithMultipleErrors(t *testing.T) {
	err := &BatchError{
		Message:  "Multiple errors occurred",
		Category: "BATCH_ERROR",
		Status:   "error",
		Context:  make(map[string][]string),
		Links:    make(map[string]string),
		Errors: []ObjectError{
			{
				Message:     "First error",
				SubCategory: "NOT_FOUND",
			},
			{
				Message:     "Second error",
				SubCategory: "VALIDATION",
			},
			{
				Message:     "Third error",
				SubCategory: "CONFLICT",
			},
		},
	}

	result := err.Error()
	assert.Contains(t, result, "Multiple errors occurred")
	assert.Contains(t, result, "First error")
	assert.Contains(t, result, "Second error")
	assert.Contains(t, result, "Third error")
}

// TestBatchError_Error_EmptyMessageNoErrors tests BatchError.Error() with no message and no errors
func TestBatchError_Error_EmptyMessageNoErrors(t *testing.T) {
	err := &BatchError{
		Message:  "",
		Category: "BATCH_ERROR",
		Status:   "error",
		Context:  make(map[string][]string),
		Links:    make(map[string]string),
		Errors:   []ObjectError{},
	}

	expectedMsg := "batch error"
	assert.Equal(t, expectedMsg, err.Error())
}
