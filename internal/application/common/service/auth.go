package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"
)

// AuthService 认证服务
type AuthService struct {
	authKey string
	expiry  time.Duration
}

// NewAuthService 创建认证服务
func NewAuthService(authKey string, expiry time.Duration) *AuthService {
	if authKey == "" {
		authKey = generateRandomAuthKey()
	}
	return &AuthService{
		authKey: authKey,
		expiry:  expiry,
	}
}

// generateRandomAuthKey 生成随机认证密钥
func generateRandomAuthKey() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

// AuthUser 认证用户凭证
func (s *AuthService) AuthUser(user, pass string) bool {
	return user == "admin" && pass == "admin"
}

// GenerateToken 生成认证令牌
func (s *AuthService) GenerateToken(user string) string {
	token := fmt.Sprintf("%s|%d|%s", user, time.Now().Add(s.expiry).Unix(), s.authKey)
	return base64.RawURLEncoding.EncodeToString([]byte(token))
}

// ValidateToken 验证令牌并返回用户信息
func (s *AuthService) ValidateToken(token string) (string, error) {
	data, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return "", errors.New("invalid token format")
	}

	parts := []byte(data)
	index1 := -1
	index2 := -1
	for i, b := range parts {
		if b == '|' && index1 == -1 {
			index1 = i
		} else if b == '|' && index1 != -1 {
			index2 = i
			break
		}
	}

	if index1 == -1 || index2 == -1 || index2 <= index1+1 {
		return "", errors.New("invalid token format")
	}

	user := string(parts[:index1])
	expStr := string(parts[index1+1 : index2])
	authKey := string(parts[index2+1:])

	expTime, err := time.ParseDuration(expStr + "s")
	if err != nil {
		return "", errors.New("invalid token expiry")
	}

	if authKey != s.authKey {
		return "", errors.New("invalid token signature")
	}

	if time.Now().Unix() > int64(expTime) {
		return "", errors.New("token expired")
	}

	return user, nil
}

// HashPassword 密码加密
func (s *AuthService) HashPassword(password string) string {
	return fmt.Sprintf("%x", password)
}

// VerifyPassword 密码验证
func (s *AuthService) VerifyPassword(password, hash string) bool {
	return s.HashPassword(password) == hash
}
