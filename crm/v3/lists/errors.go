package lists

import (
	"fmt"

	"github.com/aacc-dev/go-hubspot-sdk/client"
)

// ListNotFoundError is returned when a list is not found
type ListNotFoundError struct {
	ListID   string
	Original *client.HubSpotError
}

func (e *ListNotFoundError) Error() string {
	return fmt.Sprintf("list %s not found", e.ListID)
}

// ListValidationError is returned on validation failures
type ListValidationError struct {
	Field    string
	Message  string
	Original *client.HubSpotError
}

func (e *ListValidationError) Error() string {
	return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}

// ListAlreadyExistsError is returned when trying to create a duplicate
type ListAlreadyExistsError struct {
	ListName string
	Original *client.HubSpotError
}

func (e *ListAlreadyExistsError) Error() string {
	return fmt.Sprintf("list with name %s already exists", e.ListName)
}

// RecordNotFoundError is returned when a record is not found in a list
type RecordNotFoundError struct {
	RecordID string
	ListID   string
	Original *client.HubSpotError
}

func (e *RecordNotFoundError) Error() string {
	return fmt.Sprintf("record %s not found in list %s", e.RecordID, e.ListID)
}

// ParseListError converts a generic HubSpot error to a list-specific error
func ParseListError(err error, listID string) error {
	if hubspotErr, ok := err.(*client.HubSpotError); ok {
		switch hubspotErr.Status {
		case 404:
			return &ListNotFoundError{
				ListID:   listID,
				Original: hubspotErr,
			}
		case 400:
			if hubspotErr.Category == "VALIDATION_ERROR" {
				return &ListValidationError{
					Field:    hubspotErr.Message,
					Original: hubspotErr,
				}
			}
		case 409:
			return &ListAlreadyExistsError{
				Original: hubspotErr,
			}
		}
	}
	return err
}

// ParseRecordError converts a generic HubSpot error to a record-specific error
func ParseRecordError(err error, recordID, listID string) error {
	if hubspotErr, ok := err.(*client.HubSpotError); ok {
		if hubspotErr.Status == 404 {
			return &RecordNotFoundError{
				RecordID: recordID,
				ListID:   listID,
				Original: hubspotErr,
			}
		}
	}
	return ParseListError(err, listID)
}
