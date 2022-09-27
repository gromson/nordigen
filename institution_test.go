package nordigen

import (
	"net/http"
	"testing"
)

func TestInstitutionResource_List(t *testing.T) {
	t.Parallel()
	t.Run("listing institutions success", testInstitutionListOk)
	t.Run("listing institutions API error", testApiErrorResponse(testInstitutionListApiError))
}

func testInstitutionListOk(t *testing.T) {
	// What/Arrange
	responsePayload := `[
		{
			"id": "DIREKT_HELADEF1822",
			"name": "1822direkt",
			"bic": "HELADEF1822",
			"transaction_total_days": "730",
			"countries": ["DE"],
			"logo": "https://cdn.nordigen.com/ais/DIREKT_HELADEF1822.png"
		},
		{
			"id": "AACHENER_BANK_GENODED1AAC",
			"name": "Aachener Bank",
			"bic": "GENODED1AAC",
			"transaction_total_days": "400",
			"countries":["DE"],
			"logo": "https://cdn.nordigen.com/ais/VOLKSBANK_NIEDERGRAFSCHAFT_GENODEF1HOO.png"
		}
	]`

	srv := startServerWithAutoAuth(responsePayload, http.StatusOK)
	defer srv.Close()

	client := createTestNordigen(srv)

	underTest := client.Institution()

	// When/Act
	institutions, err := underTest.List("DE")

	// Then/Assert
	if err != nil {
		t.Fatalf("unexpected error while successful API response: %s", err)
	}

	if len(institutions) != 2 {
		t.Fatalf("number of institutions is expected to be equal to 2, %d given", len(institutions))
	}

	bank := institutions[0]

	if bank.ID != "DIREKT_HELADEF1822" {
		t.Fatalf("expected bank ID is \"DIREKT_HELADEF1822\", %s returned", bank.ID)
	}

	if bank.Name != "1822direkt" {
		t.Fatalf("expected bank name is \"1822direkt\", %s returned", bank.Name)
	}

	if bank.BIC != "HELADEF1822" {
		t.Fatalf("expected bank BIC is \"HELADEF1822\", %s returned", bank.BIC)
	}

	if bank.TransactionTotalDays != 730 {
		t.Fatalf("expected bank transaction_total_days is \"730\", %d returned", bank.TransactionTotalDays)
	}

	if len(bank.Countries) != 1 {
		t.Fatalf("expected bank countries lenth is 1, %d returned", len(bank.Countries))
	}

	if bank.Countries[0] != "DE" {
		t.Fatalf("expected bank country is \"DE\", %s returned", bank.Countries[0])
	}

	if bank.LogoUrl != "https://cdn.nordigen.com/ais/DIREKT_HELADEF1822.png" {
		t.Fatalf("expected bank logo URL is \"https://cdn.nordigen.com/ais/DIREKT_HELADEF1822.png\", %s returned", bank.LogoUrl)
	}
}

func testInstitutionListApiError(client *Nordigen) error {
	underTest := client.Institution()
	_, err := underTest.ListWithEnabledPayments("DE")
	return err
}

func TestInstitutionResource_Get(t *testing.T) {
	t.Parallel()
	t.Run("getting institutions success", testInstitutionGetOk)
	t.Run("getting institutions API error", testApiErrorResponse(testInstitutionGetApiError))
}

func testInstitutionGetOk(t *testing.T) {
	// What/Arrange
	responsePayload := `{
		"id": "DIREKT_HELADEF1822",
		"name": "1822direkt",
		"bic" :"HELADEF1822",
		"transaction_total_days": "730",
		"countries": ["DE"],
		"logo":"https://cdn.nordigen.com/ais/DIREKT_HELADEF1822.png"
	}`
	srv := startServerWithAutoAuth(responsePayload, http.StatusOK)
	defer srv.Close()

	client := createTestNordigen(srv)

	underTest := client.Institution()

	// When/Act
	institution, err := underTest.Get("DIREKT_HELADEF1822")

	// Then/Assert
	if err != nil {
		t.Fatalf("unexpected error while successful API response: %s", err)
	}

	if institution.ID != "DIREKT_HELADEF1822" {
		t.Fatalf("expected bank ID is \"DIREKT_HELADEF1822\", %s returned", institution.ID)
	}

	if institution.Name != "1822direkt" {
		t.Fatalf("expected bank name is \"1822direkt\", %s returned", institution.Name)
	}

	if institution.BIC != "HELADEF1822" {
		t.Fatalf("expected bank BIC is \"HELADEF1822\", %s returned", institution.BIC)
	}

	if institution.TransactionTotalDays != 730 {
		t.Fatalf(
			"expected bank transaction_total_days is \"730\", %d returned",
			institution.TransactionTotalDays)
	}

	if len(institution.Countries) != 1 {
		t.Fatalf("expected bank countries lenth is 1, %d returned", len(institution.Countries))
	}

	if institution.Countries[0] != "DE" {
		t.Fatalf("expected bank country is \"DE\", %s returned", institution.Countries[0])
	}

	if institution.LogoUrl != "https://cdn.nordigen.com/ais/DIREKT_HELADEF1822.png" {
		t.Fatalf(
			"expected bank logo URL is \"https://cdn.nordigen.com/ais/DIREKT_HELADEF1822.png\", %s returned",
			institution.LogoUrl)
	}
}

func testInstitutionGetApiError(client *Nordigen) error {
	underTest := client.Institution()
	_, err := underTest.Get("DIREKT_HELADEF1822")
	return err
}
