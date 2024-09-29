package handler

import (
	"github.com/Collap5e-dev/FlickPick/internal/config"
	"strings"
	"testing"
)

func TestHandler_createToken(t *testing.T) {
	expectedResult := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mj"
	secretKey := "test_key"
	username := "ivan"
	h := &Handler{
		config: &config.Config{SecretKey: secretKey},
	}
	got, err := h.createToken(username)
	if err != nil {
		t.Errorf("createToken() error = %v", err)
	}
	if !strings.Contains(got.Token, expectedResult) {
		t.Errorf("createToken() got = %v, want %v", got.Token, expectedResult)
	}

}
