package nordigen

import (
	"net/http"

	"github.com/pkg/errors"
	"gromson/nordigen/rest"
)

type nordigenResource[Response any] struct {
	nordigen *Nordigen
	generic  rest.GenericResource[Response]
}

func (r *nordigenResource[Response]) wrap(call func() (*Response, error)) (*Response, error) {
	if err := r.nordigen.ensureAuthenticated(); err != nil {
		return nil, err
	}

	res, err := call()
	if err == nil {
		return res, nil
	}

	var apiErr *rest.ApiError
	if errors.As(err, &apiErr) && apiErr.StatusCode() == http.StatusUnauthorized {
		r.nordigen.unauthenticate()
		if err := r.nordigen.ensureAuthenticated(); err != nil {
			return nil, err
		}

		return call()
	}

	return nil, err
}

func (r *nordigenResource[Response]) wrapList(call func() ([]Response, error)) ([]Response, error) {
	if err := r.nordigen.ensureAuthenticated(); err != nil {
		return nil, err
	}

	res, err := call()
	if err == nil {
		return res, nil
	}

	var apiErr *rest.ApiError
	if errors.As(err, &apiErr) && apiErr.StatusCode() == http.StatusUnauthorized {
		r.nordigen.unauthenticate()
		if err := r.nordigen.ensureAuthenticated(); err != nil {
			return nil, err
		}

		return call()
	}

	return nil, err
}
