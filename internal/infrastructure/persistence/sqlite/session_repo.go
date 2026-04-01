package sqlite

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Session represents a user session
type Session struct {
	ID        string
	UserID    string
	Claims    UserClaims
	ExpiresAt time.Time
	CreatedAt time.Time
}

// UserClaims represents the user claims stored in the session
type UserClaims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

// ToClaims converts UserClaims to JSON string
func (c *UserClaims) ToClaims() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromClaims converts JSON string to UserClaims
func FromClaims(data string) (*UserClaims, error) {
	var claims UserClaims
	if err := json.Unmarshal([]byte(data), &claims); err != nil {
		return nil, err
	}
	return &claims, nil
}

// SessionRepository SQLite session repository implementation
type SessionRepository struct {
	conn *Connection
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(conn *Connection) *SessionRepository {
	return &SessionRepository{conn: conn}
}

// Save saves a session
func (r *SessionRepository) Save(session *Session) error {
	if session.ID == "" {
		session.ID = uuid.NewString()
	}

	claimsJSON, err := session.Claims.ToClaims()
	if err != nil {
		return err
	}

	_, err = r.conn.DB().Exec(`
		INSERT OR REPLACE INTO sessions (
			id, user_id, claims, expires_at, created_at
		) VALUES (?, ?, ?, ?, ?)
		`,
		session.ID, session.UserID, claimsJSON, session.ExpiresAt, session.CreatedAt)
	return err
}

// FindByID finds a session by ID
func (r *SessionRepository) FindByID(id string) (*Session, error) {
	row := r.conn.DB().QueryRow(`
		SELECT id, user_id, claims, expires_at, created_at
		FROM sessions WHERE id = ?
		`, id)

	var temp struct {
		ID        string
		UserID    string
		Claims    string
		ExpiresAt time.Time
		CreatedAt time.Time
	}

	err := row.Scan(&temp.ID, &temp.UserID, &temp.Claims, &temp.ExpiresAt, &temp.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	claims, err := FromClaims(temp.Claims)
	if err != nil {
		return nil, err
	}

	return &Session{
		ID:        temp.ID,
		UserID:    temp.UserID,
		Claims:    *claims,
		ExpiresAt: temp.ExpiresAt,
		CreatedAt: temp.CreatedAt,
	}, nil
}

// Delete deletes a session by ID
func (r *SessionRepository) Delete(id string) error {
	_, err := r.conn.DB().Exec("DELETE FROM sessions WHERE id = ?", id)
	return err
}

// DeleteExpired deletes all expired sessions
func (r *SessionRepository) DeleteExpired(now time.Time) (int64, error) {
	result, err := r.conn.DB().Exec("DELETE FROM sessions WHERE expires_at < ?", now)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
