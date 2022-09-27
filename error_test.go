package nordigen

import (
	"net/http"
	"testing"
)

func TestApiClientErrorResponse_Error(t *testing.T) {
	t.Parallel()
	t.Run("sub structure", testApiClientErrorResponseSubStructure)
	t.Run("sub structure", testApiClientErrorResponseFlatStructure)
}

func testApiClientErrorResponseSubStructure(t *testing.T) {
	// What/Arrange
	expectedMsgOpt1 := "400 - one: two; three: four"
	expectedMsgOpt2 := "400 - three: four; one: two"

	underTest := &ApiClientErrorResponse{
		"institution_id": map[string]interface{}{
			"summary": "one",
			"detail":  "two",
		},
		"max_history_days": map[string]interface{}{
			"summary": "three",
			"detail":  "four",
		},
		"status_code": http.StatusBadRequest,
	}

	// When/Act
	msg := underTest.Error()

	// Then/Assert
	if expectedMsgOpt1 != msg && expectedMsgOpt2 != msg {
		t.Fatalf(`wanted: "%s" or "%s", got: "%s"`, expectedMsgOpt1, expectedMsgOpt2, msg)
	}
}

func testApiClientErrorResponseFlatStructure(t *testing.T) {
	// What/Arrange
	expectedMsgOpt1 := "400 - one; two"
	expectedMsgOpt2 := "400 - two; one"

	underTest := &ApiClientErrorResponse{
		"summary":     "one",
		"detail":      "two",
		"status_code": http.StatusBadRequest,
	}

	// When/Act
	msg := underTest.Error()

	// Then/Assert
	if expectedMsgOpt1 != msg && expectedMsgOpt2 != msg {
		t.Fatalf(`wanted: "%s" or "%s", got: "%s"`, expectedMsgOpt1, expectedMsgOpt2, msg)
	}
}
