package rest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"gromson/nordigen/utils"
)

var defaultHttpHeader = http.Header{"Content-Type": []string{"application/json"}}

// Client for accessing REST API
type Client struct {
	BaseUrl  *url.URL
	Header   http.Header
	LogError func(err error, message string)

	httpClient *http.Client
}

// NewClient creates new REST API client.
func NewClient(baseUrl *url.URL, header http.Header) *Client {
	return &Client{
		BaseUrl:  baseUrl,
		Header:   utils.MergeMapsOfArrays(defaultHttpHeader, header),
		LogError: defaultLogError,
	}
}

// Exec executes an HTTP request with a given body payload and writes the response body to the target
// ApiError in case of API HTTP error response. In case of the error unrelated to API other error type will be returned
func (c *Client) Exec(method, resourceID string, body io.Reader, target interface{}) error {
	req, err := c.NewRequest(method, resourceID, body)
	if err != nil {
		return errors.Wrap(err, "error creating an HTTP request")
	}

	return c.ExecuteRequest(req, target)
}

// NewRequest returns an HTTP request to the resource
func (c *Client) NewRequest(method, resourceID string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.BaseUrl.String()+resourceID, body)
	if err != nil {
		return nil, err
	}

	req.Header = utils.MergeMapsOfArrays(req.Header, c.Header)

	return req, nil
}

// ExecuteRequest executes the HTTP request and populates the result in case of a successful response or returns
// ApiError in case of API HTTP error response. In case of the error unrelated to API other error type will be returned
func (c *Client) ExecuteRequest(req *http.Request, target interface{}) error {
	res, err := c.doRequest(req)
	if err != nil {
		return errors.Wrap(err, "error while executing HTTP request")
	}
	defer c.execAndLogIfErr(res.Body.Close, "error while closing response body")

	if target == nil {
		return nil
	}

	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return errors.Wrap(err, "error evaluating a response body")
	}

	return nil
}

func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	res, err := c.http().Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error while trying to execute a request")
	}

	if res.StatusCode == http.StatusOK {
		return res, nil
	}

	if res.StatusCode != http.StatusOK {
		defer c.execAndLogIfErr(res.Body.Close, "error while closing response body")
	}

	return nil, createApiError(res)
}

func (c *Client) http() *http.Client {
	if c.httpClient == nil {
		c.httpClient = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	return c.httpClient
}

func (c *Client) execAndLogIfErr(callback func() error, message string) {
	err := callback()
	if err != nil && c.LogError != nil {
		c.LogError(err, message)
	}
}

func defaultLogError(err error, message string) {
	log.Printf("%s: %s", message, err)
}
