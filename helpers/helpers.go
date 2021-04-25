package helpers

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// HashCRC32 - returns hashed value
func HashCRC32(val string) string {
	var hasher = crc32.NewIEEE()
	hasher.Write([]byte(val))
	return fmt.Sprintf("%x", hasher.Sum32())
}

func GenerateRandomHash() string {
	max := 9000000
	min := 1000000
	// get random number between max and min
	random := rand.Intn(max-min) + min

	// hash it and raturn string
	return HashCRC32(strconv.Itoa(random))
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
