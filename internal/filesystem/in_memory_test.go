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
