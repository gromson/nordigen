package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/pkg/errors"
)

type testBody struct {
	*bytes.Buffer
}

func (b *testBody) Close() error {
	return nil
}

func Test_createApiError(t *testing.T) {
	// What/Arrange
	body := &testBody{
		bytes.NewBuffer([]byte{}),
	}

	err := json.NewEncoder(body).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"detail": "something went wrong",
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	res := &http.Response{
		Status:     "Bad Request",
		StatusCode: 400,
		Body:       body,
	}

	// When/Act
	resErr := createApiError(res)

	// Then/Arrange
	apiErr := &ApiError{}
	if !errors.As(resErr, &apiErr) {
		t.Fatalf("invalid error type created: %v", resErr)
	}

	if _, ok := (*apiErr)["error"]; !ok {
		t.Fatalf("invalid error type created: %v", apiErr)
	}

}

func TestApiError_Error(t *testing.T) {
	// What/Arrange
	expectedMsgOpt1 := `0 - {"institution_id":{"detail":"one"},"max_history_days":{"detail":"two"}}`
	expectedMsgOpt2 := `0 - {"max_history_days":{"detail": "two"},"institution_id":{"detail":"one"}}`

	underTest := &ApiError{
		"institution_id": map[string]interface{}{
			"detail": "one",
		},
		"max_history_days": map[string]interface{}{
			"detail": "two",
		},
	}

	// When/Act
	msg := underTest.Error()

	// Then/Assert
	if expectedMsgOpt1 != msg && expectedMsgOpt2 != msg {
		t.Fatalf(`wanted: "%s" or "%s", got: "%s"`, expectedMsgOpt1, expectedMsgOpt2, msg)
	}
}

type errorStatusCodeTestCase struct {
	name  string
	given interface{}
	want  int
}

var errorStatusCodeTestCases = []*errorStatusCodeTestCase{
	{"float64 as a value", float64(500), 500},
	{"int as a value", 400, 400},
	{"no value", nil, 0},
	{"unsupported value float32", float32(500), 0},
	{"unsupported value string", "ta da", 0},
}

func TestApiError_StatusCode(t *testing.T) {
	t.Parallel()
	for _, tc := range errorStatusCodeTestCases {
		t.Run(tc.name, testErrorStatusCode(tc))
	}
}

func testErrorStatusCode(tc *errorStatusCodeTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		// What/Arrange
		underTest := ApiError{}
		underTest[statusCodeKey] = tc.given

		// When/Act
		code := underTest.StatusCode()

		// Then/Assert
		if code != tc.want {
			t.Fatalf("status code in the error expected to be 500, %d given", code)
		}
	}
}
