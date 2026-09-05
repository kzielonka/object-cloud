package object_test

import (
	"bytes"
	"testing"

	"github.com/kzielonka/object-cloud/internal/object"
)

type KeyHasherFactory func(t *testing.T) object.KeyHasher

func RunKeyHasherContract(t *testing.T, newKeyHasher KeyHasherFactory) {
	t.Run("returns same hash for same key", func(t *testing.T) {
		hasher := newKeyHasher(t)
		hash1 := hasher.Hash("key")
		hash2 := hasher.Hash("key")
		if !bytes.Equal(hash1, hash2) {
			t.Fatalf("expected %s to equal %s", hash1, hash2)
		}
	})

	t.Run("returns different hash for different keys", func(t *testing.T) {
		hasher := newKeyHasher(t)
		hash1 := hasher.Hash("key1")
		hash2 := hasher.Hash("key2")
		if bytes.Equal(hash1, hash2) {
			t.Fatalf("expected %s to not equal %s", hash1, hash2)
		}
	})
}
