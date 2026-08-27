package object

import "io"


type Storage interface {
    Save(key string, data io.Reader) error
}

