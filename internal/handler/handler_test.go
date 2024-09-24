package handler

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Collap5e-dev/FlickPick/internal/config"
)

func TestHandler_createToken(t *testing.T) {
	expect := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3M"
	//expe := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc4MTMxMjksInVzZXJuYW1lIjoidXNlcm5hbWUifQ.JDlmUYCG1-0r"
	secretKey := "secretKey"
	username := "username"
	h := &Handler{
		config: &config.Config{SecretKey: secretKey},
	}
	got, err := h.createToken(username)
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(got, expect))
}
