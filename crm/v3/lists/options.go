package lists

import (
	"fmt"

	"github.com/aacc-dev/go-hubspot-sdk/client"
)

// GetListOption is a functional option for GetList methods (GetListById, GetListByName, GetListsByIDs)
type GetListOption func(*client.Request)

// WithIncludeFilters includes filter definitions in the response
func WithIncludeFilters(includeFilters bool) GetListOption {
	return func(req *client.Request) {
		if includeFilters {
			req.AddQueryParam("includeFilters", "true")
		}
	}
}

// ListMembershipsOption is a functional option for GetListMemberships
type ListMembershipsOption func(*client.Request)

// WithMembershipsLimit sets the maximum number of membership results to return
func WithMembershipsLimit(limit int) ListMembershipsOption {
	return func(req *client.Request) {
		req.AddQueryParam("limit", fmt.Sprintf("%d", limit))
	}
}

// WithMembershipsOffset sets the pagination offset for list memberships
func WithMembershipsOffset(offset string) ListMembershipsOption {
	return func(req *client.Request) {
		req.AddQueryParam("offset", offset)
	}
}
