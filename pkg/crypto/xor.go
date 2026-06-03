package crypto

import (
	"encoding/base64"
)

func DecryptXOR(b64 string, key []byte) (string, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}
	for i := range data {
		data[i] ^= key[i%len(key)]
	}
	return string(data), nil
}
