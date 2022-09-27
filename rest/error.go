package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	statusCodeKey = "status_code"
)

// ApiError represents the HTTP error responses from the API
type ApiError map[string]interface{}

func createApiError(res *http.Response) error {
	ierr := &ApiError{
		statusCodeKey: res.StatusCode,
	}

	if err := json.NewDecoder(res.Body).Decode(ierr); err != nil {
		return errors.Wrap(err, "error unmarshaling API error response")
	}

	return ierr
}

// Error returns string representation of the error
func (e *ApiError) Error() string {
	b := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(b).Encode(e); err != nil {
		return "could not interpret ApiError"
	}

	return strconv.Itoa(e.StatusCode()) + " - " + strings.Trim(b.String(), "\n")
}

// StatusCode returns the HTTP response status code
func (e *ApiError) StatusCode() int {
	code, ok := (*e)[statusCodeKey]
	if !ok {
		return 0
	}

	if statusCode, ok := code.(float64); ok {
		return int(statusCode)
	}

	if statusCode, ok := code.(int); ok {
		return statusCode
	}

	return 0
}
