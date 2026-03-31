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
			id, lease_id, type, status, amount, rent_amount, water_amount, electric_amount, other_amount, refund_deposit_amount, bill_start, bill_end, due_date, paid_at, note, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
		bill.IDField, bill.LeaseID, string(bill.Type), string(bill.Status), bill.Amount, bill.RentAmount, bill.WaterAmount, bill.ElectricAmount, bill.OtherAmount, bill.RefundDepositAmount, bill.BillStart, bill.BillEnd, bill.DueDate,
		paidAt, bill.Note, bill.CreatedAt, bill.UpdatedAt)
	return err
}

// tempBill is a temporary struct for scanning
type tempBill struct {
	ID                  string
	LeaseID             string
	Type                string
	Status              string
	Amount              int64
	RentAmount          int64
	WaterAmount         int64
	ElectricAmount      int64
	OtherAmount         int64
	RefundDepositAmount int64
	BillStart           sql.NullTime
	BillEnd             sql.NullTime
	DueDate             sql.NullTime
	PaidAt              *time.Time
	Note                string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// FindByID 根据ID查找账单
func (r *BillRepository) FindByID(id string) (*billmodel.Bill, error) {
	row := r.conn.DB().QueryRow(`
		SELECT id, lease_id, type, status, amount, rent_amount, water_amount, electric_amount, other_amount, refund_deposit_amount, bill_start, bill_end, due_date, paid_at, note, created_at, updated_at
		FROM bills WHERE id = ?
		`, id)

	var temp tempBill
	err := row.Scan(
		&temp.ID, &temp.LeaseID, &temp.Type, &temp.Status, &temp.Amount, &temp.RentAmount, &temp.WaterAmount, &temp.ElectricAmount, &temp.OtherAmount, &temp.RefundDepositAmount,
		&temp.BillStart, &temp.BillEnd, &temp.DueDate, &temp.PaidAt, &temp.Note, &temp.CreatedAt, &temp.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	billStart := time.Now()
	if temp.BillStart.Valid {
		billStart = temp.BillStart.Time
	}
	billEnd := time.Now()
	if temp.BillEnd.Valid {
		billEnd = temp.BillEnd.Time
	}
	dueDate := time.Now()
	if temp.DueDate.Valid {
		dueDate = temp.DueDate.Time
	}

	// Now construct the bill using NewBill, then set the fields
	bill := billmodel.NewBill(temp.ID, temp.LeaseID, billmodel.BillType(temp.Type), temp.Amount, billStart, billEnd, dueDate, temp.Note)
	bill.Status = billmodel.BillStatus(temp.Status)
	bill.RentAmount = temp.RentAmount
	bill.WaterAmount = temp.WaterAmount
	bill.ElectricAmount = temp.ElectricAmount
	bill.OtherAmount = temp.OtherAmount
	bill.RefundDepositAmount = temp.RefundDepositAmount
	bill.PaidAt = temp.PaidAt
	bill.CreatedAt = temp.CreatedAt
	bill.UpdatedAt = temp.UpdatedAt
	bill.ClearEvents()

	return bill, nil
}

// FindAll 查找所有账单
func (r *BillRepository) FindAll() ([]*billmodel.Bill, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, lease_id, type, status, amount, rent_amount, water_amount, electric_amount, other_amount, refund_deposit_amount, bill_start, bill_end, due_date, paid_at, note, created_at, updated_at
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
			&temp.ID, &temp.LeaseID, &temp.Type, &temp.Status, &temp.Amount, &temp.RentAmount, &temp.WaterAmount, &temp.ElectricAmount, &temp.OtherAmount, &temp.RefundDepositAmount,
			&temp.BillStart, &temp.BillEnd, &temp.DueDate, &temp.PaidAt, &temp.Note, &temp.CreatedAt, &temp.UpdatedAt)
		if err != nil {
			return nil, err
		}

		billStart := time.Now()
		if temp.BillStart.Valid {
			billStart = temp.BillStart.Time
		}
		billEnd := time.Now()
		if temp.BillEnd.Valid {
			billEnd = temp.BillEnd.Time
		}
		dueDate := time.Now()
		if temp.DueDate.Valid {
			dueDate = temp.DueDate.Time
		}

		bill := billmodel.NewBill(temp.ID, temp.LeaseID, billmodel.BillType(temp.Type), temp.Amount, billStart, billEnd, dueDate, temp.Note)
		bill.Status = billmodel.BillStatus(temp.Status)
		bill.RentAmount = temp.RentAmount
		bill.WaterAmount = temp.WaterAmount
		bill.ElectricAmount = temp.ElectricAmount
		bill.OtherAmount = temp.OtherAmount
		bill.RefundDepositAmount = temp.RefundDepositAmount
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
		SELECT id, lease_id, type, status, amount, rent_amount, water_amount, electric_amount, other_amount, refund_deposit_amount, bill_start, bill_end, due_date, paid_at, note, created_at, updated_at
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
			&temp.ID, &temp.LeaseID, &temp.Type, &temp.Status, &temp.Amount, &temp.RentAmount, &temp.WaterAmount, &temp.ElectricAmount, &temp.OtherAmount, &temp.RefundDepositAmount,
			&temp.BillStart, &temp.BillEnd, &temp.DueDate, &temp.PaidAt, &temp.Note, &temp.CreatedAt, &temp.UpdatedAt)
		if err != nil {
			return nil, err
		}

		billStart := time.Now()
		if temp.BillStart.Valid {
			billStart = temp.BillStart.Time
		}
		billEnd := time.Now()
		if temp.BillEnd.Valid {
			billEnd = temp.BillEnd.Time
		}
		dueDate := time.Now()
		if temp.DueDate.Valid {
			dueDate = temp.DueDate.Time
		}

		bill := billmodel.NewBill(temp.ID, temp.LeaseID, billmodel.BillType(temp.Type), temp.Amount, billStart, billEnd, dueDate, temp.Note)
		bill.Status = billmodel.BillStatus(temp.Status)
		bill.RentAmount = temp.RentAmount
		bill.WaterAmount = temp.WaterAmount
		bill.ElectricAmount = temp.ElectricAmount
		bill.OtherAmount = temp.OtherAmount
		bill.RefundDepositAmount = temp.RefundDepositAmount
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
		SELECT id, lease_id, type, status, amount, rent_amount, water_amount, electric_amount, other_amount, refund_deposit_amount, bill_start, bill_end, due_date, paid_at, note, created_at, updated_at
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
			&temp.ID, &temp.LeaseID, &temp.Type, &temp.Status, &temp.Amount, &temp.RentAmount, &temp.WaterAmount, &temp.ElectricAmount, &temp.OtherAmount, &temp.RefundDepositAmount,
			&temp.BillStart, &temp.BillEnd, &temp.DueDate, &temp.PaidAt, &temp.Note, &temp.CreatedAt, &temp.UpdatedAt)
		if err != nil {
			return nil, err
		}

		billStart := time.Now()
		if temp.BillStart.Valid {
			billStart = temp.BillStart.Time
		}
		billEnd := time.Now()
		if temp.BillEnd.Valid {
			billEnd = temp.BillEnd.Time
		}
		billDueDate := time.Now()
		if temp.DueDate.Valid {
			billDueDate = temp.DueDate.Time
		}

		bill := billmodel.NewBill(temp.ID, temp.LeaseID, billmodel.BillType(temp.Type), temp.Amount, billStart, billEnd, billDueDate, temp.Note)
		bill.Status = billmodel.BillStatus(temp.Status)
		bill.RentAmount = temp.RentAmount
		bill.WaterAmount = temp.WaterAmount
		bill.ElectricAmount = temp.ElectricAmount
		bill.OtherAmount = temp.OtherAmount
		bill.RefundDepositAmount = temp.RefundDepositAmount
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
