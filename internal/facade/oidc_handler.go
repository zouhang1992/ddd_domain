package facade

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zouhang1992/ddd_domain/internal/application/auth"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/middleware"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/persistence/mysql"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/persistence/sqlite"
	"go.uber.org/zap"
)

// OIDCHandler OIDC HTTP 处理器
type OIDCHandler struct {
	oidcService    *auth.OIDCService
	sessionRepo    any
	authMiddleware *middleware.AuthMiddleware
	config         auth.Config
	log            *zap.Logger

	// state 存储（生产环境应使用 Redis 或数据库）
	stateStore map[string]stateData
}

type stateData struct {
	expiry    time.Time
	returnURL string
}

// NewOIDCHandler 创建 OIDC 处理器
func NewOIDCHandler(
	oidcService *auth.OIDCService,
	sessionRepo any,
	authMiddleware *middleware.AuthMiddleware,
	config auth.Config,
	log *zap.Logger,
) *OIDCHandler {
	log.Info("OIDC Handler initialized",
		zap.Bool("dev_mode", config.DevMode),
		zap.String("issuer_url", config.IssuerURL),
		zap.String("client_id", config.ClientID),
		zap.String("redirect_url", config.RedirectURL),
	)
	return &OIDCHandler{
		oidcService:    oidcService,
		sessionRepo:    sessionRepo,
		authMiddleware: authMiddleware,
		config:         config,
		log:            log,
		stateStore:     make(map[string]stateData),
	}
}

// RegisterRoutes 注册路由
func (h *OIDCHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /oauth2/login", h.Login)
	mux.HandleFunc("GET /oauth2/callback", h.Callback)

	// 开发模式下，logout 和 userinfo 不需要认证
	if h.config.DevMode {
		mux.HandleFunc("POST /oauth2/logout", h.Logout)
		mux.HandleFunc("GET /oauth2/logout", h.Logout)
		mux.HandleFunc("GET /oauth2/userinfo", h.UserInfo)
	} else {
		mux.HandleFunc("POST /oauth2/logout", h.authMiddleware.RequireAuth(h.Logout))
		mux.HandleFunc("GET /oauth2/logout", h.authMiddleware.RequireAuth(h.Logout))
		mux.HandleFunc("GET /oauth2/userinfo", h.authMiddleware.RequireAuth(h.UserInfo))
	}
}

