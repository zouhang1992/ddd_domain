package sqlite

import (
	"database/sql"
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
)

// DepositRepository SQLite 押金仓储实现
type DepositRepository struct {
	conn *Connection
}

// NewDepositRepository 创建押金仓储
func NewDepositRepository(conn *Connection) repository.DepositRepository {
	return &DepositRepository{conn: conn}
}

// Save 保存押金
func (r *DepositRepository) Save(deposit *model.Deposit) error {
	var refundedAt, deductedAt interface{}
	if deposit.RefundedAt != nil {
		refundedAt = *deposit.RefundedAt
	}
	if deposit.DeductedAt != nil {
		deductedAt = *deposit.DeductedAt
	}

	_, err := r.conn.DB().Exec(`
		INSERT OR REPLACE INTO deposits (
			id, lease_id, amount, status, refunded_at, deducted_at, note,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		deposit.ID, deposit.LeaseID, deposit.Amount, string(deposit.Status),
		refundedAt, deductedAt, deposit.Note, deposit.CreatedAt, deposit.UpdatedAt)
	return err
}

// FindByID 根据ID查找押金
func (r *DepositRepository) FindByID(id string) (*model.Deposit, error) {
	row := r.conn.DB().QueryRow(`
		SELECT id, lease_id, amount, status, refunded_at, deducted_at, note,
			created_at, updated_at
		FROM deposits WHERE id = ?
	`, id)

	deposit := &model.Deposit{}
	var statusStr string
	var refundedAt, deductedAt interface{}
	err := row.Scan(
		&deposit.ID, &deposit.LeaseID, &deposit.Amount, &statusStr,
		&refundedAt, &deductedAt, &deposit.Note, &deposit.CreatedAt, &deposit.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	deposit.Status = model.DepositStatus(statusStr)
	if refundedAt != nil {
		if t, ok := refundedAt.(time.Time); ok {
			deposit.RefundedAt = &t
		}
	}
	if deductedAt != nil {
		if t, ok := deductedAt.(time.Time); ok {
			deposit.DeductedAt = &t
		}
	}

	return deposit, nil
}

// FindByLeaseID 根据租约ID查找押金
func (r *DepositRepository) FindByLeaseID(leaseID string) (*model.Deposit, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, lease_id, amount, status, refunded_at, deducted_at, note,
			created_at, updated_at
		FROM deposits WHERE lease_id = ? ORDER BY created_at DESC LIMIT 1
	`, leaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		deposit := &model.Deposit{}
		var statusStr string
		var refundedAt, deductedAt interface{}
		err := rows.Scan(
			&deposit.ID, &deposit.LeaseID, &deposit.Amount, &statusStr,
			&refundedAt, &deductedAt, &deposit.Note, &deposit.CreatedAt, &deposit.UpdatedAt)
		if err != nil {
			return nil, err
		}

		deposit.Status = model.DepositStatus(statusStr)
		if refundedAt != nil {
			if t, ok := refundedAt.(time.Time); ok {
				deposit.RefundedAt = &t
			}
		}
		if deductedAt != nil {
			if t, ok := deductedAt.(time.Time); ok {
				deposit.DeductedAt = &t
			}
		}

		return deposit, nil
	}
	return nil, nil
}

// FindAll 查找所有押金
func (r *DepositRepository) FindAll() ([]*model.Deposit, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, lease_id, amount, status, refunded_at, deducted_at, note,
			created_at, updated_at
		FROM deposits ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deposits []*model.Deposit
	for rows.Next() {
		deposit := &model.Deposit{}
		var statusStr string
		var refundedAt, deductedAt interface{}
		err := rows.Scan(
			&deposit.ID, &deposit.LeaseID, &deposit.Amount, &statusStr,
			&refundedAt, &deductedAt, &deposit.Note, &deposit.CreatedAt, &deposit.UpdatedAt)
		if err != nil {
			return nil, err
		}

		deposit.Status = model.DepositStatus(statusStr)
		if refundedAt != nil {
			if t, ok := refundedAt.(time.Time); ok {
				deposit.RefundedAt = &t
			}
		}
		if deductedAt != nil {
			if t, ok := deductedAt.(time.Time); ok {
				deposit.DeductedAt = &t
			}
		}

		deposits = append(deposits, deposit)
	}
	return deposits, nil
}

// Delete 删除押金
func (r *DepositRepository) Delete(id string) error {
	_, err := r.conn.DB().Exec("DELETE FROM deposits WHERE id = ?", id)
	return err
}
