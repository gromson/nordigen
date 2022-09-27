package nordigen

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestClient_Account(t *testing.T) {
	// What/Arrange
	underTest := createTestNordigen(nil)

	// When/Act
	account := underTest.Account()

	// Then/Assert
	if account == nil {
		t.Fatal("account resource expected")
	}
}

func TestAccountResource_Get(t *testing.T) {
	t.Parallel()
	t.Run("getting account success", testAccountResourceGetOk)
	t.Run("getting account API error", testApiErrorResponse(testAccountResourceGetApiError))
}

func testAccountResourceGetOk(t *testing.T) {
	// What/Arrange
	accountResponsePayload := `{
	  "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
	  "created": "2022-09-06T14:43:43.541Z",
	  "last_accessed": "2022-09-06T14:43:43.541Z",
	  "iban": "DE75512108001245126199",
	  "institution_id": "N26_NTSBDEB1",
	  "status": "DISCOVERED",
	  "owner_name": "John Doe"
	}`
	srv := startServerWithAutoAuth(accountResponsePayload, http.StatusOK)
	defer srv.Close()

	client := createTestNordigen(srv)

	underTest := client.Account()

	// When/Act
	accountResponse, err := underTest.Get(uuid.New())

	// Then/Assert
	if err != nil {
		t.Fatalf("unexpected error occurred: %s", err)
	}

	if accountResponse == nil {
		t.Fatal("account response expected, nil returned")
	}

	if accountResponse.Iban != "DE75512108001245126199" {
		t.Fatal("IBAN in the response object does not match IBAN in the response")
	}

	if accountResponse.Status != "DISCOVERED" {
		t.Fatal("status in the response object does not match the status in the response")
	}
}

func testAccountResourceGetApiError(c *Nordigen) error {
	underTest := c.Account()
	_, err := underTest.Get(uuid.New())

	return err
}
