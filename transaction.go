package nordigen

import (
	"net/url"
	"time"

	"github.com/google/uuid"
	"gromson/nordigen/rest"
)

const (
	dateFormat = "2006-01-02"
)

const (
	transactionsResourceID = "/transactions"
)

// TransactionResource access to account's balances
type TransactionResource struct {
	nordigenResource[TransactionCollectionResponse]
}

// TransactionCollectionResponse API response structure
type TransactionCollectionResponse struct {
	Transactions TransactionTypesResponse `json:"transactions"`
}

// TransactionTypesResponse structure representing part of the API response
type TransactionTypesResponse struct {
	Booked      []TransactionResponse `json:"booked"`
	Pending     []TransactionResponse `json:"pending"`
	Information []TransactionResponse `json:"information"`
}

// TransactionResponse transaction information
type TransactionResponse struct {
	ID uuid.UUID `json:"transactionId"`

	Amount Amount `json:"transactionAmount"`

	// BookingDate the date when an entry is posted to an account on the ASPSPs books
	BookingDate string `json:"bookingDate"`

	// ValueDate the date at which assets become available to the account owner in case of a credit,
	// or cease to be available to the account owner in case of a debit entry.
	// **Usage:** If entry status is pending and value date is present, then the value date refers to
	// an expected/requested value date.
	ValueDate string `json:"valueDate"`

	// MandateID Identification of Mandates, e.g. a SEPA Mandate ID
	MandateID string `json:"mandateId"`

	// CreditorId Identification of Creditors, e.g. a SEPA Creditor ID
	CreditorId string `json:"creditorId"`

	CreditorName string `json:"creditorName"`

	CreditorAccount *AccountReference `json:"creditorAccount"`

	DebtorName string `json:"debtorName"`

	DebtorAccount *AccountReference `json:"debtorAccount"`

	RemittanceInformationUnstructured string `json:"remittanceInformationUnstructured"`

	RemittanceInformationUnstructuredArray []string `json:"remittanceInformationUnstructuredArray"`

	// AdditionalInformation might be used by the ASPSP to transport additional transaction related information
	// to the PSU
	AdditionalInformation string `json:"additionalInformation"`

	// Proprietary bank transaction code as used within a community or within an ASPSP
	// e.g. for MT94x based transaction reports
	BankTransactionCode string `json:"bankTransactionCode"`
}

// AccountReference reference to an account by either
// * IBAN, of a payment accounts, or
// * BBAN, for payment accounts if there is no IBAN, or
// * the Primary Account Number (PAN) of a card, can be tokenised by the ASPSP due to PCI DSS requirements, or
// * the Primary Account Number (PAN) of a card in a masked form, or
// * an alias to access a payment account via a registered mobile phone number (MSISDN), or
// * a proprietary ID of the respective account that uniquely identifies the account for this ASPSP.
type AccountReference struct {
	Iban string `json:"iban"`

	Bban string `json:"bban"`

	// Pan Primary Account Number
	Pan string `json:"pan"`

	// MaskedPan Primary Account Number in a masked form
	MaskedPan string `json:"maskedPan"`

	// MSISDN registered mobile phone number
	MSISDN string `json:"msisdn"`

	Other string `json:"other"`

	// Currency ISO 4217 Alpha 3 currency code
	Currency string `json:"currency"`
}

// Transaction returns the resource to access account's balance data
func (r *AccountResource) Transaction(accountID uuid.UUID) *TransactionResource {
	resourceID := accountResourceID + "/" + accountID.String() + transactionsResourceID

	return &TransactionResource{
		nordigenResource[TransactionCollectionResponse]{
			nordigen: r.nordigen,
			generic: rest.GenericResource[TransactionCollectionResponse]{
				ID:     resourceID,
				Client: r.nordigen.restClient,
			},
		},
	}
}

// Get transactions for the underlying account resource
// In case API HTTP error response rest.ApiError will be returned,
func (tr *TransactionResource) Get(dateFrom *time.Time, dateTo *time.Time) (*TransactionCollectionResponse, error) {
	return tr.wrap(
		func() (*TransactionCollectionResponse, error) {
			params := url.Values{}

			if dateFrom != nil {
				params.Add("date_from", dateFrom.Format(dateFormat))
			}
			if dateTo != nil {
				params.Add("date_to", dateTo.Format(dateFormat))
			}

			return tr.generic.Get("", params)
		},
	)
}
