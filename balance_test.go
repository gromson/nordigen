package nordigen

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestBalanceResource_Get(t *testing.T) {
	t.Parallel()
	t.Run("getting balances success", testBalanceGetOk)
	t.Run("getting balances API error", testApiErrorResponse(testBalanceGetApiError))
}

func testBalanceGetOk(t *testing.T) {
	// What/Arrange
	responsePayload := `{
	  "balances": [
		{
		  "balanceAmount": {
			"amount": "657.49",
			"currency": "EUR"
		  },
		  "balanceType": "expected",
		  "lastChangeDateTime": "2022-09-06T19:11:44.753Z",
		  "referenceDate": "2022-09-07"
		},
		{
		  "balanceAmount": {
			"amount": "185.67",
			"currency": "EUR"
		  },
		  "balanceType": "expected",
		  "referenceDate": "2022-09-06"
		}
	  ]
	}`
	srv := startServerWithAutoAuth(responsePayload, http.StatusOK)
	defer srv.Close()

	client := createTestNordigen(srv)

	underTest := client.Account().Balance(uuid.New())

	// When/Act
	res, err := underTest.Get()

	// Then/Assert
	if err != nil {
		t.Fatalf("unexpected error occurred: %s", err)
	}

	if res == nil {
		t.Fatal("balances response expected, nil returned")
	}

	if len(res.Balances) != 2 {
		t.Fatal("2 balances expected")
	}

	expectedLastChangeDateTime, _ := time.Parse(time.RFC3339, "2022-09-06T19:11:44.753Z")
	if res.Balances[0].LastChangeDateTime != expectedLastChangeDateTime {
		t.Fatalf("expected last change date is %v, %v given",
			expectedLastChangeDateTime,
			res.Balances[0].ReferenceDate)
	}

	if res.Balances[0].ReferenceDate != "2022-09-07" {
		t.Fatalf("expected reference date is 2022-09-07, %v given", res.Balances[0].ReferenceDate)
	}

	if res.Balances[0].BalanceAmount.Amount != "657.49" {
		t.Fatalf("expected balance amount is 657.49, %s given", res.Balances[0].BalanceAmount.Amount)
	}
}

func testBalanceGetApiError(c *Nordigen) error {
	underTest := c.Account().Balance(uuid.New())
	_, err := underTest.Get()

	return err
}
