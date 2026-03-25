package saga

import (
	"context"
	"fmt"
	"time"
)

// Orchestrator Saga 编排器
type Orchestrator struct {
	repo Repository
}

// NewOrchestrator 创建 Saga 编排器
func NewOrchestrator(repo Repository) *Orchestrator {
	return &Orchestrator{repo: repo}
}

// Execute 执行 Saga
func (o *Orchestrator) Execute(ctx context.Context, saga *Saga) error {
	saga.State = StateRunning
	saga.UpdatedAt = now()
	if err := o.repo.Update(ctx, saga); err != nil {
		return err
	}

	for i := saga.CurrentStep; i < len(saga.Steps); i++ {
		step := saga.Steps[i]
		saga.CurrentStep = i

		if err := step.Execute(ctx, saga.Data); err != nil {
			saga.State = StateFailed
			saga.Error = err.Error()
			saga.UpdatedAt = now()
			_ = o.repo.Update(ctx, saga)

			return o.compensate(ctx, saga)
		}

		saga.UpdatedAt = now()
		if err := o.repo.Update(ctx, saga); err != nil {
			return err
		}
	}

	saga.State = StateCompleted
	saga.UpdatedAt = now()
	return o.repo.Update(ctx, saga)
}

// compensate 执行补偿
func (o *Orchestrator) compensate(ctx context.Context, saga *Saga) error {
	saga.State = StateCompensating
	saga.UpdatedAt = now()
	if err := o.repo.Update(ctx, saga); err != nil {
		return err
	}

	for i := saga.CurrentStep - 1; i >= 0; i-- {
		step := saga.Steps[i]
		if step.Compensate != nil {
			if err := step.Compensate(ctx, saga.Data); err != nil {
				return fmt.Errorf("compensation failed at step %d: %w", i, err)
			}
		}
	}

	saga.State = StateRolledBack
	saga.UpdatedAt = now()
	return o.repo.Update(ctx, saga)
}

func now() time.Time {
	return time.Now()
}
