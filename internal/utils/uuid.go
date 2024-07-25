package utils

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateShortID() string {
	newUUID := uuid.New()
	// Will replace '-' and return only 8 chars of uuid
	return strings.Replace(newUUID.String()[:8], "-", "", -1)
}
