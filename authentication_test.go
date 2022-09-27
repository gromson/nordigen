package nordigen

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

const (
	testAccessToken         = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJncm9tc29uL25vcmRpZ2VuIiwiaWF0IjoxNjU4MTczNzU3LCJleHAiOjE2ODk3MDk3MTMsImF1ZCI6InRlc3QiLCJzdWIiOiJ0ZXN0IiwidG9rZW5fdHlwZSI6ImFjY2VzcyIsImp0aSI6IjYwNDEzODUwZDM1MzJjYWQxODI5ZDc2OTQ5OWNlZGRlIiwiaWQiOiIxMjM0NSIsInNlY3JldF9pZCI6IjA1MzZhNmZiLWE0OGYtNDI4Ni1hMzJjLWU0Y2NkYzkxMTBiOCJ9.3fRUr5fn5kcB_n6Sdh1ggBThuTN9TuKcvn9P12TDmkI"
	testAccessTokenExpires  = 86400
	testRefreshToken        = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJncm9tc29uL25vcmRpZ2VuIiwiaWF0IjoxNjU4MTczNzU3LCJleHAiOjE2ODk3MDk3MTMsImF1ZCI6InRlc3QiLCJzdWIiOiJ0ZXN0IiwidG9rZW5fdHlwZSI6InJlZnJlc2giLCJqdGkiOiI2MDQxMzg1MGQzNTMyY2FkMTgyOWQ3Njk0OTljZWRmZiIsImlkIjoiMTIzNDUiLCJzZWNyZXRfaWQiOiIwNTM2YTZmYi1hNDhmLTQyODYtYTMyYy1lNGNjZGM5MTEwYjgifQ.uiR0k7a1W_67AFkhJLYQTqS5fKZtX1aWcScZPp3hbt8"
	testRefreshTokenExpires = 2592000
)

const (
	testUnauthenticatedResponsePayload = `{"summary":"Invalid token","detail":"Token is invalid or expired","status_code":401}`
)

func TestClient_authentication(t *testing.T) {
	t.Parallel()
	t.Run("sucessful authentication", testAuthenticationOk)
	t.Run("401 response", testAuthentication401)
}

func TestClient_refresh(t *testing.T) {
	t.Parallel()
	t.Run("successful refresh", testRefreshOk)
	t.Run("401 refresh", testRefresh401)
	t.Run("refresh without a refresh token", testRefreshWithoutRefreshToken)
	t.Run("refresh with an expired refresh token", testRefreshWithExpiredRefreshToken)
}

func testAuthenticationOk(t *testing.T) {
	// What/Arrange
	responsePayload := fmt.Sprintf(
		`{"access":"%s","access_expires":%d,"refresh":"%s","refresh_expires":%d}`,
		testAccessToken,
		testAccessTokenExpires,
		testRefreshToken,
		testRefreshTokenExpires)

	srv := startServer(responsePayload, http.StatusOK)
	defer srv.Close()

	underTest := createTestNordigen(srv)

	// Then/Assert
	if err := underTest.authenticate(); err != nil {
		t.Fatalf("authentication failed: %s", err)
	}

	// Then/Assert
	if strings.Split(underTest.restClient.Header.Get("Authorization"), " ")[1] != testAccessToken {
		t.Fatal("access token wasn't set")
	}

	if underTest.accessTokenExpiration.Unix() < time.Now().Unix() {
		t.Fatal("access token expiration wasn't set")
	}

	if underTest.RefreshToken != testRefreshToken {
		t.Fatal("refresh token wasn't set")
	}

	if underTest.RefreshTokenExpiration.Unix() < time.Now().Unix() {
		t.Fatal("refresh token expiration wasn't set")
	}
}

func testAuthentication401(t *testing.T) {
	srv := startServer(testUnauthenticatedResponsePayload, http.StatusUnauthorized)
	defer srv.Close()

	underTest := createTestNordigen(srv)

	err := underTest.authenticate()

	if err == nil {
		t.Fatal("error expected on 401 server response")
	}
}

func testRefreshOk(t *testing.T) {
	// What/Arrange
	responsePayload := fmt.Sprintf(`{"access":"%s","access_expires":%d}`, testAccessToken, testAccessTokenExpires)
	srv := startServer(responsePayload, http.StatusOK)
	defer srv.Close()

	underTest := createTestNordigen(srv)
	underTest.accessTokenExpiration = time.Unix(0, 0)
	underTest.RefreshToken = testRefreshToken
	underTest.RefreshTokenExpiration = time.Now().Add(time.Duration(testRefreshTokenExpires) * time.Second)

	// When/Act
	if err := underTest.refresh(); err != nil {
		t.Fatalf("refresh failed: %s", err)
	}

	// Then/Assert
	if strings.Split(underTest.restClient.Header.Get("Authorization"), " ")[1] != testAccessToken {
		t.Fatalf("token hasn't been exchanged")
	}

	if underTest.accessTokenExpiration.Unix() < time.Now().Unix() {
		t.Fatal("access token expiration wasn't set")
	}
}

func testRefresh401(t *testing.T) {
	srv := startServer(testUnauthenticatedResponsePayload, http.StatusUnauthorized)
	defer srv.Close()

	underTest := createTestNordigen(srv)
	underTest.accessTokenExpiration = time.Unix(0, 0)
	underTest.RefreshToken = testRefreshToken
	underTest.RefreshTokenExpiration = time.Now().Add(time.Duration(testRefreshTokenExpires) * time.Second)

	err := underTest.refresh()

	if err == nil {
		t.Fatal("error expected on 401 server response")
	}
}

func testRefreshWithoutRefreshToken(t *testing.T) {
	// What/Arrange
	underTest := MustNew(uuid.New(), []byte{12, 23, 34, 45, 56, 67, 78, 89, 90})
	underTest.accessTokenExpiration = time.Unix(0, 0)
	underTest.RefreshToken = ""
	underTest.RefreshTokenExpiration = time.Now().Add(24 * time.Hour)

	// When/Act
	err := underTest.refresh()

	// Then/Assert
	if err == nil {
		t.Fatal("error expected when trying refresh without a refresh token")
	}

	if err != ErrNoRefreshToken {
		t.Fatalf(
			`"%s" error expected when trying refresh without a refresh token: "%v" received`,
			ErrNoRefreshToken,
			err)
	}
}

func testRefreshWithExpiredRefreshToken(t *testing.T) {
	// What/Arrange
	underTest := MustNew(uuid.New(), []byte{12, 23, 34, 45, 56, 67, 78, 89, 90})
	underTest.accessTokenExpiration = time.Unix(0, 0)
	underTest.RefreshToken = testRefreshToken
	underTest.RefreshTokenExpiration = time.Now().Add(-1 * time.Second)

	// When/Act
	err := underTest.refresh()

	// Then/Assert
	if err == nil {
		t.Fatal("error expected when trying refresh with an expired refresh token")
	}

	if err != ErrRefreshTokeExpired {
		t.Fatalf(
			`"%s" error expected when trying refresh with an expired refresh token: "%v" received`,
			ErrNoRefreshToken,
			err)
	}
}
