package facade

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zouhang1992/ddd_domain/internal/application/auth"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/middleware"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/persistence/sqlite"
	"go.uber.org/zap"
)

// OIDCHandler OIDC HTTP 处理器
type OIDCHandler struct {
	oidcService    *auth.OIDCService
	sessionRepo    *sqlite.SessionRepository
	authMiddleware *middleware.AuthMiddleware
	config         auth.Config
	log            *zap.Logger

	// 简单的 state 存储（生产环境应使用 Redis 或数据库）
	stateStore map[string]time.Time
}

// NewOIDCHandler 创建 OIDC 处理器
func NewOIDCHandler(
	oidcService *auth.OIDCService,
	sessionRepo *sqlite.SessionRepository,
	authMiddleware *middleware.AuthMiddleware,
	config auth.Config,
	log *zap.Logger,
) *OIDCHandler {
	return &OIDCHandler{
		oidcService:    oidcService,
		sessionRepo:    sessionRepo,
		authMiddleware: authMiddleware,
		config:         config,
		log:            log,
		stateStore:     make(map[string]time.Time),
	}
}

// RegisterRoutes 注册路由
func (h *OIDCHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /oauth2/login", h.Login)
	mux.HandleFunc("GET /oauth2/callback", h.Callback)
	mux.HandleFunc("POST /oauth2/logout", h.authMiddleware.RequireAuth(h.Logout))
	mux.HandleFunc("GET /oauth2/userinfo", h.authMiddleware.RequireAuth(h.UserInfo))
}

// Login 启动 OIDC 登录流程
func (h *OIDCHandler) Login(w http.ResponseWriter, r *http.Request) {
	// 生成 state
	state, err := auth.GenerateState()
	if err != nil {
		h.log.Error("Failed to generate state", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// 存储 state（5 分钟过期）
	h.stateStore[state] = time.Now().Add(5 * time.Minute)

	// 清理过期的 state
	h.cleanupExpiredStates()

	// 获取认证 URL
	authURL, err := h.oidcService.GetAuthURL(state)
	if err != nil {
		h.log.Error("Failed to get auth URL", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.log.Info("Redirecting to OIDC provider", zap.String("state", state))
	http.Redirect(w, r, authURL, http.StatusFound)
}

// Callback 处理 OIDC 回调
func (h *OIDCHandler) Callback(w http.ResponseWriter, r *http.Request) {
	// 获取 query 参数
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	errorParam := r.URL.Query().Get("error")

	if errorParam != "" {
		h.log.Error("OIDC error", zap.String("error", errorParam))
		http.Error(w, "authentication failed", http.StatusUnauthorized)
		return
	}

	if code == "" || state == "" {
		http.Error(w, "missing code or state", http.StatusBadRequest)
		return
	}

	// 验证 state
	stateExpiry, ok := h.stateStore[state]
	if !ok {
		h.log.Warn("Invalid state", zap.String("state", state))
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}
	if time.Now().After(stateExpiry) {
		h.log.Warn("State expired", zap.String("state", state))
		delete(h.stateStore, state)
		http.Error(w, "state expired", http.StatusBadRequest)
		return
	}

	// 删除已使用的 state
	delete(h.stateStore, state)

	// 用 code 换取 tokens
	tokenSet, err := h.oidcService.ExchangeCode(code)
	if err != nil {
		h.log.Error("Failed to exchange code", zap.Error(err))
		http.Error(w, "authentication failed", http.StatusUnauthorized)
		return
	}

	// 验证 token 并提取 claims
	claims, err := h.oidcService.VerifyToken(tokenSet.IDToken)
	if err != nil {
		h.log.Error("Failed to verify token", zap.Error(err))
		http.Error(w, "authentication failed", http.StatusUnauthorized)
		return
	}

	// 序列化 claims
	claimsJSON, err := sqlite.FromClaims(claims)
	if err != nil {
		h.log.Error("Failed to serialize claims", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// 创建 Session
	session := &sqlite.Session{
		ID:          uuid.NewString(),
		UserID:      claims.Sub,
		AccessToken: tokenSet.AccessToken,
		IDToken:     tokenSet.IDToken,
		Claims:      claimsJSON,
		ExpiresAt:   time.Now().Add(h.config.SessionTTL),
	}

	if tokenSet.RefreshToken != "" {
		session.RefreshToken = sql.NullString{
			String: tokenSet.RefreshToken,
			Valid:  true,
		}
	}

	// 保存 Session
	if err := h.sessionRepo.Save(session); err != nil {
		h.log.Error("Failed to save session", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// 设置 Session Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.SessionCookieName,
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // 生产环境应设为 true
		SameSite: http.SameSiteLaxMode,
		Expires:  session.ExpiresAt,
	})

	h.log.Info("User authenticated successfully",
		zap.String("user_id", claims.Sub),
		zap.String("email", claims.Email))

	// 重定向到首页
	http.Redirect(w, r, "/", http.StatusFound)
}

// Logout 登出
func (h *OIDCHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// 获取 Session
	session := middleware.GetSessionFromContext(r.Context())
	if session != nil {
		// 删除 Session
		if err := h.sessionRepo.Delete(session.ID); err != nil {
			h.log.Error("Failed to delete session", zap.Error(err))
		}
	}

	// 清除 Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1, // 立即删除
		Expires:  time.Now().Add(-1 * time.Hour),
	})

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "logged out successfully",
	})
}

// UserInfo 获取当前用户信息
func (h *OIDCHandler) UserInfo(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"sub":         claims.Sub,
		"email":       claims.Email,
		"name":        claims.Name,
		"roles":       claims.RealmRoles,
		"permissions": claims.Permissions,
	})
}

// cleanupExpiredStates 清理过期的 state
func (h *OIDCHandler) cleanupExpiredStates() {
	now := time.Now()
	for state, expiry := range h.stateStore {
		if now.After(expiry) {
			delete(h.stateStore, state)
		}
	}
	// 如果存储太大，清理一半
	if len(h.stateStore) > 1000 {
		half := len(h.stateStore) / 2
		i := 0
		for state := range h.stateStore {
			if i >= half {
				break
			}
			delete(h.stateStore, state)
			i++
		}
	}
}
