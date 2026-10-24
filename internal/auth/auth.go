package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"time"
)

type uploadPolicy struct {
	key          string
	expiresInSec int
}

type signedUploadPolicy struct {
	uploadPolicy
	signature []byte
}

type materializedUploadPolicy struct {
	Key       string    `json:"key"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type materializedSignedUploadPolicy struct {
	materializedUploadPolicy
	Signature []byte `json:"signature"`
}

func NewUploadPolicy(key string, opts ...Option) *uploadPolicy {
	policy := &uploadPolicy{
		key:          key,
		expiresInSec: 5 * 60,
	}
	for _, opt := range opts {
		opt(policy)
	}
	return policy
}

type Option func(*uploadPolicy)

func ExpiresInSec(time int) Option {
	return func(up *uploadPolicy) {
		up.expiresInSec = time
	}
}

func generateSignature(secret string, key string, issuedAt time.Time, expiresAt time.Time) ([]byte, error) {
	materializedUp := materializedUploadPolicy{
		Key:       key,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}
	json1, err := json.Marshal(materializedUp)
	if err != nil {
		return nil, errors.New("invalid json")
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(json1)
	return mac.Sum(nil), nil
}

func (up *uploadPolicy) Issue(secret string, now time.Time) ([]byte, error) {
	expiresAt := now.Add(time.Duration(up.expiresInSec) * time.Second)

	signature, err := generateSignature(secret, up.key, now, expiresAt)
	if err != nil {
		return nil, errors.New("can not generate signature")
	}
	materializedSup := materializedSignedUploadPolicy{
		Key:       up.key,
		IssuedAt:  now,
		ExpiresAt: expiresAt,
		Signature: signature,
	}
	json2, err := json.Marshal(materializedSup)
	if err != nil {
		return nil, errors.New("invalid signed json")
	}
	return json2, nil
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
	issuedAt  time.Time
	expiresAt time.Time
}

func (t *uploadToken) Decrypt(secret string) *decryptedToken {
	var decodedToken = materializedSignedUploadPolicy{}
	err := json.Unmarshal(t.token, &decodedToken)
	if err != nil {
		return &decryptedToken{
			isValid:   false,
			key:       "",
			issuedAt:  time.Now(),
			expiresAt: time.Now(),
		}
	}
	signature, err := generateSignature(secret, decodedToken.Key, decodedToken.IssuedAt, decodedToken.ExpiresAt)
	if err != nil || !hmac.Equal(signature, decodedToken.Signature) {
		return &decryptedToken{
			isValid:   false,
			key:       decodedToken.Key,
			issuedAt:  decodedToken.IssuedAt,
			expiresAt: decodedToken.ExpiresAt,
		}
	}
	return &decryptedToken{
		isValid:   true,
		key:       decodedToken.Key,
		issuedAt:  decodedToken.IssuedAt,
		expiresAt: decodedToken.ExpiresAt,
	}
}

func (dt *decryptedToken) IsValid(atTime time.Time) bool {
	return dt.isValid && (dt.issuedAt.Before(atTime) || dt.issuedAt.Equal(atTime)) && (dt.expiresAt.After(atTime) || dt.expiresAt.Equal(atTime))
}

func (dt *decryptedToken) Key() string {
	return dt.key
}
