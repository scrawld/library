package telegramauth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Credentials 存储 Telegram 登录授权数据
type Credentials struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	PhotoURL  string `json:"photo_url"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
}

var (
	ErrInvalidHash  = errors.New("telegram: hash mismatch")
	ErrAuthTooOld   = errors.New("telegram: auth_date is too old")
	ErrMissingToken = errors.New("telegram: bot token is empty")
)

// dataCheckString 按 Telegram 规范构造待校验字符串。
// 只包含非空字段，按字段名字母序排列，以 "\n" 连接。
func (c *Credentials) dataCheckString() string {
	parts := make([]string, 0, 6)
	parts = append(parts,
		"id="+strconv.FormatInt(c.ID, 10),
		"first_name="+c.FirstName,
		"auth_date="+strconv.FormatInt(c.AuthDate, 10),
	)

	if c.LastName != "" {
		parts = append(parts, "last_name="+c.LastName)
	}
	if c.Username != "" {
		parts = append(parts, "username="+c.Username)
	}
	if c.PhotoURL != "" {
		parts = append(parts, "photo_url="+c.PhotoURL)
	}

	sort.Strings(parts)
	return strings.Join(parts, "\n")
}

// Verify 校验 Telegram 授权数据的合法性
// token: Telegram Bot Token
// maxAge: 允许的最大授权时效，0 表示不限制
func (c *Credentials) Verify(token string, maxAge time.Duration) error {
	if token == "" {
		return ErrMissingToken
	}

	if maxAge > 0 && time.Since(time.Unix(c.AuthDate, 0)) > maxAge {
		return ErrAuthTooOld
	}

	secret := sha256.Sum256([]byte(token))
	mac := hmac.New(sha256.New, secret[:])
	mac.Write([]byte(c.dataCheckString()))

	hashBytes, _ := hex.DecodeString(c.Hash)
	if !hmac.Equal(mac.Sum(nil), hashBytes) {
		return ErrInvalidHash
	}
	return nil
}
