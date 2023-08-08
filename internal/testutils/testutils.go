package testutils

import (
	"encoding/binary"

	"github.com/dense-analysis/ranges"
	"github.com/google/uuid"
)

// GenerateString create a string of a given length repeating a character.
func GenerateString(char rune, length int) string {
	return ranges.String(ranges.Take[rune](ranges.Repeat(char), 65))
}

func UUIDFromInt(n uint64) uuid.UUID {
	bytes := make([]byte, 16)
	binary.BigEndian.PutUint64(bytes[8:], n)

	return uuid.Must(uuid.FromBytes(bytes))
}
