package auth_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/kzielonka/object-cloud/internal/auth"
)

func TestPresignedUrl_Verify(t *testing.T) {
	presigned_url := auth.NewUploadPolicy("key", auth.ExpiresInSec(3600))
	timeNow, _ := time.Parse(time.Kitchen, "10:00AM")
	token, err := presigned_url.Issue("secret", timeNow)
	if err != nil {
		t.Fatalf("can not issue token")
	}
	decrypted_token := auth.NewUploadToken(token).Decrypt("secret")
	if decrypted_token.IsValid(timeNow) != true {
		t.Fatalf("token should be valid")
	}
	if decrypted_token.Key() != "key" {
		t.Fatalf("expected key %s but got %s", "key", decrypted_token.Key())
	}
}

func TestPresignedUrl_FailsOnInvalidSecret(t *testing.T) {
	presigned_url := auth.NewUploadPolicy("key", auth.ExpiresInSec(3600))
	timeNow, _ := time.Parse(time.Kitchen, "10:00AM")
	token, err := presigned_url.Issue("secret", timeNow)
	if err != nil {
		t.Fatalf("can not issue token")
	}
	decrypted_token := auth.NewUploadToken(token).Decrypt("secret2")
	if decrypted_token.IsValid(timeNow) != false {
		t.Fatalf("token should be invalid")
	}
}

func TestPresignedUrl_FailsOnTemperedKey(t *testing.T) {
	secret := "secret"
	presigned_url := auth.NewUploadPolicy("key", auth.ExpiresInSec(3600))
	timeNow, _ := time.Parse(time.Kitchen, "10:00AM")
	token, err := presigned_url.Issue(secret, timeNow)
	if err != nil {
		t.Fatalf("can not issue token")
	}
	decrypted_token := auth.NewUploadToken(token).Decrypt(secret)
	if decrypted_token.IsValid(timeNow) == false {
		t.Fatalf("token should be valid")
	}
	var objmap map[string]interface{}
	if err := json.Unmarshal(token, &objmap); err != nil {
		t.Fatalf("can not marshall token")
	}
	objmap["key"] = "new key"
	temperedToken, _ := json.Marshal(objmap)
	decrypted_token = auth.NewUploadToken(temperedToken).Decrypt(secret)
	if decrypted_token.IsValid(timeNow) != false {
		t.Fatalf("token should be invalid")
	}
}

func TestPresignedUrl_FailsWhenExpired(t *testing.T) {
	secret := "secret"
	presigned_url := auth.NewUploadPolicy("key", auth.ExpiresInSec(5*60))
	timeNow, _ := time.Parse(time.Kitchen, "10:00AM")
	token, err := presigned_url.Issue(secret, timeNow)
	if err != nil {
		t.Fatalf("can not issue token")
	}
	decrypted_token := auth.NewUploadToken(token).Decrypt(secret)
	if decrypted_token.IsValid(timeNow.Add(4*time.Minute)) == false {
		t.Fatalf("token should be valid")
	}
	if decrypted_token.IsValid(timeNow.Add(4*time.Minute+59*time.Second)) == false {
		t.Fatalf("token should be valid")
	}
	if decrypted_token.IsValid(timeNow.Add(5*time.Minute+1*time.Second)) == true {
		t.Fatalf("token should be invalid")
	}
	if decrypted_token.IsValid(timeNow.Add(time.Duration(-1)*time.Second)) == true {
		t.Fatalf("token should be invalid")
	}
}
