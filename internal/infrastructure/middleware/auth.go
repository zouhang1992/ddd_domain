package middleware

import (
	"context"
	"net/http"

	"github.com/zouhang1992/ddd_domain/internal/application/auth"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/persistence/sqlite"
	"go.uber.org/zap"
)

type contextKey string

const (
	// UserContextKey 用户信息上下文键
	UserContextKey contextKey = "user"
	// SessionContextKey Session 信息上下文键
	SessionContextKey contextKey = "session"
	// SessionCookieName Session Cookie 名称
	SessionCookieName = "session_id"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	sessionRepo *sqlite.SessionRepository
	oidcService *auth.OIDCService
	log         *zap.Logger
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(
	sessionRepo *sqlite.SessionRepository,
	oidcService *auth.OIDCService,
	log *zap.Logger,
) *AuthMiddleware {
	return &AuthMiddleware{
		sessionRepo: sessionRepo,
		oidcService: oidcService,
		log:         log,
	}
}

// RequireAuth 要求认证的中间件
func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从 Cookie 获取 Session ID
		cookie, err := r.Cookie(SessionCookieName)
		if err != nil {
			m.log.Debug("No session cookie found", zap.Error(err))
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// 查找 Session
		session, err := m.sessionRepo.FindByID(cookie.Value)
		if err != nil {
			m.log.Error("Failed to find session", zap.Error(err))
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if session == nil {
			m.log.Debug("Session not found", zap.String("session_id", cookie.Value))
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// 解析 claims
		claims, err := sqlite.ToClaims(session.Claims)
		if err != nil {
			m.log.Error("Failed to parse claims", zap.Error(err))
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// 将用户信息注入上下文
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		ctx = context.WithValue(ctx, SessionContextKey, session)

		// 继续处理
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// OptionalAuth 可选认证的中间件
func (m *AuthMiddleware) OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 尝试从 Cookie 获取 Session ID
		cookie, err := r.Cookie(SessionCookieName)
		if err == nil && cookie != nil {
			// 查找 Session
			session, err := m.sessionRepo.FindByID(cookie.Value)
			if err == nil && session != nil {
				// 解析 claims
				claims, err := sqlite.ToClaims(session.Claims)
				if err == nil {
					// 将用户信息注入上下文
					ctx := context.WithValue(r.Context(), UserContextKey, claims)
					ctx = context.WithValue(ctx, SessionContextKey, session)
					r = r.WithContext(ctx)
				}
			}
		}

		// 继续处理（即使没有认证）
		next.ServeHTTP(w, r)
	}
}

// GetUserFromContext 从上下文获取用户信息
func GetUserFromContext(ctx context.Context) *auth.UserClaims {
	if claims, ok := ctx.Value(UserContextKey).(*auth.UserClaims); ok {
		return claims
	}
	return nil
}

// GetSessionFromContext 从上下文获取 Session 信息
func GetSessionFromContext(ctx context.Context) *sqlite.Session {
	if session, ok := ctx.Value(SessionContextKey).(*sqlite.Session); ok {
		return session
	}
	return nil
}