// Login 启动 OIDC 登录流程
func (h *OIDCHandler) Login(w http.ResponseWriter, r *http.Request) {
	// 获取 return_url
	returnURL := r.URL.Query().Get("return_url")
	if returnURL == "" {
		returnURL = "/"
	}

	// 开发模式：直接创建 mock session，跳过 OIDC
	if h.config.DevMode {
		h.devModeLogin(w, r, returnURL)
		return
	}

	// 生成 state
	state, err := auth.GenerateState()
	if err != nil {
		h.log.Error("Failed to generate state", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// 存储 state（5 分钟过期）
	h.stateStore[state] = stateData{
		expiry:    time.Now().Add(5 * time.Minute),
		returnURL: returnURL,
	}

	// 清理过期的 state
	h.cleanupExpiredStates()

	// 获取认证 URL
	authURL, err := h.oidcService.GetAuthURL(state)
	if err != nil {
		h.log.Error("Failed to get auth URL", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.log.Info("Redirecting to OIDC provider", zap.String("state", state), zap.String("return_url", returnURL))
	http.Redirect(w, r, authURL, http.StatusFound)
}

// devModeLogin 开发模式登录（直接创建 mock session）
func (h *OIDCHandler) devModeLogin(w http.ResponseWriter, r *http.Request, returnURL string) {
	// 创建 mock claims
	claims := &auth.UserClaims{
		Sub:         "dev-user-id",
		Email:       "dev@example.com",
		Name:        "Developer",
		RealmRoles:  []string{"user", "admin"},
		Permissions: []string{"read", "write", "delete"},
		Exp:         time.Now().Add(24 * time.Hour).Unix(),
	}

	h.saveSessionAndRedirect(w, r, claims, nil, returnURL)
}

// saveSessionAndRedirect saves the session using the appropriate repository and redirects
func (h *OIDCHandler) saveSessionAndRedirect(w http.ResponseWriter, r *http.Request, claims *auth.UserClaims, tokenSet *auth.TokenSet, returnURL string) {
	sessionID := uuid.NewString()

	switch repo := h.sessionRepo.(type) {
	case *sqlite.SessionRepository:
		// 序列化 claims
		claimsJSON, err := sqlite.FromClaims(claims)
		if err != nil {
			h.log.Error("Failed to serialize claims", zap.Error(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		// 创建 Session
		session := &sqlite.Session{
			ID:        sessionID,
			UserID:    claims.Sub,
			Claims:    claimsJSON,
			ExpiresAt: time.Now().Add(h.config.SessionTTL),
		}

		if tokenSet != nil {
			session.AccessToken = tokenSet.AccessToken
			session.IDToken = tokenSet.IDToken
			if tokenSet.RefreshToken != "" {
				session.RefreshToken = sql.NullString{
					String: tokenSet.RefreshToken,
					Valid:  true,
				}
			}
		}

		// 保存 Session
		if err := repo.Save(session); err != nil {
			h.log.Error("Failed to save session", zap.Error(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

	case *mysql.SessionRepository:
		// 序列化 claims
		claimsJSON, err := mysql.FromClaims(claims)
		if err != nil {
			h.log.Error("Failed to serialize claims", zap.Error(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		// 创建 Session
		session := &mysql.Session{
			ID:        sessionID,
			UserID:    claims.Sub,
			Claims:    claimsJSON,
			ExpiresAt: time.Now().Add(h.config.SessionTTL),
		}

		if tokenSet != nil {
			session.AccessToken = tokenSet.AccessToken
			session.IDToken = tokenSet.IDToken
			if tokenSet.RefreshToken != "" {
				session.RefreshToken = sql.NullString{
					String: tokenSet.RefreshToken,
					Valid:  true,
				}
			}
		}

		// 保存 Session
		if err := repo.Save(session); err != nil {
			h.log.Error("Failed to save session", zap.Error(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	default:
		h.log.Error("Unsupported session repository type")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// 设置 Session Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.SessionCookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(h.config.SessionTTL),
	})

	h.log.Info("Login successful",
		zap.String("user_id", claims.Sub),
		zap.String("email", claims.Email),
		zap.String("return_url", returnURL))

	// 重定向到 returnURL 或前端首页
	redirectURL := h.config.FrontendURL
	if returnURL != "/" {
		redirectURL += returnURL
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
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
	stateData, ok := h.stateStore[state]
	if !ok {
		h.log.Warn("Invalid state", zap.String("state", state))
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}
	if time.Now().After(stateData.expiry) {
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

	h.saveSessionAndRedirect(w, r, claims, tokenSet, stateData.returnURL)
}

// Logout 登出
func (h *OIDCHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// 开发模式：记录登出日志
	if h.config.DevMode {
		h.log.Info("Dev mode logout requested")
	}

	var endSessionURL string
	var idToken string

	// 获取 Session
	sessionAny := middleware.GetSessionFromContext(r.Context())
	if sessionAny != nil {
		var sessionID string

		switch s := sessionAny.(type) {
		case *sqlite.Session:
			idToken = s.IDToken
			sessionID = s.ID
		case *mysql.Session:
			idToken = s.IDToken
			sessionID = s.ID
		}

		h.log.Info("Session found for logout",
			zap.String("session_id", sessionID),
			zap.Bool("has_id_token", idToken != ""))

		// 删除 Session
		switch repo := h.sessionRepo.(type) {
		case *sqlite.SessionRepository:
			if err := repo.Delete(sessionID); err != nil {
				h.log.Error("Failed to delete session", zap.Error(err))
			} else {
				h.log.Info("Session deleted", zap.String("session_id", sessionID))
			}
		case *mysql.SessionRepository:
			if err := repo.Delete(sessionID); err != nil {
				h.log.Error("Failed to delete session", zap.Error(err))
			} else {
				h.log.Info("Session deleted", zap.String("session_id", sessionID))
			}
		}
	} else {
		h.log.Warn("No session found for logout")
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

	// 获取 OIDC 单点登出 URL
	h.log.Info("Preparing end session URL",
		zap.Bool("dev_mode", h.config.DevMode),
		zap.Bool("has_id_token", idToken != ""))
	if !h.config.DevMode && idToken != "" {
		postLogoutRedirectURI := h.config.FrontendURL
		if url, err := h.oidcService.GetEndSessionURL(idToken, postLogoutRedirectURI); err == nil {
			endSessionURL = url
			h.log.Info("Got end session URL, redirecting", zap.String("url", endSessionURL))
			http.Redirect(w, r, endSessionURL, http.StatusFound)
			return
		} else {
			h.log.Warn("Failed to get end session URL", zap.Error(err))
		}
	} else {
		h.log.Info("Skipping OIDC end session",
			zap.Bool("dev_mode", h.config.DevMode),
			zap.Bool("has_id_token", idToken != ""))
	}

	// 如果没有 OIDC 单点登出，直接重定向到前端
	h.log.Info("Redirecting to frontend home")
	http.Redirect(w, r, h.config.FrontendURL, http.StatusFound)
}

// UserInfo 获取当前用户信息
func (h *OIDCHandler) UserInfo(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("UserInfo request received", zap.Bool("dev_mode", h.config.DevMode))

	// 开发模式：如果没有 session，直接返回 mock 用户
	if h.config.DevMode {
		h.log.Debug("Returning mock user (dev mode)")
		session := middleware.GetSessionFromContext(r.Context())
		if session == nil {
			// 创建临时 mock 用户
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"sub":         "dev-user-id",
				"email":       "dev@example.com",
				"name":        "Developer",
				"roles":       []string{"user", "admin"},
				"permissions": []string{"read", "write", "delete"},
			})
			return
		}
	}

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
		"hello":       claims.Hello,
	})
}

// cleanupExpiredStates 清理过期的 state
func (h *OIDCHandler) cleanupExpiredStates() {
	now := time.Now()
	for state, data := range h.stateStore {
		if now.After(data.expiry) {
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
