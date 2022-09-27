package nordigen

import (
	"time"

	"github.com/google/uuid"
	"gromson/nordigen/rest"
)

const (
	requisitionResourceID = "/requisitions"
)

type RequisitionResource struct {
	nordigenResource[RequisitionResponse]
}

// CreateRequisitionRequest requisition creation request
type CreateRequisitionRequest struct {
	// Redirect URL to your application after end-user authorization with ASPSP
	Redirect string `json:"redirect"`
	// InstitutionID
	InstitutionID string `json:"institution_id"`
	// Agreement ID of the end user agreement
	Agreement uuid.UUID `json:"agreement"`
	// Reference additional ID to identify the end user
	Reference string `json:"reference"`
	// UserLanguage a two-letter country code (ISO 639-1)
	UserLanguage string `json:"user_language"`
	// Ssn optional field to verify ownership of the account
	Ssn string `json:"ssn"`
	// AccountSelection option to enable account selection view for the end user
	AccountSelection bool `json:"account_selection"`
	// RedirectImmediate enable redirect back to the client after account list received
	RedirectImmediate bool `json:"redirect_immediate"`
}

// RequisitionResponse information about requisition
type RequisitionResponse struct {
	ID uuid.UUID `json:"id"`

	Created time.Time `json:"created"`

	// RedirectUrl to your application after end-user authorization with ASPSP
	RedirectUrl string `json:"redirect"`

	Status string `json:"status"`

	InstitutionID string `json:"institution_id"`

	AgreementID uuid.UUID `json:"agreement"`

	// Reference additional ID to identify the end user
	Reference string `json:"reference"`

	Accounts []uuid.UUID `json:"accounts"`

	UserLanguage string `json:"user_language"`

	// Link to initiate authorization with Institution
	Link string `json:"link"`

	// Ssn optional field to verify ownership of the account
	Ssn string `json:"ssn"`

	// AccountSelection option to enable account selection view for the end user
	AccountSelection bool `json:"account_selection"`

	// RedirectImmediate enable redirect back to the client after account list received
	RedirectImmediate bool `json:"redirect_immediate"`
}

type RequisitionCollectionResponse = CollectionResponse[RequisitionResponse]

// Requisition access to requisition resource
func (n *Nordigen) Requisition() *RequisitionResource {
	return &RequisitionResource{
		nordigenResource[RequisitionResponse]{
			nordigen: n,
			generic: rest.GenericResource[RequisitionResponse]{
				ID:     requisitionResourceID,
				Client: n.restClient,
			},
		},
	}
}

// Get a requisition with a given ID
// In case API HTTP error response rest.ApiError will be returned,
func (r *RequisitionResource) Get(ID uuid.UUID) (*RequisitionResponse, error) {
	return r.wrap(
		func() (*RequisitionResponse, error) {
			return r.generic.Get(ID.String(), nil)
		},
	)
}

// List returns RequisitionResource for access to the list of requisitions.
// In case API HTTP error response rest.ApiError will be returned,
func (r *RequisitionResource) List() (*RequisitionCollectionResponse, error) {
	return newCollectionResponse(r.nordigen, &r.generic, nil)
}

// Create a new requisition.
// In case API HTTP error response rest.ApiError will be returned,
func (r *RequisitionResource) Create(payload *CreateRequisitionRequest) (*RequisitionResponse, error) {
	return r.wrap(
		func() (*RequisitionResponse, error) {
			return r.generic.Post("", payload)
		},
	)
}

// Delete a requisition with the given ID
// In case API HTTP error response rest.ApiError will be returned,
func (r *RequisitionResource) Delete(ID uuid.UUID) error {
	_, err := r.wrap(
		func() (*RequisitionResponse, error) {
			return nil, r.generic.Delete(ID.String())
		},
	)

	return err
}
