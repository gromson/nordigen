package nordigen

import (
	"time"

	"github.com/google/uuid"
	"gromson/nordigen/rest"
)

const (
	balanceResourceID = "/balances"
)

// BalanceResource access to account's balances
type BalanceResource struct {
	nordigenResource[BalanceCollectionResponse]
}

// BalanceCollectionResponse API response structure
type BalanceCollectionResponse struct {
	Balances []BalanceResponse `json:"balances"`
}

// BalanceResponse balance information
type BalanceResponse struct {
	BalanceAmount      Amount    `json:"balanceAmount"`
	BalanceType        string    `json:"balanceType"`
	LastChangeDateTime time.Time `json:"lastChangeDateTime"`
	ReferenceDate      string    `json:"referenceDate"`
}

// Amount of money and currency
type Amount struct {
	Amount string `json:"amount"`

	// Currency ISO 4217 Alpha 3 currency code
	Currency string `json:"currency"`
}

// Balance returns the resource to access account's balance data
func (r *AccountResource) Balance(accountID uuid.UUID) *BalanceResource {
	resourceID := accountResourceID + "/" + accountID.String() + balanceResourceID

	return &BalanceResource{
		nordigenResource[BalanceCollectionResponse]{
			r.nordigen,
			rest.GenericResource[BalanceCollectionResponse]{
				ID:     resourceID,
				Client: r.nordigen.restClient,
			},
		},
	}
}

// Get balances for the underlying account resource
// In case API HTTP error response rest.ApiError will be returned,
func (b *BalanceResource) Get() (*BalanceCollectionResponse, error) {
	return b.wrap(
		func() (*BalanceCollectionResponse, error) {
			return b.generic.Get("", nil)
		},
	)
}
