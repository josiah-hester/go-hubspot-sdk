package associations

import (
	"fmt"

	"github.com/josiah-hester/go-hubspot-sdk/client"
)

// AssociationOption represents a functional option for association requests
type AssociationOption func(*client.Request)

// WithLimit sets the maximum number of results per page
func WithLimit(limit int) AssociationOption {
	return func(req *client.Request) {
		req.AddQueryParam("limit", fmt.Sprintf("%d", limit))
	}
}

// WithAfter sets the paging cursor
func WithAfter(after string) AssociationOption {
	return func(req *client.Request) {
		req.AddQueryParam("after", after)
	}
}
