package sqlite

import (
	"database/sql"
	"time"

	depositmodel "github.com/zouhang1992/ddd_domain/internal/domain/deposit/model"
	depositrepo "github.com/zouhang1992/ddd_domain/internal/domain/deposit/repository"
)

// DepositRepository SQLite 押金仓储实现
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
	ReturnedAt *time.Time
	Note       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Save 保存押金
func (r *DepositRepository) Save(deposit *depositmodel.Deposit) error {
	var returnedAt interface{}
	if deposit.ReturnedAt != nil {
		returnedAt = *deposit.ReturnedAt
	}

	_, err := r.conn.DB().Exec(`
		INSERT OR REPLACE INTO deposits (
			id, lease_id, amount, status, returned_at, note,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`,
		deposit.ID(), deposit.LeaseID, deposit.Amount, string(deposit.Status),
		returnedAt, deposit.Note, deposit.CreatedAt, deposit.UpdatedAt)
	return err
}

// FindByID 根据ID查找押金
func (r *DepositRepository) FindByID(id string) (*depositmodel.Deposit, error) {
	row := r.conn.DB().QueryRow(`
		SELECT id, lease_id, amount, status, returned_at, note,
			created_at, updated_at
		FROM deposits WHERE id = ?
		`, id)

	var temp tempDeposit
	err := row.Scan(
		&temp.ID, &temp.LeaseID, &temp.Amount, &temp.Status,
		&temp.ReturnedAt, &temp.Note, &temp.CreatedAt, &temp.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	deposit := depositmodel.NewDeposit(temp.ID, temp.LeaseID, temp.Amount, temp.Note)
	deposit.Status = depositmodel.DepositStatus(temp.Status)
	deposit.ReturnedAt = temp.ReturnedAt
	deposit.CreatedAt = temp.CreatedAt
	deposit.UpdatedAt = temp.UpdatedAt

	return deposit, nil
}

// FindByLeaseID 根据租约ID查找押金
func (r *DepositRepository) FindByLeaseID(leaseID string) (*depositmodel.Deposit, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, lease_id, amount, status, returned_at, note,
			created_at, updated_at
		FROM deposits WHERE lease_id = ? ORDER BY created_at DESC LIMIT 1
		`, leaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var temp tempDeposit
		err := rows.Scan(
			&temp.ID, &temp.LeaseID, &temp.Amount, &temp.Status,
			&temp.ReturnedAt, &temp.Note, &temp.CreatedAt, &temp.UpdatedAt)
		if err != nil {
			return nil, err
		}

		deposit := depositmodel.NewDeposit(temp.ID, temp.LeaseID, temp.Amount, temp.Note)
		deposit.Status = depositmodel.DepositStatus(temp.Status)
		deposit.ReturnedAt = temp.ReturnedAt
		deposit.CreatedAt = temp.CreatedAt
		deposit.UpdatedAt = temp.UpdatedAt

		return deposit, nil
	}
	return nil, nil
}

// FindAll 查找所有押金
func (r *DepositRepository) FindAll() ([]*depositmodel.Deposit, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, lease_id, amount, status, returned_at, note,
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
			&temp.ReturnedAt, &temp.Note, &temp.CreatedAt, &temp.UpdatedAt)
		if err != nil {
			return nil, err
		}

		deposit := depositmodel.NewDeposit(temp.ID, temp.LeaseID, temp.Amount, temp.Note)
		deposit.Status = depositmodel.DepositStatus(temp.Status)
		deposit.ReturnedAt = temp.ReturnedAt
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
