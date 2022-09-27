package nordigen

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTransactionResource_Get(t *testing.T) {
	t.Parallel()
	t.Run("getting transactions", testTransactionsGetOk)
	t.Run("getting transactions API error", testApiErrorResponse(testTransactionsGetApiError))
}

func testTransactionsGetOk(t *testing.T) {
	// What/Arrange
	responsePayload := `
		{
			"transactions": {
				"booked": [
					{
						"transactionId": "06de5c3d-aecd-4e58-9d5f-797c7c8a16e8",
						"bookingDate": "2022-09-18",
						"valueDate": "2022-09-18",
						"transactionAmount": {
							"amount": "-3.9",
							"currency": "EUR"
						},
						"creditorName": "Caffe",
						"remittanceInformationUnstructured": "-",
						"remittanceInformationUnstructuredArray": [
							"-"
						],
						"additionalInformation": "95ce7224-486f-4bd3-a1e0-d5366c4f8837",
						"bankTransactionCode": "DNOCFRESACME"
					},
					{
						"transactionId": "935d7d4a-5451-483f-89dd-ea1472115126",
						"bookingDate": "2022-09-12",
						"valueDate": "2022-09-12",
						"transactionAmount": {
							"amount": "75.0",
							"currency": "EUR"
						},
						"debtorName": "Correction Department",
						"debtorAccount": {
							"iban": "BG18RZBB91550123456789"
						},
						"remittanceInformationUnstructured": "43e4dd40-b502-4c2c-a7c0-8aa52d0d4767",
						"remittanceInformationUnstructuredArray": [
							"Some information about the transaction"
						],
						"bankTransactionCode": "PMNT-RCDT-ESCT"
					}
				],
				"pending": [
					{
						"transactionId": "45ea805a-b772-432b-a802-d65aee4b3cb5",
						"bookingDate": "2022-09-18",
						"valueDate": "2022-09-18",
						"transactionAmount": {
							"amount": "-1.99",
							"currency": "EUR"
						},
						"creditorName": "PAYPAL *GOOGLE GOOGLE",
						"additionalInformation": "bb46ddfc-8d63-458c-aae9-bd0386af0b99",
						"bankTransactionCode": "YBDSWERTPMDC"
					},
					{
						"transactionId": "52c427fa-2bc0-4b43-a43c-06a1d709eabe",
						"bookingDate": "2022-09-17",
						"valueDate": "2022-09-17",
						"transactionAmount": {
							"amount": "-15.69",
							"currency": "EUR"
						},
						"creditorName": "Test Creditor",
						"additionalInformation": "19659417-03e3-466c-b97b-b98eda4ee811",
						"bankTransactionCode": "YYYY-ZZZZ-DDDD"
					}
				]
			}
		}`
	srv := startServerWithAutoAuth(responsePayload, http.StatusOK)
	defer srv.Close()

	client := createTestNordigen(srv)

	underTest := client.Account().Transaction(uuid.New())

	from := time.Now().Add(-72 * time.Hour)
	to := time.Now()

	// When/Act
	res, err := underTest.Get(&from, &to)

	// Then/Assert
	if err != nil {
		t.Fatalf("unexpected error occurred: %s", err)
	}

	if res == nil {
		t.Fatal("transactions response expected, nil returned")
	}

	if len(res.Transactions.Booked) != 2 {
		t.Fatalf("2 booked transactions expected, %d received", len(res.Transactions.Booked))
	}

	if len(res.Transactions.Pending) != 2 {
		t.Fatalf("2 pending transactions expected, %d received", len(res.Transactions.Pending))
	}

	if len(res.Transactions.Information) != 0 {
		t.Fatalf("0 information transactions expected, %d received", len(res.Transactions.Information))
	}

	if res.Transactions.Booked[0].Amount.Amount != "-3.9" {
		t.Fatalf(
			`the amount of the first booked transaction expected to be "-3.9", %s received`,
			res.Transactions.Booked[0].Amount.Amount)
	}

	if res.Transactions.Booked[1].DebtorName != "Correction Department" {
		t.Fatalf(
			`debtor name of the second booked transaction expected to be "Correction Department", %s received`,
			res.Transactions.Booked[1].DebtorName)
	}

	if res.Transactions.Booked[1].Amount.Amount != "75.0" {
		t.Fatalf(
			`the amount of the second booked transaction expected to be "75.0", %s received`,
			res.Transactions.Booked[1].Amount.Amount)
	}

	if res.Transactions.Booked[1].RemittanceInformationUnstructuredArray[0] != "Some information about the transaction" {
		t.Fatalf(
			`the first value of the remittance info array for the second booked transaction expected to be
			"Some information about the transaction", %s received`,
			res.Transactions.Booked[1].RemittanceInformationUnstructuredArray[0])
	}
}

func testTransactionsGetApiError(c *Nordigen) error {
	underTest := c.Account().Transaction(uuid.New())
	_, err := underTest.Get(nil, nil)

	return err
}
