package util

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// GenerateUUIDOrFallback returns a UUID string.
// If UUID generation fails, it falls back to a timestamp-based string like "uuid_fallback_1751942234747270000".
func GenerateUUIDSafe() string {
	id, err := uuid.NewRandom()
	if err != nil {
		fallback := fmt.Sprintf("uuid_fallback_%d", time.Now().UnixNano())
		log.Printf("WARN: UUID generation failed, using fallback ID: %s (error: %v)", fallback, err)
		return fallback
	}
	return id.String()
}
