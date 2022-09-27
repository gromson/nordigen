package typ

import (
	"encoding/json"
	"testing"
)

func TestHexBytes_MarshalText(t *testing.T) {
	underTest := HexBytes{255, 42, 36, 12}
	expectedOutput := "ff2a240c"

	output, err := underTest.MarshalText()
	if err != nil {
		t.Fatalf("error while marshaling HexBytes: %s", err)
	}

	if expectedOutput != string(output) {
		t.Fatal("expected marshalled value of HexBytes doesn't match the actual output")
	}
}

func TestHexBytes_UnmarshalText(t *testing.T) {
	input := "4c7b0221"
	expectedResult := []byte{76, 123, 2, 33}

	underTest := &HexBytes{}
	err := underTest.UnmarshalText([]byte(input))
	if err != nil {
		t.Fatalf("error while unmarshaling HexBytes from text: %s", err)
	}

	for i, b := range *underTest {
		if b != expectedResult[i] {
			t.Fatal("")
		}
	}
}

func TestHexBytes_MarshalJSON(t *testing.T) {
	testStruct := struct {
		Key HexBytes `json:"key"`
	}{
		Key: HexBytes{23, 42, 56, 1, 0, 123},
	}

	expectedResult := `{"key":"172a3801007b"}`

	result, err := json.Marshal(testStruct)
	if err != nil {
		t.Fatal(err)
	}

	if expectedResult != string(result) {
		t.Fatalf("expected and actual results don't match: extected: %s, actual: %s", expectedResult, result)
	}
}

func TestHexBytes_UnmarshalJSON(t *testing.T) {
	jsonData := []byte(`{"key":"172a3801007b"}`)
	testStruct := struct {
		Key HexBytes `json:"key"`
	}{}

	expectedData := []byte{23, 42, 56, 1, 0, 123}

	if err := json.Unmarshal(jsonData, &testStruct); err != nil {
		t.Fatal(err)
	}

	if len(expectedData) != len(testStruct.Key) {
		t.Fatalf(
			"the length of the expected unmarshalled data doesn't match the length of an actual data: %d != %d",
			len(expectedData),
			len(testStruct.Key))
	}

	for i, b := range testStruct.Key {
		if expectedData[i] != b {
			t.Fatalf(
				"expected data does not match actual result; expected: %v; actual: %v",
				expectedData,
				testStruct.Key)
		}
	}
}
