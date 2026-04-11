package persistence

import (
	"time"

	"github.com/zouhang1992/ddd_domain/internal/application/auth"
)

// Session 通用Session接口
type Session interface {
	GetID() string
	GetUserID() string
	GetClaims() []byte
	GetExpiresAt() time.Time
}

// SessionRepository 通用Session仓储接口
type SessionRepository interface {
	Save(session Session) error
	FindByID(id string) (Session, error)
	Delete(id string) error
	DeleteExpired() (int64, error)
}

// ToClaims 通用函数类型
type ToClaimsFunc func(data []byte) (*auth.UserClaims, error)

// FromClaims 通用函数类型
type FromClaimsFunc func(claims *auth.UserClaims) ([]byte, error)
