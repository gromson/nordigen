package nordigen

import (
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestClient_EndUserAgreement(t *testing.T) {
	// What/Arrange
	underTest := createTestNordigen(nil)

	// When/Act
	agreement := underTest.EndUserAgreement()

	// Then/Assert
	if agreement == nil {
		t.Fatal("agreement resource expected")
	}
}

func TestEndUserAgreement_Get(t *testing.T) {
	t.Parallel()
	t.Run("getting end user agreement success", testEndUserAgreementGetOk)
	t.Run("getting end user agreement API error", testApiErrorResponse(testEndUserAgreementGetApiError))
}

func testEndUserAgreementGetOk(t *testing.T) {
	// What/Arrange
	responsePayload := `{
		"id": "cb4b85a7-61c9-42d5-be3d-4cb64886bf8c",
		"created": "2022-07-11T22:26:22.250723Z",
		"max_historical_days": 180,
		"access_valid_for_days": 30,
		"access_scope": [
			"balances",
			"details",
			"transactions"
		],
		"accepted": "2022-07-11T22:31:15.950235Z",
		"institution_id": "N26_NTSBDEB1"
	}`
	srv := startServerWithAutoAuth(responsePayload, http.StatusOK)
	defer srv.Close()

	client := createTestNordigen(srv)

	underTest := client.EndUserAgreement()

	// When/Act
	response, err := underTest.Get(uuid.New())

	// Then/Assert
	if err != nil {
		t.Fatalf("unexpected error occurred: %s", err)
	}

	if response == nil {
		t.Fatal("account response expected, nil returned")
	}

	if response.ID.String() != "cb4b85a7-61c9-42d5-be3d-4cb64886bf8c" {
		t.Fatal("ID in the response object does not match ID in the response")
	}

	if response.Created.Format(time.RFC3339Nano) != "2022-07-11T22:26:22.250723Z" {
		t.Fatal("creation time in the response object doesn't match the creation time in the response")
	}

	if len(response.AccessScopes) != 3 {
		t.Fatal("number of access scopes is not equal to the expected amount")
	}
}

func testEndUserAgreementGetApiError(c *Nordigen) error {
	underTest := c.EndUserAgreement()
	_, err := underTest.Get(uuid.New())

	return err
}

func TestEndUserAgreement_List(t *testing.T) {
	t.Parallel()
	t.Run("listing end user agreement success", testEndUserAgreementListOk)
	t.Run("listing end user agreement API error", testApiErrorResponse(testEndUserAgreementListApiError))
}

func testEndUserAgreementListOk(t *testing.T) {
	// What/Arrange
	srv := startListServerWithAutoAuth(func() interface{} {
		accepted := time.Now().Add(-2 * 24 * time.Hour)
		return EndUserAgreementResponse{
			ID:                 uuid.MustParse("cb4b85a7-61c9-42d5-be3d-4cb64886bf8c"),
			Created:            time.Now(),
			MaxHistoricalDays:  180,
			AccessValidForDays: 30,
			AccessScopes:       []string{"balances", "details", "transactions"},
			Accepted:           &accepted,
			InstitutionID:      "TEST_INSTITUTION",
		}
	})
	defer srv.Close()

	client := createTestNordigen(srv)

	underTest := client.EndUserAgreement()

	// When/Act
	response, err := underTest.List()

	// Then/Assert
	if err != nil {
		t.Fatalf("unexpected error occurred: %s", err)
	}

	assertCollectionResponse(t, response)
}

func testEndUserAgreementListApiError(c *Nordigen) error {
	underTest := c.EndUserAgreement()
	_, err := underTest.List()

	return err
}

func TestEndUserAgreement_Create(t *testing.T) {
	t.Parallel()
	t.Run("creating end user agreement success", testEndUserAgreementCreateOk)
	t.Run("creating end user agreement API error", testApiErrorResponse(testEndUserAgreementCreateApiError))
}

func testEndUserAgreementCreateOk(t *testing.T) {
	// What/Arrange
	responsePayload := `{
		"id": "cb4b85a7-61c9-42d5-be3d-4cb64886bf8c",
		"created": "2022-07-11T22:26:22.250723Z",
		"max_historical_days": 180,
		"access_valid_for_days": 30,
		"access_scope": [
			"balances",
			"details",
			"transactions"
		],
		"accepted": "2022-07-11T22:31:15.950235Z",
		"institution_id": "N26_NTSBDEB1"
	}`
	srv := startServerWithAutoAuth(responsePayload, http.StatusOK)
	defer srv.Close()

	client := createTestNordigen(srv)

	underTest := client.EndUserAgreement()

	// When/Act
	response, err := underTest.Create(&CreateAgreementRequest{
		InstitutionID:      "N26_NTSBDEB1",
		MaxHistoricalDays:  180,
		AccessValidForDays: 30,
		AccessScope:        []string{"balances", "details", "transactions"},
	})

	// Then/Assert
	if err != nil {
		t.Fatalf("unexpected error occurred: %s", err)
	}

	if response == nil {
		t.Fatal("account response expected, nil returned")
	}

	if response.ID.String() != "cb4b85a7-61c9-42d5-be3d-4cb64886bf8c" {
		t.Fatal("ID in the response object does not match ID in the response")
	}

	if response.Created.Format(time.RFC3339Nano) != "2022-07-11T22:26:22.250723Z" {
		t.Fatal("creation time in the response object doesn't match the creation time in the response")
	}

	if len(response.AccessScopes) != 3 {
		t.Fatal("number of access scopes is not equal to the expected amount")
	}
}

func testEndUserAgreementCreateApiError(c *Nordigen) error {
	underTest := c.EndUserAgreement()
	_, err := underTest.Create(&CreateAgreementRequest{
		InstitutionID:      "N26_NTSBDEB1",
		MaxHistoricalDays:  180,
		AccessValidForDays: 30,
		AccessScope:        []string{"balances", "details", "transactions"},
	})

	return err
}

func TestEndUserAgreement_Accept(t *testing.T) {
	t.Parallel()
	t.Run("accepting end user agreement success", testEndUserAgreementAcceptOk)
	t.Run("accepting end user agreement API error", testApiErrorResponse(testEndUserAgreementAcceptApiError))
}

func testEndUserAgreementAcceptOk(t *testing.T) {
	// What/Arrange
	responsePayload := `{
		"id": "cb4b85a7-61c9-42d5-be3d-4cb64886bf8c",
		"created": "2022-07-11T22:26:22.250723Z",
		"max_historical_days": 180,
		"access_valid_for_days": 30,
		"access_scope": [
			"balances",
			"details",
			"transactions"
		],
		"accepted": "2022-07-11T22:31:15.950235Z",
		"institution_id": "N26_NTSBDEB1"
	}`
	srv := startServerWithAutoAuth(responsePayload, http.StatusOK)
	defer srv.Close()

	client := createTestNordigen(srv)

	underTest := client.EndUserAgreement()

	// When/Act
	response, err := underTest.Accept(
		uuid.MustParse("cb4b85a7-61c9-42d5-be3d-4cb64886bf8c"),
		&AcceptEndUserAgreementRequest{
			UserAgent: "TestAgent",
			IPAddress: net.IPv4(42, 42, 42, 42),
		})

	// Then/Assert
	if err != nil {
		t.Fatalf("unexpected error occurred: %s", err)
	}

	if response == nil {
		t.Fatal("account response expected, nil returned")
	}

	if response.ID.String() != "cb4b85a7-61c9-42d5-be3d-4cb64886bf8c" {
		t.Fatal("ID in the response object does not match ID in the response")
	}

	if response.Created.Format(time.RFC3339Nano) != "2022-07-11T22:26:22.250723Z" {
		t.Fatal("creation time in the response object doesn't match the creation time in the response")
	}

	if len(response.AccessScopes) != 3 {
		t.Fatal("number of access scopes is not equal to the expected amount")
	}
}

func testEndUserAgreementAcceptApiError(c *Nordigen) error {
	underTest := c.EndUserAgreement()
	_, err := underTest.Accept(
		uuid.MustParse("cb4b85a7-61c9-42d5-be3d-4cb64886bf8c"),
		&AcceptEndUserAgreementRequest{
			UserAgent: "TestAgent",
			IPAddress: net.IPv4(42, 42, 42, 42),
		})

	return err
}

func TestEndUserAgreement_Delete(t *testing.T) {
	t.Parallel()
	t.Run("deleting end user agreement success", testEndUserAgreementDeleteOk)
	t.Run("deleting end user agreement API error", testApiErrorResponse(testEndUserAgreementDeleteApiError))
}

func testEndUserAgreementDeleteOk(t *testing.T) {
	// What/Arrange
	srv := startServerWithAutoAuth("", http.StatusOK)
	defer srv.Close()

	client := createTestNordigen(srv)

	underTest := client.EndUserAgreement()

	// When/Act
	err := underTest.Delete(uuid.New())

	// Then/Assert
	if err != nil {
		t.Fatalf("unexpected error occurred: %s", err)
	}
}

func testEndUserAgreementDeleteApiError(c *Nordigen) error {
	underTest := c.EndUserAgreement()
	return underTest.Delete(uuid.New())
}
