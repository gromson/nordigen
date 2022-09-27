package nordigen

import (
	"github.com/google/uuid"
	"gromson/nordigen/rest"
)

const (
	accountDetailsResourceID = "/details"
)

// AccountDetailsResource access to account details
type AccountDetailsResource struct {
	nordigenResource[AccountDetailsResponse]
}

// AccountDetailsResponse API response structure
type AccountDetailsResponse struct {
	Account AccountDetailsInfoResponse `json:"account"`
}

// AccountDetailsInfoResponse account details information
type AccountDetailsInfoResponse struct {
	// ResourceID shall be filled, if addressable resource are created by the ASPSP
	ResourceID uuid.UUID `json:"resourceId"`

	Iban string `json:"iban"`

	// Currency ISO 4217 Alpha 3 currency code
	Currency string `json:"currency"`

	// Name of the legal account owner.
	// If there is more than one owner, then e.g. two names might be noted here.
	// For a corporate account, the corporate name is used for this attribute.
	// Even if supported by the ASPSP, the provision of this field might depend on the fact whether an explicit
	// consent to this specific additional account information has been given by the PSU.
	OwnerName string `json:"ownerName"`

	// Name of the account, as assigned by the ASPSP, in agreement with the account owner
	// in order to provide an additional means of identification of the account
	Name string `json:"name"`

	// Product name of the bank for this account, proprietary definition
	Product string `json:"product"`

	// CashAccountType ExternalCashAccountType1Code from ISO 20022
	CashAccountType string `json:"cashAccountType"`

	// Status undocumented. One of the possible values "enabled"
	Status string `json:"status"`

	// Usage specifies the usage of the account:
	// * PRIV: private personal account
	// * ORGA: professional account
	Usage string `json:"usage"`
}

// Details returns the resource to access account's balance data
func (r *AccountResource) Details(accountID uuid.UUID) *AccountDetailsResource {
	resourceID := accountResourceID + "/" + accountID.String() + accountDetailsResourceID

	return &AccountDetailsResource{
		nordigenResource[AccountDetailsResponse]{
			r.nordigen,
			rest.GenericResource[AccountDetailsResponse]{
				ID:     resourceID,
				Client: r.nordigen.restClient,
			},
		},
	}
}

// Get details for the underlying account resource
// In case API HTTP error response rest.ApiError will be returned,
func (d *AccountDetailsResource) Get() (*AccountDetailsResponse, error) {
	return d.wrap(
		func() (*AccountDetailsResponse, error) {
			return d.generic.Get("", nil)
		},
	)
}
