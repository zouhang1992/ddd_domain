package sqlite

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Session Session 数据模型
type Session struct {
	ID           string
	UserID       string
	AccessToken  string
	RefreshToken sql.NullString
	IDToken      string
	Claims       []byte
	ExpiresAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// SessionRepository Session 仓储实现
type SessionRepository struct {
	conn *Connection
}

// NewSessionRepository 创建 Session 仓储
func NewSessionRepository(conn *Connection) *SessionRepository {
	return &SessionRepository{conn: conn}
}

// Save 保存 Session
func (r *SessionRepository) Save(session *Session) error {
	if session.ID == "" {
		session.ID = uuid.NewString()
	}
	now := time.Now()
	if session.CreatedAt.IsZero() {
		session.CreatedAt = now
	}
	session.UpdatedAt = now

	_, err := r.conn.DB().Exec(`
		INSERT OR REPLACE INTO sessions (
			id, user_id, access_token, refresh_token, id_token,
			claims, expires_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
		session.ID, session.UserID, session.AccessToken,
		session.RefreshToken, session.IDToken, session.Claims,
		session.ExpiresAt, session.CreatedAt, session.UpdatedAt)
	return err
}

// FindByID 根据 ID 查找 Session
func (r *SessionRepository) FindByID(id string) (*Session, error) {
	row := r.conn.DB().QueryRow(`
		SELECT id, user_id, access_token, refresh_token, id_token,
			claims, expires_at, created_at, updated_at
		FROM sessions WHERE id = ?
		`, id)

	var session Session
	var refreshToken sql.NullString
	err := row.Scan(
		&session.ID, &session.UserID, &session.AccessToken,
		&refreshToken, &session.IDToken, &session.Claims,
		&session.ExpiresAt, &session.CreatedAt, &session.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	session.RefreshToken = refreshToken
	return &session, nil
}

// Delete 删除 Session
func (r *SessionRepository) Delete(id string) error {
	_, err := r.conn.DB().Exec("DELETE FROM sessions WHERE id = ?", id)
	return err
}

// DeleteExpired 删除过期的 Session
func (r *SessionRepository) DeleteExpired() (int64, error) {
	result, err := r.conn.DB().Exec("DELETE FROM sessions WHERE expires_at < ?", time.Now())
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// UserClaims 用户 Claims
type UserClaims struct {
	Sub          string                 `json:"sub"`
	Email        string                 `json:"email"`
	Name         string                 `json:"name"`
	RealmRoles   []string               `json:"realm_roles,omitempty"`
	ResourceRoles map[string][]string    `json:"resource_roles,omitempty"`
	Permissions  []string               `json:"permissions,omitempty"`
	Exp          int64                  `json:"exp"`
	Extra        map[string]any `json:"-"`
}

// ToClaims 将 JSON 转换为 UserClaims
func ToClaims(data []byte) (*UserClaims, error) {
	var claims UserClaims
	if err := json.Unmarshal(data, &claims); err != nil {
		return nil, err
	}
	return &claims, nil
}

// FromClaims 将 UserClaims 转换为 JSON
func FromClaims(claims *UserClaims) ([]byte, error) {
	return json.Marshal(claims)
}
