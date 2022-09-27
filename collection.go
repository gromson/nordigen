package nordigen

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
	"gromson/nordigen/rest"
	"gromson/nordigen/utils"
)

const (
	defaultListResponseLimit = 1000
)

// CollectionResponse represents a response with a list of items
type CollectionResponse[Response any] struct {
	nordigen *Nordigen
	count    int
	next     int
	limit    int
	offset   int
	results  []Response
	resource *rest.GenericResource[Response]
	params   url.Values
}

func newCollectionResponse[Response any](n *Nordigen, r *rest.GenericResource[Response], params url.Values) (*CollectionResponse[Response], error) {
	collection := &CollectionResponse[Response]{
		nordigen: n,
		count:    0,
		next:     0,
		limit:    defaultListResponseLimit,
		offset:   0,
		results:  make([]Response, 0, 2),
		resource: r,
		params:   params,
	}

	if err := collection.get(); err != nil {
		return nil, err
	}

	return collection, nil
}

// Next return the next item from the collection
// In case of API error rest.ApiError returned
func (c *CollectionResponse[Response]) Next() (*Response, error) {
	if c.resource == nil || c.resource.Client == nil {
		return nil, nil
	}

	if c.next == c.count {
		return nil, nil
	}

	if len(c.results) == 0 || c.next >= c.offset {
		if err := c.get(); err != nil {
			return nil, errors.Wrap(err, "error getting next item")
		}
	}

	defer func() { c.next++ }()

	return &c.results[c.next], nil
}

// Count returns the number of existing end user agreements
func (c *CollectionResponse[Response]) Count() int {
	return c.count
}

func (c *CollectionResponse[Response]) get() error {
	if err := c.nordigen.ensureAuthenticated(); err != nil {
		return err
	}

	err := c.exec()
	if err == nil {
		return nil
	}

	var apiErr *rest.ApiError
	if errors.As(err, &apiErr) && apiErr.StatusCode() == http.StatusUnauthorized {
		c.nordigen.unauthenticate()
		if err := c.nordigen.ensureAuthenticated(); err != nil {
			return err
		}

		return c.exec()
	}

	return err
}

func (c *CollectionResponse[Response]) exec() error {
	collectionParams := url.Values{}
	collectionParams.Add("limit", strconv.Itoa(c.limit))
	collectionParams.Add("offset", strconv.Itoa(c.offset))

	if c.params == nil {
		c.params = url.Values{}
	}

	p := utils.MergeMapsOfArrays(c.params, collectionParams)

	path := c.resource.ID + "?" + p.Encode()

	results := struct {
		Count   int        `json:"count"`
		Results []Response `json:"results"`
	}{}

	if err := c.resource.Client.Exec(http.MethodGet, path, nil, &results); err != nil {
		return errors.Wrap(err, "error evaluating collection response")
	}

	if len(c.results) != results.Count {
		tmpResults := make([]Response, len(c.results))
		copy(tmpResults, c.results)
		c.results = make([]Response, results.Count)
		copy(c.results, tmpResults)
	}

	copy(c.results[c.offset:], results.Results)

	c.offset += c.limit
	c.count = results.Count

	return nil
}
