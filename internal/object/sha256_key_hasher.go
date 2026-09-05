package object

import (
	"crypto/sha256"
)

type sha256KeyHasher struct{}

func NewSHA256KeyHasher() *sha256KeyHasher {
	return &sha256KeyHasher{}
}

func (kh *sha256KeyHasher) Hash(key string) []byte {
	sum := sha256.Sum256([]byte(key))
	return sum[:]
}
