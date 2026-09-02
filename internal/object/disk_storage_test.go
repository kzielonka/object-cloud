package object_test

import (
	"github.com/kzielonka/object-cloud/internal/object"
	"testing"
)

func TestDiskFileSystem_Contract(t *testing.T) {
	RunFileSystemContract(t, func(t *testing.T) object.FileSystem {
		testDir := t.TempDir()
		return object.NewDiskFileSystem(testDir)
	})
}
