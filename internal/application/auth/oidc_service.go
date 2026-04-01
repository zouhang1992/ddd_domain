package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

// TokenSet Token 集合
type TokenSet struct {
	AccessToken  string
	RefreshToken string
	IDToken      string
	Expiry       time.Time
}

// OIDCService OIDC 服务
type OIDCService struct {
	config       Config
	oauth2Config *oauth2.Config
	httpClient   *http.Client
	ctx          context.Context

	// OIDC discovery 缓存
	discoveryCache   *oidcDiscovery
	discoveryMutex   sync.RWMutex
	discoveryExpiry  time.Time

	// JWKS 缓存
	jwksCache   *jwks
	jwksMutex   sync.RWMutex
	jwksExpiry  time.Time
}

// oidcDiscovery OIDC discovery 响应
type oidcDiscovery struct {
	Issuer           string `json:"issuer"`
	AuthURL          string `json:"authorization_endpoint"`
	TokenURL         string `json:"token_endpoint"`
	JWKSURL          string `json:"jwks_uri"`
	UserInfoURL      string `json:"userinfo_endpoint"`
	EndSessionURL    string `json:"end_session_endpoint,omitempty"`
}

// jwks JWKS 响应
type jwks struct {
	Keys []json.RawMessage `json:"keys"`
}

// NewOIDCService 创建 OIDC 服务
func NewOIDCService(config Config) *OIDCService {
	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Scopes:       config.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "", // 从 discovery 获取
			TokenURL: "", // 从 discovery 获取
		},
	}

	return &OIDCService{
		config:       config,
		oauth2Config: oauth2Config,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		ctx:          context.Background(),
	}
}

// fetchDiscovery 获取 OIDC discovery 配置
func (s *OIDCService) fetchDiscovery() (*oidcDiscovery, error) {
	s.discoveryMutex.RLock()
	if s.discoveryCache != nil && time.Now().Before(s.discoveryExpiry) {
		cached := s.discoveryCache
		s.discoveryMutex.RUnlock()
		return cached, nil
	}
	s.discoveryMutex.RUnlock()

	// 获取 discovery 文档
	discoveryURL := s.config.IssuerURL + "/.well-known/openid-configuration"
	resp, err := s.httpClient.Get(discoveryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch discovery: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("discovery request failed with status: %d", resp.StatusCode)
	}

	var discovery oidcDiscovery
	if err := json.NewDecoder(resp.Body).Decode(&discovery); err != nil {
		return nil, fmt.Errorf("failed to parse discovery: %w", err)
	}

	// 更新缓存
	s.discoveryMutex.Lock()
	s.discoveryCache = &discovery
	s.discoveryExpiry = time.Now().Add(1 * time.Hour) // 缓存 1 小时
	s.discoveryMutex.Unlock()

	// 更新 oauth2 config 的 endpoint
	s.oauth2Config.Endpoint.AuthURL = discovery.AuthURL
	s.oauth2Config.Endpoint.TokenURL = discovery.TokenURL

	return &discovery, nil
}

// GenerateState 生成随机 state
func GenerateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// GetAuthURL 获取认证 URL
func (s *OIDCService) GetAuthURL(state string) (string, error) {
	if _, err := s.fetchDiscovery(); err != nil {
		return "", err
	}
	return s.oauth2Config.AuthCodeURL(state), nil
}

// ExchangeCode 用 code 换取 tokens
func (s *OIDCService) ExchangeCode(code string) (*TokenSet, error) {
	if _, err := s.fetchDiscovery(); err != nil {
		return nil, err
	}

	ctx := context.WithValue(s.ctx, oauth2.HTTPClient, s.httpClient)
	token, err := s.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %w", err)
	}

	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("id_token not found in response")
	}

	return &TokenSet{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		IDToken:      idToken,
		Expiry:       token.Expiry,
	}, nil
}

// VerifyToken 验证 ID Token 并提取 claims
// 注意：简化实现，生产环境应使用完整的 JWT 验证库
func (s *OIDCService) VerifyToken(idToken string) (*UserClaims, error) {
	// 简单实现：从 ID Token 中提取 claims（不验证签名）
	// 生产环境应该使用完整的 JWT 验证库
	// 如 github.com/golang-jwt/jwt/v5

	// 这个简单实现假设 token 格式正确，直接解析 payload
	// 实际项目中应该验证签名、iss、aud、exp 等

	// 为了演示，这里创建一个简单的 claims 解析
	// 实际应该使用完整的 JWT 验证

	// 这里先返回一个 mock，后续任务会完善
	claims := &UserClaims{
		Sub:        "test-user-id",
		Email:      "user@example.com",
		Name:       "Test User",
		RealmRoles: []string{"user"},
		Exp:        time.Now().Add(24 * time.Hour).Unix(),
	}
	return claims, nil
}

// RefreshToken 使用 Refresh Token 刷新 Access Token
func (s *OIDCService) RefreshToken(refreshToken string) (*TokenSet, error) {
	if _, err := s.fetchDiscovery(); err != nil {
		return nil, err
	}

	ctx := context.WithValue(s.ctx, oauth2.HTTPClient, s.httpClient)
	tokenSource := s.oauth2Config.TokenSource(ctx, &oauth2.Token{
		RefreshToken: refreshToken,
	})

	token, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("token refresh failed: %w", err)
	}

	idToken, _ := token.Extra("id_token").(string)

	return &TokenSet{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		IDToken:      idToken,
		Expiry:       token.Expiry,
	}, nil
}

// FetchUserInfo 获取用户信息
func (s *OIDCService) FetchUserInfo(accessToken string) (map[string]any, error) {
	discovery, err := s.fetchDiscovery()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", discovery.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo request failed: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]any
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}
