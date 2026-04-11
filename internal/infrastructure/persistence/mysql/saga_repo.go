package mysql

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/zouhang1992/ddd_domain/internal/infrastructure/saga"
)

// SagaRepository Saga MySQL 仓储实现
type SagaRepository struct {
	conn *Connection
}

// NewSagaRepository 创建 Saga 仓储
func NewSagaRepository(conn *Connection) *SagaRepository {
	return &SagaRepository{conn: conn}
}

// Save 保存 Saga
func (r *SagaRepository) Save(ctx context.Context, s *saga.Saga) error {
	data, err := json.Marshal(s.Data)
	if err != nil {
		return err
	}

	_, err = r.conn.DB().ExecContext(ctx, `
		INSERT INTO sagas (id, name, state, current_step, error, data, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			name = VALUES(name),
			state = VALUES(state),
			current_step = VALUES(current_step),
			error = VALUES(error),
			data = VALUES(data),
			updated_at = VALUES(updated_at)
	`, s.ID, s.Name, string(s.State), s.CurrentStep, s.Error, data, s.CreatedAt, s.UpdatedAt)

	return err
}

// FindByID 根据 ID 查找 Saga
func (r *SagaRepository) FindByID(ctx context.Context, id string) (*saga.Saga, error) {
	row := r.conn.DB().QueryRowContext(ctx, `
		SELECT id, name, state, current_step, error, data, created_at, updated_at
		FROM sagas WHERE id = ?
	`, id)

	s := &saga.Saga{}
	var dataBytes []byte
	var stateStr string

	err := row.Scan(
		&s.ID,
		&s.Name,
		&stateStr,
		&s.CurrentStep,
		&s.Error,
		&dataBytes,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	s.State = saga.State(stateStr)
	s.Data = make(map[string]any)
	if len(dataBytes) > 0 {
		if err := json.Unmarshal(dataBytes, &s.Data); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// Update 更新 Saga
func (r *SagaRepository) Update(ctx context.Context, s *saga.Saga) error {
	return r.Save(ctx, s)
}
