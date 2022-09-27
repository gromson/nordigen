package nordigen

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
)

const (
	testApiEntryListAmount        = 6
	testApiEntryListResponseLimit = 2
)

type invalidSecretsTestCase struct {
	desc      string
	secretId  string
	secretKey string
}

var invalidSecretTestCases = []*invalidSecretsTestCase{
	{"invalid secret ID symbol", "b6789fd6-95ee-4093-a03b-b003b7d7858z", "eafc3b"},
	{"secret ID is empty string", "", "eafc3b"},
	{"invalid secret UUID format", "95ee-4093-a03b", "eafc3b"},
	{"invalid secret key symbol", "b6789fd6-95ee-4093-a03b-b003b7d7858c", "eafg3b"},
	{"secret key is empty string", "b6789fd6-95ee-4093-a03b-b003b7d7858c", ""},
	{"invalid secret key format", "b6789fd6-95ee-4093-a03b-b003b7d7858c", "eaf-c3b"},
	{"invalid secret key format 2", "b6789fd6-95ee-4093-a03b-b003b7d7858c", "ea-fc3b"},
}

func TestNewClient(t *testing.T) {
	t.Parallel()
	t.Run("client creation success", testNewClientSuccess)
	for _, tc := range invalidSecretTestCases {
		t.Run("client creation invalid secret", testNewClientError(tc))
	}
}

func testNewClientSuccess(t *testing.T) {
	// What/Arrange
	secretId := "b6789fd6-95ee-4093-a03b-b003b7d7858a"
	secretKey := "eafc3b" // 234, 252, 59
	secretKeyRaw := []byte{234, 252, 59}

	// When/Act
	c, err := New(secretId, secretKey)

	// Then/Assert
	if err != nil {
		t.Fatalf("error creating a client: %s", err)
	}

	if c.SecretID.String() != secretId {
		t.Fatalf(`invalid secretId has been set: "%s" expected, "%s" set`, secretId, c.SecretID)
	}

	for i, b := range c.SecretKey {
		if b != secretKeyRaw[i] {
			t.Fatalf("invalid secretKey has been set: %v expected, %v set", secretKeyRaw, c.SecretKey)
		}
	}
}

func testNewClientError(tc *invalidSecretsTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		// What/Arrange
		secretId := tc.secretId
		secretKey := tc.secretKey

		// When/Act
		_, err := New(secretId, secretKey)

		// Then/Assert
		if err == nil {
			t.Fatalf(`error expected when invalid parameter provided: ("%s", "%s")`, secretId, secretKey)
		}
	}
}

func testApiErrorResponse(act func(c *Nordigen) error) func(t *testing.T) {
	return func(t *testing.T) {
		// What/Arrange
		srv := startServerWithAutoAuth("", http.StatusServiceUnavailable)
		defer srv.Close()

		underTest := createTestNordigen(srv)

		// When/Act
		err := act(underTest)

		// Then/Assert
		if err == nil {
			t.Fatal("error expected if API returns an error response")
		}
	}
}

func createTestNordigen(srv *httptest.Server) *Nordigen {
	c := MustNew(uuid.New(), []byte{12, 23, 42})
	c.accessTokenExpiration = time.Now().Add(time.Second)

	if srv != nil {
		srvUrl, _ := url.Parse(srv.URL)
		c.restClient.BaseUrl = srvUrl
	}

	return c
}

func startServer(responsePayload string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		if _, err := w.Write([]byte(responsePayload)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
}

func startServerWithAutoAuth(responsePayload string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token/new/" {
			authenticate(w)
			return
		}

		w.WriteHeader(statusCode)
		if _, err := w.Write([]byte(responsePayload)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
}

func startListServerWithAutoAuth(gen func() interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token/new/" {
			authenticate(w)
			return
		}

		params := r.URL.Query()

		offset := 0
		limit := testApiEntryListResponseLimit
		var err error

		if offsetString := params.Get("offset"); offsetString != "" {
			if offset, err = strconv.Atoi(offsetString); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		if limitString := params.Get("limit"); limitString != "" {
			if limit, err = strconv.Atoi(limitString); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		results := make([]interface{}, 0, 2)

		if offset+limit <= testApiEntryListAmount {
			results = append(results, gen(), gen())
		}

		payload := struct {
			Count    int           `json:"count"`
			Next     *string       `json:"next"`
			Previous *string       `json:"previous"`
			Results  []interface{} `json:"results"`
		}{
			Count:    testApiEntryListAmount,
			Next:     nil,
			Previous: nil,
			Results:  results,
		}

		if err := json.NewEncoder(w).Encode(payload); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
}

func assertCollectionResponse[Response any](t *testing.T, response *CollectionResponse[Response]) {
	if response.Count() != testApiEntryListAmount {
		t.Fatalf("%d list item expected to be returned from the server, %d returned",
			testApiEntryListAmount,
			response.Count())
	}

	for i := 0; i < response.Count(); i++ {
		item, err := response.Next()
		if err != nil {
			t.Fatalf("unexpected error occured while getting next item from the list: %s", err)
		}

		if item == nil {
			t.Fatalf("item is expected to be a Response object but it's nil")
		}
	}

	item, err := response.Next()

	if err != nil {
		t.Fatalf("unexpected error occured while getting next item from the list: %s", err)
	}

	if item != nil {
		t.Fatalf("item is expected to be nil, %v", item)
	}
}

func authenticate(w http.ResponseWriter) {
	responsePayload := `{
		"access": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZXhwIjoxNjY0MjAxMDk3LCJqdGkiOiIxMjM0NTY3ODkwcXdlcnR5IiwiaWQiOjEzMTM0LCJzZWNyZXRfaWQiOiJkMmE1YzU2My0yYmI2LTQ2M2UtOTMyOC02ODhiZTc0MWIzMWEiLCJhbGxvd2VkX2NpZHJzIjpbIjAuMC4wLjAvMCIsIjo6LzAiXX0.m0AE4HrKFa38k6bBKV7LMaovEOEzyOevvvjDzLXCDBk",
		"access_expires": 86400,
		"refresh": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImV4cCI6MTY2NjcwNjY5NywianRpIjoiMDk4NzY1NDMyMXl0cmV3cSIsImlkIjoxMzEzNCwic2VjcmV0X2lkIjoiMjM3M2E2ODItY2FhYi00ODk4LWIzODEtMWY1NjhiMGExOTk5IiwiYWxsb3dlZF9jaWRycyI6WyIwLjAuMC4wLzAiLCI6Oi8wIl19.fuS3k1UW2PVFzt0K4pCVrd5pwIfUsum6lTlswkiwEzk",
		"refresh_expires": 2592000
	}`

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(responsePayload)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
