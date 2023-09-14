package totp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateAndValidateTOTP(t *testing.T) {
	secret, url, err := Generate("ziy", "ziy@163.com")
	require.NoError(t, err, "generate basic TOTP")
	require.Equal(t, 32, len(secret), "Secret is 32 bytes long as base32.")
	require.Contains(t, url, "issuer=ziy")

	passcode, err := GenerateCode(secret)
	require.NoError(t, err)

	valid := Validate(passcode, secret)
	require.True(t, valid)
}
