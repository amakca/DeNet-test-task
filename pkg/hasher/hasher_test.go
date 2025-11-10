package hasher

import (
	"crypto/sha1"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSHA1Hasher_Hash_IsDeterministic(t *testing.T) {
	h := NewSHA1Hasher("salt")
	got1 := h.Hash("password")
	got2 := h.Hash("password")
	assert.Equal(t, got1, got2, "expected deterministic hash")
}

func TestSHA1Hasher_Hash_MatchesAlgorithm(t *testing.T) {
	h := NewSHA1Hasher("pepper")
	got := h.Hash("p@ssW0rd")

	hash := sha1.New()
	hash.Write([]byte("p@ssW0rd"))
	expected := fmt.Sprintf("%x", hash.Sum([]byte("pepper")))

	assert.Equal(t, expected, got, "unexpected hash")
}
