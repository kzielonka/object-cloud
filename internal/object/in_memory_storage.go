package object

import "io"
import "fmt"
import "errors"
import "bytes"


type inMemoryStorage struct {
    savedFiles map[string][]byte  
}

func InMemoryStorage() *inMemoryStorage {
	  return &inMemoryStorage { 
  	  	savedFiles: make(map[string][]byte),
	  }
}

func (s *inMemoryStorage) Save(key string, data io.Reader) error {
  	bytes, err := io.ReadAll(data)
  	if err != nil {
    		return fmt.Errorf("failed to read data: %w", err)
  	}
  	s.savedFiles[key] = bytes
		return nil
}

func (s *inMemoryStorage) Load(key string) (error, io.Reader) {
	data := s.savedFiles[key]
	if data != nil {
		return errors.New("no object"), nil
	}
  return nil, bytes.NewReader(data)
}
