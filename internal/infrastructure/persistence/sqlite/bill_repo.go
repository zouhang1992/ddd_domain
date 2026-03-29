package sqlite

import (
	"database/sql"
	"time"

	billmodel "github.com/zouhang1992/ddd_domain/internal/domain/bill/model"
	billrepo "github.com/zouhang1992/ddd_domain/internal/domain/bill/repository"
)

// BillRepository SQLite 账单仓储实现
type BillRepository struct {
	conn *Connection
}

// NewBillRepository 创建账单仓储
func NewBillRepository(conn *Connection) billrepo.BillRepository {
	return &BillRepository{conn: conn}
}

// Save 保存账单
func (r *BillRepository) Save(bill *billmodel.Bill) error {
	var paidAt interface{}
	if bill.PaidAt != nil {
		paidAt = *bill.PaidAt
	}

	_, err := r.conn.DB().Exec(`
		INSERT OR REPLACE INTO bills (
			id, lease_id, type, amount, due_date, paid_at, note, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
		bill.ID(), bill.LeaseID, string(bill.Type), bill.Amount, bill.DueDate,
		paidAt, bill.Note, bill.CreatedAt, bill.UpdatedAt)
	return err
}

// tempBill is a temporary struct for scanning
type tempBill struct {
	ID        string
	LeaseID   string
	Type      string
	Amount    int64
	DueDate   time.Time
	PaidAt    *time.Time
	Note      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// FindByID 根据ID查找账单
func (r *BillRepository) FindByID(id string) (*billmodel.Bill, error) {
	row := r.conn.DB().QueryRow(`
		SELECT id, lease_id, type, amount, due_date, paid_at, note, created_at, updated_at
		FROM bills WHERE id = ?
		`, id)

	var temp tempBill
	err := row.Scan(
		&temp.ID, &temp.LeaseID, &temp.Type, &temp.Amount,
		&temp.DueDate, &temp.PaidAt, &temp.Note, &temp.CreatedAt, &temp.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Now construct the bill using NewBill, then set the fields
	bill := billmodel.NewBill(temp.ID, temp.LeaseID, billmodel.BillType(temp.Type), temp.Amount, temp.DueDate, temp.Note)
	bill.PaidAt = temp.PaidAt
	bill.CreatedAt = temp.CreatedAt
	bill.UpdatedAt = temp.UpdatedAt

	return bill, nil
}

// FindAll 查找所有账单
func (r *BillRepository) FindAll() ([]*billmodel.Bill, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, lease_id, type, amount, due_date, paid_at, note, created_at, updated_at
		FROM bills ORDER BY created_at DESC
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []*billmodel.Bill
	for rows.Next() {
		var temp tempBill
		err := rows.Scan(
			&temp.ID, &temp.LeaseID, &temp.Type, &temp.Amount,
			&temp.DueDate, &temp.PaidAt, &temp.Note, &temp.CreatedAt, &temp.UpdatedAt)
		if err != nil {
			return nil, err
		}

		bill := billmodel.NewBill(temp.ID, temp.LeaseID, billmodel.BillType(temp.Type), temp.Amount, temp.DueDate, temp.Note)
		bill.PaidAt = temp.PaidAt
		bill.CreatedAt = temp.CreatedAt
		bill.UpdatedAt = temp.UpdatedAt

		bills = append(bills, bill)
	}
	return bills, nil
}

// FindByLeaseID 根据租约ID查找账单
func (r *BillRepository) FindByLeaseID(leaseID string) ([]*billmodel.Bill, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, lease_id, type, amount, due_date, paid_at, note, created_at, updated_at
		FROM bills WHERE lease_id = ? ORDER BY created_at DESC
		`, leaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []*billmodel.Bill
	for rows.Next() {
		var temp tempBill
		err := rows.Scan(
			&temp.ID, &temp.LeaseID, &temp.Type, &temp.Amount,
			&temp.DueDate, &temp.PaidAt, &temp.Note, &temp.CreatedAt, &temp.UpdatedAt)
		if err != nil {
			return nil, err
		}

		bill := billmodel.NewBill(temp.ID, temp.LeaseID, billmodel.BillType(temp.Type), temp.Amount, temp.DueDate, temp.Note)
		bill.PaidAt = temp.PaidAt
		bill.CreatedAt = temp.CreatedAt
		bill.UpdatedAt = temp.UpdatedAt

		bills = append(bills, bill)
	}
	return bills, nil
}

// FindUnpaidBillsDueBefore 查找到期前未支付的账单
func (r *BillRepository) FindUnpaidBillsDueBefore(dueDate time.Time) ([]*billmodel.Bill, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, lease_id, type, amount, due_date, paid_at, note, created_at, updated_at
		FROM bills WHERE due_date <= ? AND paid_at IS NULL ORDER BY due_date ASC
		`, dueDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []*billmodel.Bill
	for rows.Next() {
		var temp tempBill
		err := rows.Scan(
			&temp.ID, &temp.LeaseID, &temp.Type, &temp.Amount,
			&temp.DueDate, &temp.PaidAt, &temp.Note, &temp.CreatedAt, &temp.UpdatedAt)
		if err != nil {
			return nil, err
		}

		bill := billmodel.NewBill(temp.ID, temp.LeaseID, billmodel.BillType(temp.Type), temp.Amount, temp.DueDate, temp.Note)
		bill.PaidAt = temp.PaidAt
		bill.CreatedAt = temp.CreatedAt
		bill.UpdatedAt = temp.UpdatedAt

		bills = append(bills, bill)
	}
	return bills, nil
}

// Delete 删除账单
func (r *BillRepository) Delete(id string) error {
	_, err := r.conn.DB().Exec("DELETE FROM bills WHERE id = ?", id)
	return err
}
