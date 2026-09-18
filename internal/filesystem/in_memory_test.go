package filesystem_test

import (
	"testing"

	"github.com/kzielonka/object-cloud/internal/filesystem"
	"github.com/kzielonka/object-cloud/internal/object"
	"github.com/kzielonka/object-cloud/internal/object/objecttest"
)

func TestInMemory_Contract(t *testing.T) {
	objecttest.RunFileSystemContract(t, func(t *testing.T) object.FileSystem {
		return filesystem.NewInMemory()
	})
}

func TestInMemory_CreateDirAndHasDirectory(t *testing.T) {
	fs := filesystem.NewInMemory()
	dir1 := "dir1"
	dir2 := "dir2"

	if fs.HasDirectory(dir1) || fs.HasDirectory(dir2) {
		t.Fatalf("expected directories not to exist initially")
	}

	if err := fs.CreateDir(dir1); err != nil {
		t.Fatalf("unexpected error creating dir1: %v", err)
	}
	if !fs.HasDirectory(dir1) {
		t.Errorf("expected dir1 to exist")
	}
	if fs.HasDirectory(dir2) {
		t.Errorf("expected dir2 not to exist yet")
	}

	if err := fs.CreateDir(dir2); err != nil {
		t.Fatalf("unexpected error creating dir2: %v", err)
	}
	if !fs.HasDirectory(dir1) {
		t.Errorf("expected dir1 to still exist")
	}
	if !fs.HasDirectory(dir2) {
		t.Errorf("expected dir2 to exist")
	}
}
