package totp

import (
	"time"

	"github.com/pquerna/otp/totp"
)

// Generate 生成新秘钥
func Generate(issuer, accountName string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
	})
	if err != nil {
		return "", "", err
	}
	return key.Secret(), key.URL(), nil
}

// GenerateCode 使用当前时间创建TOTP令牌
func GenerateCode(secret string) (string, error) {
	return totp.GenerateCode(secret, time.Now().UTC())
}

// Validate 使用当前时间验证TOTP
func Validate(passcode string, secret string) bool {
	return totp.Validate(passcode, secret)
}
