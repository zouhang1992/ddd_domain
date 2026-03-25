package saga

import (
	"context"
	"time"
)

// State Saga 状态
type State string

const (
	StatePending      State = "pending"
	StateRunning      State = "running"
	StateCompleted    State = "completed"
	StateFailed       State = "failed"
	StateCompensating State = "compensating"
	StateRolledBack   State = "rolled_back"
)

// Step Saga 步骤
type Step struct {
	Name       string
	Execute    func(ctx context.Context, data map[string]any) error
	Compensate func(ctx context.Context, data map[string]any) error
}

// Saga Saga 定义
type Saga struct {
	ID          string
	Name        string
	State       State
	Steps       []Step
	Data        map[string]any
	CurrentStep int
	Error       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewSaga 创建新 Saga
func NewSaga(id, name string, steps []Step) *Saga {
	return &Saga{
		ID:        id,
		Name:      name,
		State:     StatePending,
		Steps:     steps,
		Data:      make(map[string]any),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Repository Saga 仓储接口
type Repository interface {
	Save(ctx context.Context, saga *Saga) error
	FindByID(ctx context.Context, id string) (*Saga, error)
	Update(ctx context.Context, saga *Saga) error
}
