package nordigen

import (
	"time"

	"github.com/google/uuid"
	"gromson/nordigen/rest"
)

const (
	accountResourceID = "/accounts"
)

// AccountResource access to account related information
type AccountResource struct {
	nordigenResource[AccountResponse]
}

// AccountResponse basic account information
type AccountResponse struct {
	ID            uuid.UUID `json:"id"`
	Created       time.Time `json:"created"`
	LastAccessed  time.Time `json:"last_accessed"`
	Iban          string    `json:"iban"`
	InstitutionID string    `json:"institution_id"`
	Status        string    `json:"status"`
}

// Account returns the resource to access account related data
func (n *Nordigen) Account() *AccountResource {
	return &AccountResource{
		nordigenResource[AccountResponse]{
			n,
			rest.GenericResource[AccountResponse]{
				ID:     accountResourceID,
				Client: n.restClient,
			},
		},
	}
}

// Get an account with a given ID
// In case API HTTP error response rest.ApiError will be returned,
func (r *AccountResource) Get(ID uuid.UUID) (*AccountResponse, error) {
	return r.wrap(
		func() (*AccountResponse, error) {
			return r.generic.Get(ID.String(), nil)
		},
	)
}
