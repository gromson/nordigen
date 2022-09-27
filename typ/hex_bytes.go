package typ

import (
	"encoding/hex"
	"encoding/json"

	"github.com/pkg/errors"
)

// HexBytes byte slice that marshals to and unmarshals from a hex-encoded string
type HexBytes []byte

func (v HexBytes) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(v)), nil
}

func (v *HexBytes) UnmarshalText(text []byte) error {
	dst := make([]byte, hex.DecodedLen(len(text)))
	n, err := hex.Decode(dst, text)
	if err != nil {
		return err
	}

	*v = dst[:n]

	return nil
}

func (v HexBytes) MarshalJSON() ([]byte, error) {
	text, err := v.MarshalText()
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling HexBytes value to JSON")
	}

	return json.Marshal(string(text))
}

func (v *HexBytes) UnmarshalJSON(b []byte) error {
	return v.UnmarshalText(b[1 : len(b)-1])
}
