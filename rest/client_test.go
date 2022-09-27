package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

const (
	testApiEntryListAmount        = 6
	testApiEntryListResponseLimit = 2
)

func createTestClient(srv *httptest.Server) *Client {
	c := NewClient(&url.URL{}, nil)

	srvUrl, _ := url.Parse(srv.URL)
	c.BaseUrl = srvUrl

	return c
}

func startServer(responsePayload []byte, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		if _, err := w.Write(responsePayload); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
}

func startListServer(gen func() interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		offset := 0
		limit := testApiEntryListResponseLimit
		var err error

		if offsetString := params.Get("offset"); offsetString != "" {
			if offset, err = strconv.Atoi(offsetString); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		if limitString := params.Get("limit"); limitString != "" {
			if limit, err = strconv.Atoi(limitString); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		results := make([]interface{}, 0, limit)

		for i := 0; i < limit; i++ {
			if offset+limit <= testApiEntryListAmount {
				results = append(results, gen())
			}
		}

		if err := json.NewEncoder(w).Encode(results); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
}

func testApiErrorResponse(act func(c *Client) error) func(t *testing.T) {
	return func(t *testing.T) {
		// What/Arrange
		srv := startServer(nil, http.StatusServiceUnavailable)
		defer srv.Close()

		underTest := createTestClient(srv)

		// When/Act
		err := act(underTest)

		// Then/Assert
		if err == nil {
			t.Fatal("error expected if API returns an error response")
		}
	}
}
