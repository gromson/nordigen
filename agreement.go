package nordigen

import (
	"net"
	"time"

	"github.com/google/uuid"
	"gromson/nordigen/rest"
)

const (
	endUserAgreementResourceID = "/agreements/enduser"
)

type EndUserAgreementResource struct {
	nordigenResource[EndUserAgreementResponse]
}

// CreateAgreementRequest the request for creating an end user agreement
type CreateAgreementRequest struct {
	InstitutionID      string   `json:"institution_id"`
	MaxHistoricalDays  int      `json:"max_historical_days,string"`
	AccessValidForDays int      `json:"access_valid_for_days"`
	AccessScope        []string `json:"access_scope"`
}

// AcceptEndUserAgreementRequest the request for accepting the end user agreement
type AcceptEndUserAgreementRequest struct {
	UserAgent string `json:"user_agent"`
	IPAddress net.IP `json:"ip_address"`
}

// EndUserAgreementResponse the information about the end user agreement
type EndUserAgreementResponse struct {
	ID                 uuid.UUID  `json:"id"`
	Created            time.Time  `json:"created"`
	MaxHistoricalDays  int        `json:"max_historical_days"`
	AccessValidForDays int        `json:"access_valid_for_days"`
	AccessScopes       []string   `json:"access_scope"`
	Accepted           *time.Time `json:"accepted"`
	InstitutionID      string     `json:"institution_id"`
}

type EndUserAgreementCollectionResponse = CollectionResponse[EndUserAgreementResponse]

// EndUserAgreement access to end user agreement resource
func (n *Nordigen) EndUserAgreement() *EndUserAgreementResource {
	return &EndUserAgreementResource{
		nordigenResource[EndUserAgreementResponse]{
			n,
			rest.GenericResource[EndUserAgreementResponse]{
				ID:     endUserAgreementResourceID,
				Client: n.restClient,
			},
		},
	}
}

// Get an end user agreement with a given ID
// In case API HTTP error response rest.ApiError will be returned,
func (r *EndUserAgreementResource) Get(ID uuid.UUID) (*EndUserAgreementResponse, error) {
	return r.wrap(
		func() (*EndUserAgreementResponse, error) {
			return r.generic.Get(ID.String(), nil)
		},
	)
}

// List returns RequisitionResource for access to the list of requisitions.
// In case API HTTP error response rest.ApiError will be returned,
func (r *EndUserAgreementResource) List() (*EndUserAgreementCollectionResponse, error) {
	return newCollectionResponse(r.nordigen, &r.generic, nil)
}

// Create a new end user agreement.
// In case API HTTP error response rest.ApiError will be returned,
func (r *EndUserAgreementResource) Create(payload *CreateAgreementRequest) (*EndUserAgreementResponse, error) {
	return r.wrap(
		func() (*EndUserAgreementResponse, error) {
			return r.generic.Post("", payload)
		},
	)
}

// Accept a new end user agreement.
// In case API HTTP error response rest.ApiError will be returned,
func (r *EndUserAgreementResource) Accept(ID uuid.UUID, payload *AcceptEndUserAgreementRequest) (*EndUserAgreementResponse, error) {
	return r.wrap(
		func() (*EndUserAgreementResponse, error) {
			return r.generic.Put(ID.String()+"/accept", payload)
		},
	)
}

// Delete an end user agreement with the given ID
// In case API HTTP error response rest.ApiError will be returned,
func (r *EndUserAgreementResource) Delete(ID uuid.UUID) error {
	_, err := r.wrap(
		func() (*EndUserAgreementResponse, error) {
			return nil, r.generic.Delete(ID.String())
		},
	)

	return err
}
