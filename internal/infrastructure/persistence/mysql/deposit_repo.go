package mysql

import (
	"database/sql"
	"time"

	depositmodel "github.com/zouhang1992/ddd_domain/internal/domain/deposit/model"
	depositrepo "github.com/zouhang1992/ddd_domain/internal/domain/deposit/repository"
)

// DepositRepository MySQL 押金仓储实现
type DepositRepository struct {
	conn *Connection
}

// NewDepositRepository 创建押金仓储
func NewDepositRepository(conn *Connection) depositrepo.DepositRepository {
	return &DepositRepository{conn: conn}
}

// tempDeposit is a temporary struct for scanning
type tempDeposit struct {
	ID         string
	LeaseID    string
	Amount     int64
	Status     string
	RefundedAt *time.Time
	DeductedAt *time.Time
	Note       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Save 保存押金
func (r *DepositRepository) Save(deposit *depositmodel.Deposit) error {
	var refundedAt interface{}
	if deposit.RefundedAt != nil {
		refundedAt = *deposit.RefundedAt
	}
	var deductedAt interface{}
	if deposit.DeductedAt != nil {
		deductedAt = *deposit.DeductedAt
	}

	_, err := r.conn.DB().Exec(`
		INSERT INTO deposits (
			id, lease_id, amount, status, refunded_at, deducted_at, note,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			lease_id = VALUES(lease_id),
			amount = VALUES(amount),
			status = VALUES(status),
			refunded_at = VALUES(refunded_at),
			deducted_at = VALUES(deducted_at),
			note = VALUES(note),
			updated_at = VALUES(updated_at)
		`,
		deposit.IDField, deposit.LeaseID, deposit.Amount, string(deposit.Status),
		refundedAt, deductedAt, deposit.Note, deposit.CreatedAt, deposit.UpdatedAt)
	return err
}

// FindByID 根据ID查找押金
func (r *DepositRepository) FindByID(id string) (*depositmodel.Deposit, error) {
	row := r.conn.DB().QueryRow(`
		SELECT id, lease_id, amount, status, refunded_at, deducted_at, note,
			created_at, updated_at
		FROM deposits WHERE id = ?
		`, id)

	var temp tempDeposit
	err := row.Scan(
		&temp.ID, &temp.LeaseID, &temp.Amount, &temp.Status,
		&temp.RefundedAt, &temp.DeductedAt, &temp.Note, &temp.CreatedAt, &temp.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	deposit := depositmodel.NewDeposit(temp.ID, temp.LeaseID, temp.Amount, temp.Note)
	deposit.Status = depositmodel.DepositStatus(temp.Status)
	deposit.RefundedAt = temp.RefundedAt
	deposit.DeductedAt = temp.DeductedAt
	deposit.CreatedAt = temp.CreatedAt
	deposit.UpdatedAt = temp.UpdatedAt
	deposit.ClearEvents()

	return deposit, nil
}

// FindByLeaseID 根据租约ID查找押金
func (r *DepositRepository) FindByLeaseID(leaseID string) (*depositmodel.Deposit, error) {
	row := r.conn.DB().QueryRow(`
		SELECT id, lease_id, amount, status, refunded_at, deducted_at, note,
			created_at, updated_at
		FROM deposits WHERE lease_id = ? ORDER BY created_at DESC LIMIT 1
		`, leaseID)

	var temp tempDeposit
	err := row.Scan(
		&temp.ID, &temp.LeaseID, &temp.Amount, &temp.Status,
		&temp.RefundedAt, &temp.DeductedAt, &temp.Note, &temp.CreatedAt, &temp.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	deposit := depositmodel.NewDeposit(temp.ID, temp.LeaseID, temp.Amount, temp.Note)
	deposit.Status = depositmodel.DepositStatus(temp.Status)
	deposit.RefundedAt = temp.RefundedAt
	deposit.DeductedAt = temp.DeductedAt
	deposit.CreatedAt = temp.CreatedAt
	deposit.UpdatedAt = temp.UpdatedAt

	return deposit, nil
}

// FindAll 查找所有押金
func (r *DepositRepository) FindAll() ([]*depositmodel.Deposit, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, lease_id, amount, status, refunded_at, deducted_at, note,
			created_at, updated_at
		FROM deposits ORDER BY created_at DESC
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deposits []*depositmodel.Deposit
	for rows.Next() {
		var temp tempDeposit
		err := rows.Scan(
			&temp.ID, &temp.LeaseID, &temp.Amount, &temp.Status,
			&temp.RefundedAt, &temp.DeductedAt, &temp.Note, &temp.CreatedAt, &temp.UpdatedAt)
		if err != nil {
			return nil, err
		}

		deposit := depositmodel.NewDeposit(temp.ID, temp.LeaseID, temp.Amount, temp.Note)
		deposit.Status = depositmodel.DepositStatus(temp.Status)
		deposit.RefundedAt = temp.RefundedAt
		deposit.DeductedAt = temp.DeductedAt
		deposit.CreatedAt = temp.CreatedAt
		deposit.UpdatedAt = temp.UpdatedAt

		deposits = append(deposits, deposit)
	}
	return deposits, nil
}

// Delete 删除押金
func (r *DepositRepository) Delete(id string) error {
	_, err := r.conn.DB().Exec("DELETE FROM deposits WHERE id = ?", id)
	return err
}
