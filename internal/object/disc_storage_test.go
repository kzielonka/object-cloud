package object_test

import (
	"github.com/kzielonka/object-cloud/internal/object"
	"testing"
)

func TestDiscFileSystem_Contract(t *testing.T) {
	RunStorageContractTests(t, func(t *testing.T) object.FileSystem {
		testDir := t.TempDir()
		return object.DiscFileSystem(testDir)
	})
}
