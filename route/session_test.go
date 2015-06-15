package route

import (
	"crypto/hmac"
	"crypto/sha256"
	"testing"
)

func TestComputeHMAC(t *testing.T) {
	msgMAC := computeHMAC(TestMsg, TestKey)
	mac := hmac.New(sha256.New, []byte(TestKey))
	mac.Write([]byte(TestMsg))
	expectedMAC := mac.Sum(nil)
	if !hmac.Equal(msgMAC, expectedMAC) {
		t.Errorf("want %v\n, got %v", expectedMAC, msgMAC)
	}
}

func TestCheckMAC(t *testing.T) {
	mac := hmac.New(sha256.New, []byte(TestKey))
	mac.Write([]byte(TestMsg))
	expectedMAC := mac.Sum(nil)
	if !checkMAC(TestMsg, expectedMAC, TestKey) {
		t.Errorf("want %b, got %b", true, false)
	}
}
