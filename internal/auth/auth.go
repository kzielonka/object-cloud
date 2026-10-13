package auth

import (
	"encoding/json"
	"errors"
	"time"
)

type uploadPolicy struct {
	key          string
	expiresInSec int
}

type materializedUploadPolicy struct {
	Key       string    `json:"key"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewUploadPolicy(key string, opts ...Option) *uploadPolicy {
	return &uploadPolicy{
		key: key,
	}
}

type Option func(*uploadPolicy)

func ExpiresInSec(time int) Option {
	return func(up *uploadPolicy) {
		up.expiresInSec = time
	}
}

func (up *uploadPolicy) Issue(secret string, now time.Time) ([]byte, error) {
	materialized_up := materializedUploadPolicy{
		Key:       up.key,
		ExpiresAt: now.Add(time.Duration(up.expiresInSec) * time.Second),
	}
	json, err := json.Marshal(materialized_up)
	if err != nil {
		return nil, errors.New("wtf")
	}
	return json, nil
}

type uploadToken struct {
	token []byte
}

func NewUploadToken(token []byte) *uploadToken {
	return &uploadToken{
		token: token,
	}
}

type decryptedToken struct {
	isValid   bool
	key       string
	expiresAt time.Time
}

func (t *uploadToken) Decrypt(secret string) *decryptedToken {
	var decodedToken = materializedUploadPolicy{}
	err := json.Unmarshal(t.token, &decodedToken)
	if err != nil {
		return &decryptedToken{
			isValid:   false,
			key:       "",
			expiresAt: time.Now(),
		}

	}
	return &decryptedToken{
		isValid:   true,
		key:       decodedToken.Key,
		expiresAt: time.Now(),
	}
}

func (dt *decryptedToken) IsValid() bool {
	return dt.isValid
}

func (dt *decryptedToken) Key() string {
	return dt.key
}
