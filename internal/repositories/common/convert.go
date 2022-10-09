package common

import (
	"encoding/base64"

	"github.com/google/uuid"
)

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

func EncodeUUIDToBase64(uuid uuid.UUID) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(uuid.String()))
	return encoded
}
