package nordigen

import (
	"testing"

	"gromson/nordigen/rest"
)

func TestGenericCollectionResponse_Next(t *testing.T) {
	t.Parallel()
	t.Run("collection Next Ok", testCollectionResponseNextOk)
	t.Run("collection Next no client", testCollectionResponseNextNoClient)
}

func testCollectionResponseNextOk(t *testing.T) {
	// What/Arrange
	srv := startListServerWithAutoAuth(generateResourceEntryForListResponse)
	defer srv.Close()

	nordigen := createTestNordigen(srv)

	resource := &rest.GenericResource[any]{
		ID:     "/test",
		Client: nordigen.restClient,
	}

	// When/Act
	underTest := &CollectionResponse[any]{
		nordigen: nordigen,
		count:    0,
		next:     0,
		limit:    testApiEntryListResponseLimit,
		offset:   0,
		results:  make([]any, 0, 2),
		resource: resource,
		params:   nil,
	}
	err := underTest.get()

	// Then/Assert
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if underTest.Count() != testApiEntryListAmount {
		t.Fatalf("collection Count should return %d, %d returned", testApiEntryListAmount, underTest.Count())
	}

	for i := 0; i < testApiEntryListAmount; i++ {
		a, err := underTest.Next()
		if err != nil {
			t.Fatalf("unexpected error while getting the next resource: %s", err)
		}

		if a == nil {
			t.Fatal("expected requisition wasn't received")
		}
	}

	a, err := underTest.Next()
	if err != nil {
		t.Fatalf("unexpected error while getting the next non-existing requisition: %s", err)
	}

	if a != nil {
		t.Fatal("unexpected requisition found")
	}
}

func testCollectionResponseNextNoClient(t *testing.T) {
	// What/Arrange
	underTest := CollectionResponse[interface{}]{}

	// When/Act
	a, err := underTest.Next()

	// Then/Assert
	if err != nil {
		t.Fatalf("unexpected error occurred: %s", err)
	}

	if a != nil {
		t.Fatalf("returned value expected to be nil, %v given", a)
	}
}

func generateResourceEntryForListResponse() interface{} {
	return map[string]interface{}{
		"id":    1,
		"title": "one",
	}
}
