package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kzielonka/object-cloud/internal/filesystem"
	"github.com/kzielonka/object-cloud/internal/object"
	"github.com/kzielonka/object-cloud/internal/server"
)

func setUpStore(t *testing.T) object.Store {
	t.Helper()
	store, err := object.NewStore(
		object.WithFileSystem(filesystem.NewInMemory()),
		object.WithDir("/dir"),
	)
	if err != nil {
		t.Fatalf("received error %s", err)
	}
	return store
}

func TestServer_UploadsAndDownloadsFile(t *testing.T) {
	store := setUpStore(t)
	s := server.NewServer(store)
	fileContent := "file content"
	fileContentReader := strings.NewReader(fileContent)

	// 1. Upload
	req := httptest.NewRequest(http.MethodPut, "/objects/key", fileContentReader)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	// 2. Download
	req = httptest.NewRequest(http.MethodGet, "/objects/key", nil)
	rec2 := httptest.NewRecorder()
	s.ServeHTTP(rec2, req)

	if rec2.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec2.Code)
	}

	if rec2.Body.String() != fileContent {
		t.Fatalf("expected body to be %q, but got %q", fileContent, rec2.Body.String())
	}
}

func TestServer_GetNonExistentFile(t *testing.T) {
	store := setUpStore(t)
	s := server.NewServer(store)

	req := httptest.NewRequest(http.MethodGet, "/objects/missing-key", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}

	if got := strings.TrimSpace(rec.Body.String()); got != "not found" {
		t.Fatalf("expected body to be %q, but got %q", "not found", got)
	}
}
