package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
)

type testResource struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type resourceSuccessGetTestCase struct {
	name            string
	id              string
	params          url.Values
	responsePayload *testResource
}

var resourceSuccessGetTestCases = []resourceSuccessGetTestCase{
	{"without ID, without values", "", nil, &testResource{ID: 0, Title: "No ID"}},
	{"with ID, without values", "123", nil, &testResource{ID: 123, Title: "One Two Three"}},
	{"with ID, with values", "123", url.Values(map[string][]string{"foo": {"bar"}}), &testResource{ID: 1, Title: "Foo: Bar"}},
}

func TestGenericResource_Get(t *testing.T) {
	t.Parallel()
	for _, tc := range resourceSuccessGetTestCases {
		t.Run("getting resource success", testGetResourceOk(tc))
	}

	t.Run("getting resource API error", testApiErrorResponse(testGetResourceApiError))
}

func testGetResourceOk(tc resourceSuccessGetTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		// What/Arrange
		responseBody := bytes.NewBuffer([]byte{})
		if err := json.NewEncoder(responseBody).Encode(tc.responsePayload); err != nil {
			t.Fatal(err)
		}

		srv := startServer(responseBody.Bytes(), http.StatusOK)
		defer srv.Close()

		client := createTestClient(srv)

		underTest := &GenericResource[testResource]{
			ID:     "/test",
			Client: client,
		}

		// When/Act
		res, err := underTest.Get(tc.id, tc.params)

		// Then/Assert
		if err != nil {
			t.Fatalf("unexpected error occurred: %s", err)
		}

		if res == nil {
			t.Fatal("resource must not be nil")
		}

		if res.ID != tc.responsePayload.ID {
			t.Fatalf(`expected resource ID "%d", "%d" returned`, tc.responsePayload.ID, res.ID)
		}

		if res.Title != tc.responsePayload.Title {
			t.Fatalf(`expected resource title "%s", "%s" returned`, tc.responsePayload.Title, res.Title)
		}
	}
}

func testGetResourceApiError(client *Client) error {
	underTest := &GenericResource[testResource]{
		ID:     "/test",
		Client: client,
	}
	_, err := underTest.Get("testID", nil)

	return err
}

func TestGenericResource_Post(t *testing.T) {
	t.Parallel()
	t.Run("Post resource success", testPostResourceOk)
	t.Run("Post resource failure", testApiErrorResponse(testPostRequisitionApiError))
}

func testPostResourceOk(t *testing.T) {
	// What/Arrange
	inputEntry := testResource{ID: 123, Title: "New Entry"}
	responsePayload := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(responsePayload).Encode(inputEntry); err != nil {
		t.Fatal(err)
	}

	srv := startServer(responsePayload.Bytes(), http.StatusOK)
	defer srv.Close()

	client := createTestClient(srv)

	underTest := GenericResource[testResource]{
		ID:     "/test",
		Client: client,
	}

	// When/Act
	createdEntry, err := underTest.Post("", inputEntry)

	// Then/Assert
	if err != nil {
		t.Fatalf("unexpected error while creating a createdEntry: %s", err)
	}

	if inputEntry.ID != createdEntry.ID {
		t.Fatal("created entry ID doesn't match input entry ID")
	}
}

func testPostRequisitionApiError(client *Client) error {
	underTest := &GenericResource[testResource]{
		ID:     "/test",
		Client: client,
	}
	_, err := underTest.Post("", testResource{ID: 123, Title: "New Entry"})

	return err
}

func TestGenericResource_Put(t *testing.T) {
	t.Parallel()
	t.Run("Put resource success", testPutResourceOk)
	t.Run("Put resource failure", testApiErrorResponse(testPutRequisitionApiError))
}

func testPutResourceOk(t *testing.T) {
	// What/Arrange
	inputEntry := testResource{Title: "Existing Entry"}
	responsePayload := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(responsePayload).Encode(inputEntry); err != nil {
		t.Fatal(err)
	}

	srv := startServer(responsePayload.Bytes(), http.StatusOK)
	defer srv.Close()

	client := createTestClient(srv)

	underTest := GenericResource[testResource]{
		ID:     "/test",
		Client: client,
	}

	// When/Act
	createdEntry, err := underTest.Put("", inputEntry)

	// Then/Assert
	if err != nil {
		t.Fatalf("unexpected error while creating a createdEntry: %s", err)
	}

	if inputEntry.ID != createdEntry.ID {
		t.Fatal("created entry ID doesn't match input entry ID")
	}
}

func testPutRequisitionApiError(client *Client) error {
	underTest := &GenericResource[testResource]{
		ID:     "/test",
		Client: client,
	}
	_, err := underTest.Put("", testResource{Title: "Existing Entry"})

	return err
}

func TestGenericResource_List(t *testing.T) {
	t.Parallel()
	t.Run("List resources success", testListResourcesOk)
	t.Run("List resources server error", testApiErrorResponse(testGetRequisitionsListApiError))
}

func testListResourcesOk(t *testing.T) {
	// What/Arrange
	srv := startListServer(generateResourceEntryForListResponse)
	defer srv.Close()

	client := createTestClient(srv)

	underTest := GenericResource[testResource]{
		ID:     "/test",
		Client: client,
	}

	// When/Act
	collection, err := underTest.List(nil)

	// Then/Assert
	if err != nil {
		t.Fatalf("unexpected error occurred: %s", err)
	}

	if len(collection) == 1 {
		t.Fatal("collection is not expected to be empty")
	}
}

func testGetRequisitionsListApiError(client *Client) error {
	underTest := GenericResource[testResource]{
		ID:     "/test",
		Client: client,
	}
	_, err := underTest.List(nil)
	return err
}

func TestGenericResource_Delete(t *testing.T) {
	t.Parallel()
	t.Run("deleting resource success", testDeleteResourceOk)
	t.Run("deleting resource API error", testApiErrorResponse(testDeleteResourceApiError))
}

func testDeleteResourceOk(t *testing.T) {
	// What/Arrange
	srv := startServer(nil, http.StatusOK)
	defer srv.Close()

	client := createTestClient(srv)

	underTest := GenericResource[testResource]{
		ID:     "/test",
		Client: client,
	}

	// When/Act
	err := underTest.Delete("testID")

	// Then/Assert
	if err != nil {
		t.Fatalf("unexpected error occurred: %s", err)
	}
}

func testDeleteResourceApiError(client *Client) error {
	underTest := GenericResource[testResource]{
		ID:     "/test",
		Client: client,
	}

	return underTest.Delete("testID")
}

func generateResourceEntryForListResponse() interface{} {
	return map[string]interface{}{
		"id":    1,
		"title": "one",
	}
}
