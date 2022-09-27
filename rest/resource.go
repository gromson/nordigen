package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type GenericResource[Response any] struct {
	// REST resource ID (in other words URL path e.g. /users)
	ID     string
	Client *Client
}

// Post sends POST request to the resource.
// In case API returns HTTP error response ApiError will be returned
func (r *GenericResource[Response]) Post(ID string, payload interface{}) (*Response, error) {
	return r.exec(http.MethodPost, ID, nil, payload)
}

// Put sends PUT request to the resource.
// In case API returns HTTP error response ApiError will be returned
func (r *GenericResource[Response]) Put(ID string, payload interface{}) (*Response, error) {
	return r.exec(http.MethodPut, ID, nil, payload)
}

// List returns GenericCollectionResponse for access to the List of requisitions.
// In case API returns HTTP error response ApiError will be returned
func (r *GenericResource[Response]) List(params url.Values) ([]Response, error) {
	path := r.preparePath("", params)

	res := make([]Response, 0, 2)
	if err := r.Client.Exec(http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// Get a resource with a given ID or without ID.
// In case API returns HTTP error response ApiError will be returned
func (r *GenericResource[Response]) Get(ID string, params url.Values) (*Response, error) {
	return r.exec(http.MethodGet, ID, params, nil)
}

// Delete a resource with the given ID or without ID
// In case API returns HTTP error response ApiError will be returned
func (r *GenericResource[Response]) Delete(ID string) error {
	return r.Client.Exec(http.MethodDelete, resourceID(r.ID, ID), nil, nil)
}

func (r *GenericResource[Response]) exec(method, ID string, params url.Values, payload interface{}) (*Response, error) {
	path := r.preparePath(ID, params)
	body, err := r.prepareBody(payload)
	if err != nil {
		return nil, err
	}

	res := new(Response)
	if err := r.Client.Exec(method, path, body, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *GenericResource[Response]) preparePath(ID string, params url.Values) string {
	path := resourceID(r.ID, ID)
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	return path
}

func (r *GenericResource[Response]) prepareBody(payload interface{}) (io.Reader, error) {
	if payload == nil {
		return nil, nil
	}

	body := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return nil, errors.Wrap(err, "error marshaling a resource POST request")
	}

	return body, nil
}

func resourceID(resourceRoot, ID string) string {
	if ID != "" {
		resourceRoot += "/" + ID
	}

	return resourceRoot
}
