package filesystem_test

import (
	"strings"
	"testing"

	"github.com/kzielonka/object-cloud/internal/filesystem"
	"github.com/kzielonka/object-cloud/internal/object"
	"github.com/kzielonka/object-cloud/internal/object/objecttest"
)

func TestDisk_Contract(t *testing.T) {
	objecttest.RunFileSystemContract(t, func(t *testing.T) object.FileSystem {
		testDir := t.TempDir()
		return filesystem.NewDisk(testDir)
	})
}

func TestDisk_SaveFileErrorHandling(t *testing.T) {
	t.Run("returns error when target directory is unwritable", func(t *testing.T) {
		fs := filesystem.NewDisk("/invalid/non-existent/path")
		err := fs.SaveFile("test-file", strings.NewReader("data"))
		if err == nil {
			t.Fatalf("expected error when directory does not exist, got nil")
		}
	})
}
