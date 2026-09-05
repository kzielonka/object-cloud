package object

type KeyHasher interface {
	Hash(key string) []byte
}
