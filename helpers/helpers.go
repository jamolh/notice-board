package helpers

import (
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// GetEnv - lookup by key
// if fatal=true and didn't find
// by key anything finish with
// fatal error
func GetEnv(key string, fatal bool) string {
	value, found := os.LookupEnv(key)
	if !found {
		log.Println("LookupEnv not found by key:", key)
		if fatal {
			os.Exit(1)
		}
	}

	return value
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
