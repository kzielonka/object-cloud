package server

import (
	"errors"
	"io"
	"net/http"

	"github.com/kzielonka/object-cloud/internal/object"
)

func NewServer(store object.Store) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("PUT /objects/{key...}", func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")
		if key == "" {
			http.Error(w, "missing key", http.StatusBadRequest)
			return
		}

		if err := store.Upload(key, r.Body); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	mux.HandleFunc("GET /objects/{key...}", func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")
		if key == "" {
			http.Error(w, "missing key", http.StatusBadRequest)
			return
		}

		fileReader, err := store.Download(key)
		if err != nil {
			if errors.Is(err, object.ErrNotFound) {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		defer fileReader.Close()

		w.WriteHeader(http.StatusOK)
		_, _ = io.Copy(w, fileReader)
	})

	return mux
}
