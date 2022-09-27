package nordigen

import (
	"net/url"

	"gromson/nordigen/rest"
)

const institutionResourceID = "/institutions"

// InstitutionResource resource for accessing institutions
type InstitutionResource struct {
	nordigenResource[InstitutionResponse]
}

// InstitutionResponse information about an institution
type InstitutionResponse struct {
	ID                   string   `json:"id"`
	Name                 string   `json:"name"`
	BIC                  string   `json:"bic"`
	TransactionTotalDays int      `json:"transaction_total_days,string"`
	Countries            []string `json:"countries"`
	LogoUrl              string   `json:"logo"`
}

// Institution provides access to institutions
func (n *Nordigen) Institution() *InstitutionResource {
	return &InstitutionResource{
		nordigenResource[InstitutionResponse]{
			nordigen: n,
			generic: rest.GenericResource[InstitutionResponse]{
				ID:     institutionResourceID,
				Client: n.restClient,
			},
		},
	}
}

// List returns a list of institutions. If the country (ISO 3166 two-character country code) argument
// is an empty string the result will contain institutions from all available countries.
// In case API HTTP error response rest.ApiError will be returned,
func (r *InstitutionResource) List(country string) ([]InstitutionResponse, error) {
	params := url.Values{}
	if country != "" {
		params.Add("country", country)
	}

	return r.list(params)
}

// ListWithEnabledPayments returns a list of institutions with enabled payments.
// If the country (ISO 3166 two-character country code) argument
// is an empty string the result will contain institutions from all available countries.
// In case API HTTP error response rest.ApiError will be returned,
func (r *InstitutionResource) ListWithEnabledPayments(country string) ([]InstitutionResponse, error) {
	params := url.Values{}
	params.Add("payments_enabled", "true")

	if country != "" {
		params.Add("country", country)
	}

	return r.list(params)
}

// ListWithDisabledPayments returns a list of institutions with disabled payments.
// If the country (ISO 3166 two-character country code) argument
// is an empty string the result will contain institutions from all available countries.
// In case API HTTP error response rest.ApiError will be returned,
func (r *InstitutionResource) ListWithDisabledPayments(country string) ([]InstitutionResponse, error) {
	params := url.Values{}
	params.Add("payments_enabled", "false")

	if country != "" {
		params.Add("country", country)
	}

	return r.list(params)
}

// In case API HTTP error response rest.ApiError will be returned,
func (r *InstitutionResource) list(params url.Values) ([]InstitutionResponse, error) {
	return r.wrapList(
		func() ([]InstitutionResponse, error) {
			return r.generic.List(params)
		},
	)
}

// Get returns an institutions with the given ID.
// In case API HTTP error response rest.ApiError will be returned,
func (r *InstitutionResource) Get(ID string) (*InstitutionResponse, error) {
	return r.wrap(
		func() (*InstitutionResponse, error) {
			return r.generic.Get(ID, nil)
		},
	)
}
