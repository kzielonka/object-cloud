package object_test

import (
	"testing"

	"github.com/kzielonka/object-cloud/internal/object"
)

func TestSHA256KeyHasher(t *testing.T) {
	RunKeyHasherContract(t, func(t *testing.T) object.KeyHasher {
		return object.NewSHA256KeyHasher()
	})
}
