package filesystem_test

import (
	"testing"

	"github.com/kzielonka/object-cloud/internal/filesystem"
	"github.com/kzielonka/object-cloud/internal/object"
)

func TestInMemory_Contract(t *testing.T) {
	object.RunFileSystemContract(t, func(t *testing.T) object.FileSystem {
		return filesystem.NewInMemory()
	})
}
