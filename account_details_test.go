package nordigen

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestBalanceDetails_Get(t *testing.T) {
	t.Parallel()
	t.Run("getting account details success", testAccountDetailsGetOk)
	t.Run("getting account details API error", testApiErrorResponse(testAccountDetailsGetApiError))
}

func testAccountDetailsGetOk(t *testing.T) {
	// What/Arrange
	responsePayload := `{
		"account": {
			"resourceId": "0e5556c0-5960-40b7-9374-514215857c2f",
			"currency": "EUR",
			"ownerName": "John Doe",
			"name": "Main",
			"product": "Space",
			"cashAccountType": "CACC",
			"status": "enabled",
			"usage": "PRIV"
		}
	}`
	srv := startServerWithAutoAuth(responsePayload, http.StatusOK)
	defer srv.Close()

	client := createTestNordigen(srv)

	underTest := client.Account().Details(uuid.New())

	// When/Act
	res, err := underTest.Get()

	// Then/Assert
	if err != nil {
		t.Fatalf("unexpected error occurred: %s", err)
	}

	if res == nil {
		t.Fatal("details response expected, nil returned")
	}

	if res.Account.ResourceID.String() != "0e5556c0-5960-40b7-9374-514215857c2f" {
		t.Fatalf(
			`expected resourceId is "0e5556c0-5960-40b7-9374-514215857c2f", "%s" is returned`,
			res.Account.ResourceID)
	}

	if res.Account.OwnerName != "John Doe" {
		t.Fatalf(`expected ownerName is "John Doe", "%s" is returned`, res.Account.OwnerName)
	}

	if res.Account.Name != "Main" {
		t.Fatalf(`expected name is "Main", "%s" is returned`, res.Account.Name)
	}

	if res.Account.Product != "Space" {
		t.Fatalf(`expected product is "Space", "%s" is returned`, res.Account.Product)
	}

	if res.Account.CashAccountType != "CACC" {
		t.Fatalf(`expected cashAccountType is "CACC", "%s" is returned`, res.Account.CashAccountType)
	}

	if res.Account.Status != "enabled" {
		t.Fatalf(`expected status is "enabled", "%s" is returned`, res.Account.Status)
	}

	if res.Account.Usage != "PRIV" {
		t.Fatalf(`expected usage is "PRIV", "%s" is returned`, res.Account.Usage)
	}
}

func testAccountDetailsGetApiError(c *Nordigen) error {
	underTest := c.Account().Details(uuid.New())
	_, err := underTest.Get()

	return err
}
