package nordigen

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// ApiClientErrorResponse represents the 400 and 404 error responses from the API
type ApiClientErrorResponse map[string]interface{}

// Error returns string representation of the error
func (e *ApiClientErrorResponse) Error() string {
	msg := make([]string, 0, 1)
	for _, v := range *e {
		if detail, ok := v.(map[string]interface{}); ok {
			msg = append(msg, newApiErrorDetailFromMap(detail).string())
			continue
		}

		if val, ok := v.(string); ok {
			msg = append(msg, val)
			continue
		}
	}

	return strconv.Itoa(e.StatusCode()) + " - " + strings.Join(msg, "; ")
}

// StatusCode returns the HTTP response status code
func (e *ApiClientErrorResponse) StatusCode() int {
	code, ok := (*e)["status_code"]
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

// Fields returns the list of field names that have errors
func (e *ApiClientErrorResponse) Fields() []string {
	fieldNames := make([]string, 0, len(*e))
	cntr := 0
	for name := range *e {
		fieldNames[cntr] = name
		cntr++
	}

	return fieldNames
}

// ApiServerErrorResponse represents HTTP errors responses other than 400 and 404
type ApiServerErrorResponse struct {
	apiErrorDetail
	StatusCode int `json:"status_code"`
}

// Error returns string representation of the error
func (e *ApiServerErrorResponse) Error() string {
	return strconv.Itoa(e.StatusCode) + " - " + e.string()
}

type apiErrorDetail struct {
	Summary string `json:"summary"`
	Detail  string `json:"detail"`
}

func (d *apiErrorDetail) string() string {
	return d.Summary + ": " + d.Detail
}

func newApiErrorDetailFromMap(m map[string]interface{}) *apiErrorDetail {
	summary, ok := m["summary"].(string)
	if !ok {
		summary = ""
	}

	detail, ok := m["detail"].(string)
	if !ok {
		detail = ""
	}

	return &apiErrorDetail{
		Summary: summary,
		Detail:  detail,
	}
}

func createApiServerError(res *http.Response) error {
	return createApiError(res, &ApiServerErrorResponse{})
}

func createApiClientError(res *http.Response) error {
	return createApiError(res, &ApiClientErrorResponse{})
}

func createApiError(res *http.Response, ierr error) error {
	if err := json.NewDecoder(res.Body).Decode(ierr); err != nil {
		return errors.Wrap(err, "error unmarshaling API error response")
	}

	return ierr
}
