package nordigen

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gromson/nordigen/typ"
)

var ErrNoRefreshToken = errors.New("no refresh token configured")
var ErrRefreshTokeExpired = errors.New("refresh token expired")

type tokenRequest struct {
	SecretId  uuid.UUID    `json:"secret_id"`
	SecretKey typ.HexBytes `json:"secret_key"`
}

type tokensResponse struct {
	Access         string `json:"access"`
	AccessExpires  int    `json:"access_expires"`
	Refresh        string `json:"refresh"`
	RefreshExpires int    `json:"refresh_expires"`
}

type refreshRequest struct {
	Refresh string `json:"refresh"`
}

type refreshResponse struct {
	Access        string `json:"access"`
	AccessExpires int    `json:"access_expires"`
}

// authenticate requests an access and refresh token and configures the client.
// In case of API error ApiServerErrorResponse returned
func (n *Nordigen) authenticate() error {
	n.unauthenticate()

	body := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(body).Encode(tokenRequest{SecretId: n.SecretID, SecretKey: n.SecretKey}); err != nil {
		return errors.Wrap(err, "error marshaling new token request")
	}

	tokens := tokensResponse{}
	if err := n.restClient.Exec(http.MethodPost, "/token/new/", body, &tokens); err != nil {
		return errors.Wrap(err, "error executing authentication request")
	}

	n.restClient.Header.Set("Authorization", "Bearer "+tokens.Access)
	n.accessTokenExpiration = time.Now().Add(time.Duration(tokens.AccessExpires) * time.Second).
		Add(n.TokenExpirationBuffer)
	n.RefreshToken = tokens.Refresh
	n.RefreshTokenExpiration = time.Now().Add(time.Duration(tokens.RefreshExpires) * time.Second).
		Add(n.TokenExpirationBuffer)

	return nil
}

// refresh requests a new access token and reconfigures the client.
// In case of API error ApiServerErrorResponse will be returned.
// If the refresh token is not configured in the client ErrNoRefreshToken will be returned.
// If the refresh token is expired ErrRefreshTokeExpired will be returned
func (n *Nordigen) refresh() error {
	if n.RefreshToken == "" {
		return ErrNoRefreshToken
	}

	if n.RefreshTokenExpiration.Unix() < time.Now().Unix() {
		return ErrRefreshTokeExpired
	}

	n.unauthenticate()

	body := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(body).Encode(refreshRequest{Refresh: n.RefreshToken}); err != nil {
		return errors.Wrap(err, "error marshaling a refresh request")
	}

	token := refreshResponse{}
	if err := n.restClient.Exec(http.MethodPost, "/token/refresh", body, &token); err != nil {
		return errors.Wrap(err, "error executing token refresh request")
	}

	n.restClient.Header.Set("Authorization", "Bearer "+token.Access)
	n.accessTokenExpiration = time.Now().Add(time.Duration(token.AccessExpires) * time.Second).
		Add(n.TokenExpirationBuffer)

	return nil
}

func (n *Nordigen) unauthenticate() {
	n.restClient.Header.Del("Authorization")
	n.accessTokenExpiration = time.Unix(0, 0)
}

func (n *Nordigen) ensureAuthenticated() error {
	if n.restClient.Header.Get("Authorization") != "" && n.accessTokenExpiration.Unix() > time.Now().Unix() {
		return nil
	}

	if n.RefreshToken != "" && n.RefreshTokenExpiration.Unix() > time.Now().Unix() {
		return n.refresh()
	}

	return n.authenticate()
}
