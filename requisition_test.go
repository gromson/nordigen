package nordigen

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestClient_Requisition(t *testing.T) {
	// What/Arrange
	underTest := createTestNordigen(nil)

	// When/Act
	requisition := underTest.Requisition()

	// Then/Assert
	if requisition == nil {
		t.Fatal("requisition resource expected")
	}
}

func TestRequisitionResource_Get(t *testing.T) {
	t.Parallel()
	t.Run("getting requisition success", testRequisitionGetOk)
	t.Run("getting requisition API error", testApiErrorResponse(testRequisitionGetApiError))
}

func testRequisitionGetOk(t *testing.T) {
	// What/Arrange
	responsePayload := `{
		"id": "1071addf-e971-4fcf-9cb4-83e10e050eff",
		"created": "2022-07-11T22:29:09.470230Z",
		"redirect": "https://redirect.com",
		"status": "EX",
		"institution_id": "TEST_INSTITUTION",
		"agreement": "3d0e2267-40a4-47ed-a7d6-8b4a830c1cfb",
		"reference": "124151",
		"accounts": [
			"9febb941-0886-4d03-991f-98111c68bf26",
			"713e989a-0ba3-4015-a64b-863942711f14"
		],
		"user_language": "EN",
		"link": "https://ob.nordigen.com/psd2/start/e530bfd2-e05f-465e-9d16-d4859002fdbc/TEST_INSTITUTION",
		"ssn": null,
		"account_selection": true,
		"redirect_immediate": true
	}`
	srv := startServerWithAutoAuth(responsePayload, http.StatusOK)
	defer srv.Close()

	client := createTestNordigen(srv)

	underTest := client.Requisition()

	// When/Act
	response, err := underTest.Get(uuid.New())

	// Then/Assert
	if err != nil {
		t.Fatalf("unexpected error occurred: %s", err)
	}

	if response == nil {
		t.Fatal("account response expected, nil returned")
	}

	if response.ID.String() != "1071addf-e971-4fcf-9cb4-83e10e050eff" {
		t.Fatal("ID in the response object does not match ID in the response")
	}

	if response.Created.Format(time.RFC3339Nano) != "2022-07-11T22:29:09.47023Z" {
		t.Fatal("creation time in the response object doesn't match the creation time in the response")
	}

	if response.RedirectUrl != "https://redirect.com" {
		t.Fatal("redirect in the response object doesn't match the redirect in the response")
	}

	if response.Status != "EX" {
		t.Fatal("status in the response object doesn't match the status in the response")
	}

	if response.InstitutionID != "TEST_INSTITUTION" {
		t.Fatal("institution_id in the response object doesn't match the institution_id in the response")
	}

	if response.AgreementID.String() != "3d0e2267-40a4-47ed-a7d6-8b4a830c1cfb" {
		t.Fatal("agreement in the response object doesn't match the agreement in the response")
	}

	if response.Reference != "124151" {
		t.Fatal("reference in the response object doesn't match the reference in the response")
	}

	if len(response.Accounts) != 2 {
		t.Fatal("number of accounts in the response object doesn't match the number of accounts in the response")
	}

	if response.UserLanguage != "EN" {
		t.Fatal("user language in the response object doesn't match the user language in the response")
	}

	if response.Link != "https://ob.nordigen.com/psd2/start/e530bfd2-e05f-465e-9d16-d4859002fdbc/TEST_INSTITUTION" {
		t.Fatal("link in the response object doesn't match the link in the response")
	}

	if response.Ssn != "" {
		t.Fatal("ssn in the response object doesn't match the ssn in the response")
	}

	if response.AccountSelection != true {
		t.Fatal("account_selection in the response object doesn't match the account_selection in the response")
	}

	if response.RedirectImmediate != true {
		t.Fatal("redirect_immediate in the response object doesn't match the redirect_immediate in the response")
	}
}

func testRequisitionGetApiError(c *Nordigen) error {
	underTest := c.Requisition()
	_, err := underTest.Get(uuid.New())

	return err
}
