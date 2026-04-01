package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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
	discoveryCache  *oidcDiscovery
	discoveryMutex  sync.RWMutex
	discoveryExpiry time.Time

	// JWKS 缓存
	jwksCache  *jwks
	jwksMutex  sync.RWMutex
	jwksExpiry time.Time
}

// oidcDiscovery OIDC discovery 响应
type oidcDiscovery struct {
	Issuer        string `json:"issuer"`
	AuthURL       string `json:"authorization_endpoint"`
	TokenURL      string `json:"token_endpoint"`
	JWKSURL       string `json:"jwks_uri"`
	UserInfoURL   string `json:"userinfo_endpoint"`
	EndSessionURL string `json:"end_session_endpoint,omitempty"`
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
func (s *OIDCService) VerifyToken(idToken string) (*UserClaims, error) {
	// 分割 JWT（格式：header.payload.signature）
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	// 解码 payload（第二部分）
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode token: %w", err)
	}

	// 解析 JSON 到 UserClaims
	var claims UserClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	// 解析 resource_access 获取角色（Keycloak 的角色结构）
	var rawClaims map[string]any
	if err := json.Unmarshal(payload, &rawClaims); err == nil {
		// 尝试从 realm_access.roles 获取角色
		if realmAccess, ok := rawClaims["realm_access"].(map[string]any); ok {
			if roles, ok := realmAccess["roles"].([]any); ok {
				claims.RealmRoles = make([]string, 0, len(roles))
				for _, r := range roles {
					if roleStr, ok := r.(string); ok {
						claims.RealmRoles = append(claims.RealmRoles, roleStr)
					}
				}
			}
		}
		// 尝试从 resource_access 获取角色
		if resourceAccess, ok := rawClaims["resource_access"].(map[string]any); ok {
			claims.ResourceRoles = make(map[string][]string)
			for clientId, access := range resourceAccess {
				if accessMap, ok := access.(map[string]any); ok {
					if roles, ok := accessMap["roles"].([]any); ok {
						clientRoles := make([]string, 0, len(roles))
						for _, r := range roles {
							if roleStr, ok := r.(string); ok {
								clientRoles = append(clientRoles, roleStr)
							}
						}
						claims.ResourceRoles[clientId] = clientRoles
					}
				}
			}
		}
	}

	return &claims, nil
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

// GetEndSessionURL 获取 OIDC 单点登出 URL
func (s *OIDCService) GetEndSessionURL(idTokenHint string, postLogoutRedirectURI string) (string, error) {
	discovery, err := s.fetchDiscovery()
	if err != nil {
		return "", err
	}

	if discovery.EndSessionURL == "" {
		return "", fmt.Errorf("end_session_endpoint not available")
	}

	endSessionURL := discovery.EndSessionURL
	params := []string{}

	if idTokenHint != "" {
		params = append(params, fmt.Sprintf("id_token_hint=%s", idTokenHint))
	}
	if postLogoutRedirectURI != "" {
		params = append(params, fmt.Sprintf("post_logout_redirect_uri=%s", postLogoutRedirectURI))
	}

	if len(params) > 0 {
		sep := "?"
		if len(discovery.EndSessionURL) > 0 && discovery.EndSessionURL[len(discovery.EndSessionURL)-1] == '?' {
			sep = ""
		} else if len(discovery.EndSessionURL) > 0 && discovery.EndSessionURL[len(discovery.EndSessionURL)-1] != '&' {
			sep = "?"
			if strings.Contains(discovery.EndSessionURL, "?") {
				sep = "&"
			}
		}
		endSessionURL += sep + strings.Join(params, "&")
	}

	return endSessionURL, nil
}
