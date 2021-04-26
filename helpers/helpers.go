package helpers

import (
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// RemoveNonLetter - remove all expect letters
func RemoveNonLetter(s string) string {
	return strings.Map(
		func(r rune) rune {
			if unicode.IsLetter(r) {
				return r
			}
			return -1
		},
		s,
	)
}

func IsValidUUID(s string) bool {
	if _, err := uuid.Parse(s); err != nil {
		return false
	}
	return true
}
