package object_test

import (
	"testing"
	"strings"
	"github.com/kzielonka/object-cloud/internal/object"
	"io"
)

func TestInMemoryStorage_ReturnErrorWhenKeyIsNotSet(t *testing.T) {
  	store := object.InMemoryStorage()
    err, _ := store.Load("key")
		if err != nil {
			t.Fatalf("error expected")
  	}
}

func TestInMemoryStorage_SetsKey(t *testing.T) {
  	store := object.InMemoryStorage()
		initialData := "hello"
	  store.Save("key", strings.NewReader(initialData)) 
    _, stream := store.Load("key")
		data := io.ReadAll(stream)
		if data != initialData {
			t.Errorf("data %s is not equal to %s", data, initialData)
		}
}
