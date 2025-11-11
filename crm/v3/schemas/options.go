package schemas

import "github.com/josiah-hester/go-hubspot-sdk/client"

// SchemaOption is a functional option for the Schemas API
type SchemaOption func(*client.Request)

// WithArchived includes archived schemas in results
func WithArchived() SchemaOption {
	return func(req *client.Request) {
		req.AddQueryParam("archived", "true")
	}
}
