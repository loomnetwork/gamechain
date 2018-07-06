package battleground

import (
	"crypto/rand"

	"github.com/google/uuid"
)

const UUIDBytes = 16

func generateUUID() (string, error) {
	buffer := make([]byte, UUIDBytes)

	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	if genUuid, err := uuid.FromBytes(buffer); err != nil {
		return "", err
	} else {
		return genUuid.String(), nil
	}

}
