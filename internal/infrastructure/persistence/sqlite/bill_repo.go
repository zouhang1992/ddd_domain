package sqlite

import (
	"database/sql"
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
)

// BillRepository SQLite 账单仓储实现
type BillRepository struct {
	conn *Connection
}

// NewBillRepository 创建账单仓储
func NewBillRepository(conn *Connection) repository.BillRepository {
	return &BillRepository{conn: conn}
}

// Save 保存账单
func (r *BillRepository) Save(bill *model.Bill) error {
	var paidAt interface{}
	if bill.PaidAt != nil {
		paidAt = *bill.PaidAt
	}

	_, err := r.conn.DB().Exec(`
		INSERT OR REPLACE INTO bills (
			id, lease_id, type, status, amount, rent_amount, water_amount,
			electric_amount, other_amount, paid_at, note, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		bill.ID, bill.LeaseID, string(bill.Type), string(bill.Status), bill.Amount,
		bill.RentAmount, bill.WaterAmount, bill.ElectricAmount, bill.OtherAmount,
		paidAt, bill.Note, bill.CreatedAt, bill.UpdatedAt)
	return err
}

// FindByID 根据ID查找账单
func (r *BillRepository) FindByID(id string) (*model.Bill, error) {
	row := r.conn.DB().QueryRow(`
		SELECT id, lease_id, type, status, amount, rent_amount, water_amount,
			electric_amount, other_amount, paid_at, note, created_at, updated_at
		FROM bills WHERE id = ?
	`, id)

	bill := &model.Bill{}
	var typeStr, statusStr string
	var paidAt interface{}
	err := row.Scan(
		&bill.ID, &bill.LeaseID, &typeStr, &statusStr, &bill.Amount,
		&bill.RentAmount, &bill.WaterAmount, &bill.ElectricAmount, &bill.OtherAmount,
		&paidAt, &bill.Note, &bill.CreatedAt, &bill.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	bill.Type = model.BillType(typeStr)
	bill.Status = model.BillStatus(statusStr)
	if paidAt != nil {
		if t, ok := paidAt.(time.Time); ok {
			bill.PaidAt = &t
		}
	}

	return bill, nil
}

// FindAll 查找所有账单
func (r *BillRepository) FindAll() ([]*model.Bill, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, lease_id, type, status, amount, rent_amount, water_amount,
			electric_amount, other_amount, paid_at, note, created_at, updated_at
		FROM bills ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []*model.Bill
	for rows.Next() {
		bill := &model.Bill{}
		var typeStr, statusStr string
		var paidAt interface{}
		err := rows.Scan(
			&bill.ID, &bill.LeaseID, &typeStr, &statusStr, &bill.Amount,
			&bill.RentAmount, &bill.WaterAmount, &bill.ElectricAmount, &bill.OtherAmount,
			&paidAt, &bill.Note, &bill.CreatedAt, &bill.UpdatedAt)
		if err != nil {
			return nil, err
		}

		bill.Type = model.BillType(typeStr)
		bill.Status = model.BillStatus(statusStr)
		if paidAt != nil {
			if t, ok := paidAt.(time.Time); ok {
				bill.PaidAt = &t
			}
		}

		bills = append(bills, bill)
	}
	return bills, nil
}

// FindByLeaseID 根据租约ID查找账单
func (r *BillRepository) FindByLeaseID(leaseID string) ([]*model.Bill, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, lease_id, type, status, amount, rent_amount, water_amount,
			electric_amount, other_amount, paid_at, note, created_at, updated_at
		FROM bills WHERE lease_id = ? ORDER BY created_at DESC
	`, leaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []*model.Bill
	for rows.Next() {
		bill := &model.Bill{}
		var typeStr, statusStr string
		var paidAt interface{}
		err := rows.Scan(
			&bill.ID, &bill.LeaseID, &typeStr, &statusStr, &bill.Amount,
			&bill.RentAmount, &bill.WaterAmount, &bill.ElectricAmount, &bill.OtherAmount,
			&paidAt, &bill.Note, &bill.CreatedAt, &bill.UpdatedAt)
		if err != nil {
			return nil, err
		}

		bill.Type = model.BillType(typeStr)
		bill.Status = model.BillStatus(statusStr)
		if paidAt != nil {
			if t, ok := paidAt.(time.Time); ok {
				bill.PaidAt = &t
			}
		}

		bills = append(bills, bill)
	}
	return bills, nil
}

// FindByRoomID 根据房间ID查找账单
func (r *BillRepository) FindByRoomID(roomID string) ([]*model.Bill, error) {
	rows, err := r.conn.DB().Query(`
		SELECT b.id, b.lease_id, b.type, b.status, b.amount, b.rent_amount, b.water_amount,
			b.electric_amount, b.other_amount, b.paid_at, b.note, b.created_at, b.updated_at
		FROM bills b
		INNER JOIN leases l ON b.lease_id = l.id
		WHERE l.room_id = ?
		ORDER BY b.created_at DESC
	`, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []*model.Bill
	for rows.Next() {
		bill := &model.Bill{}
		var typeStr, statusStr string
		var paidAt interface{}
		err := rows.Scan(
			&bill.ID, &bill.LeaseID, &typeStr, &statusStr, &bill.Amount,
			&bill.RentAmount, &bill.WaterAmount, &bill.ElectricAmount, &bill.OtherAmount,
			&paidAt, &bill.Note, &bill.CreatedAt, &bill.UpdatedAt)
		if err != nil {
			return nil, err
		}

		bill.Type = model.BillType(typeStr)
		bill.Status = model.BillStatus(statusStr)
		if paidAt != nil {
			if t, ok := paidAt.(time.Time); ok {
				bill.PaidAt = &t
			}
		}

		bills = append(bills, bill)
	}
	return bills, nil
}

// FindByMonth 根据月份查找账单
func (r *BillRepository) FindByMonth(year int, month time.Month) ([]*model.Bill, error) {
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	rows, err := r.conn.DB().Query(`
		SELECT id, lease_id, type, status, amount, rent_amount, water_amount,
			electric_amount, other_amount, paid_at, note, created_at, updated_at
		FROM bills
		WHERE paid_at >= ? AND paid_at < ?
		ORDER BY paid_at DESC
	`, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []*model.Bill
	for rows.Next() {
		bill := &model.Bill{}
		var typeStr, statusStr string
		var paidAt interface{}
		err := rows.Scan(
			&bill.ID, &bill.LeaseID, &typeStr, &statusStr, &bill.Amount,
			&bill.RentAmount, &bill.WaterAmount, &bill.ElectricAmount, &bill.OtherAmount,
			&paidAt, &bill.Note, &bill.CreatedAt, &bill.UpdatedAt)
		if err != nil {
			return nil, err
		}

		bill.Type = model.BillType(typeStr)
		bill.Status = model.BillStatus(statusStr)
		if paidAt != nil {
			if t, ok := paidAt.(time.Time); ok {
				bill.PaidAt = &t
			}
		}

		bills = append(bills, bill)
	}
	return bills, nil
}

// Delete 删除账单
func (r *BillRepository) Delete(id string) error {
	_, err := r.conn.DB().Exec("DELETE FROM bills WHERE id = ?", id)
	return err
}
