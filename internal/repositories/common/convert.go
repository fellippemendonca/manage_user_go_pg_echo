package common

import (
	"encoding/base64"

	"github.com/google/uuid"
)

// DecodeBase64ToUUID is responsible to unmask UUIDs previously converted to base 64
func DecodeBase64ToUUID(encoded string) (uuid.UUID, error) {
	if encoded == "" {
		return uuid.Nil, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return uuid.UUID{}, err
	}

	return uuid.Parse(string(decoded))
}

// EncodeBase64ToUUID is responsible to mask UUIDs to base 64
func EncodeUUIDToBase64(uuid uuid.UUID) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(uuid.String()))
	return encoded
}
