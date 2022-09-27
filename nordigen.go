package nordigen

import (
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gromson/nordigen/rest"
	"gromson/nordigen/typ"
)

const (
	baseUrl        = "https://ob.nordigen.com/api"
	defaultVersion = "v2"

	defaultTokenExpirationBuffer = -1 * time.Minute
)

// Nordigen for accessing Nordigen API
type Nordigen struct {
	// SecretID for accessing Nordigen API
	SecretID uuid.UUID
	// SecretKey decoded hex-string value
	SecretKey              []byte
	RefreshToken           string
	RefreshTokenExpiration time.Time
	// TokenExpirationBuffer for treating a token as expired earlier than it actually become expired
	// to compensate the delay between receiving it from the API and setting expiration date to a client.
	// The value must be negative
	TokenExpirationBuffer time.Duration
	accessToken           string
	accessTokenExpiration time.Time
	// BaseUrl of restClient must be "https://ob.nordigen.com/api".
	restClient *rest.Client
}

// New creates new Nordigen client.
func New(secretID string, secretKey string) (*Nordigen, error) {
	secretUUID, err := uuid.Parse(secretID)
	if err != nil {
		return nil, errors.Wrap(err, "invalid secretID format provided")
	}

	if secretKey == "" {
		return nil, errors.New("secretKey can't be empty")
	}

	secretKeyData := typ.HexBytes{}
	if err := secretKeyData.UnmarshalText([]byte(secretKey)); err != nil {
		return nil, errors.Wrap(err, "invalid secret key format")
	}

	return MustNew(secretUUID, secretKeyData), nil
}

// MustNew creates new Nordigen client. Secret key is a decoded hex-string i.e. "ff2a24" -> [255, 42, 36].
func MustNew(secretID uuid.UUID, secretKey []byte) *Nordigen {
	apiUrl, err := url.Parse(baseUrl + "/" + defaultVersion)
	if err != nil {
		panic(err)
	}

	return &Nordigen{
		SecretID:               secretID,
		SecretKey:              secretKey,
		RefreshToken:           "",
		RefreshTokenExpiration: time.Unix(0, 0),
		TokenExpirationBuffer:  defaultTokenExpirationBuffer,
		accessTokenExpiration:  time.Unix(0, 0),
		restClient:             rest.NewClient(apiUrl, nil),
	}
}
