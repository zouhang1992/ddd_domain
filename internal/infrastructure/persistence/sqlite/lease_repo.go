package sqlite

import (
	"database/sql"
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
)

// LeaseRepository SQLite 租约仓储实现
type LeaseRepository struct {
	conn *Connection
}

// NewLeaseRepository 创建租约仓储
func NewLeaseRepository(conn *Connection) repository.LeaseRepository {
	return &LeaseRepository{conn: conn}
}

// Save 保存租约
func (r *LeaseRepository) Save(lease *model.Lease) error {
	var lastChargeAt interface{}
	if lease.LastChargeAt != nil {
		lastChargeAt = *lease.LastChargeAt
	}

	_, err := r.conn.DB().Exec(`
		INSERT OR REPLACE INTO leases (
			id, room_id, landlord_id, tenant_name, tenant_phone,
			start_date, end_date, rent_amount, deposit_amount, status, note, last_charge_at,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		lease.ID, lease.RoomID, lease.LandlordID, lease.TenantName, lease.TenantPhone,
		lease.StartDate, lease.EndDate, lease.RentAmount, lease.DepositAmount, string(lease.Status), lease.Note, lastChargeAt,
		lease.CreatedAt, lease.UpdatedAt)
	return err
}

// FindByID 根据ID查找租约
func (r *LeaseRepository) FindByID(id string) (*model.Lease, error) {
	row := r.conn.DB().QueryRow(`
		SELECT id, room_id, landlord_id, tenant_name, tenant_phone,
			start_date, end_date, rent_amount, deposit_amount, status, note, last_charge_at,
			created_at, updated_at
		FROM leases WHERE id = ?
	`, id)

	lease := &model.Lease{}
	var statusStr string
	var lastChargeAt interface{}
	err := row.Scan(
		&lease.ID, &lease.RoomID, &lease.LandlordID, &lease.TenantName, &lease.TenantPhone,
		&lease.StartDate, &lease.EndDate, &lease.RentAmount, &lease.DepositAmount, &statusStr, &lease.Note, &lastChargeAt,
		&lease.CreatedAt, &lease.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	lease.Status = model.LeaseStatus(statusStr)
	if lastChargeAt != nil {
		if t, ok := lastChargeAt.(time.Time); ok {
			lease.LastChargeAt = &t
		}
	}

	return lease, nil
}

// FindAll 查找所有租约
func (r *LeaseRepository) FindAll() ([]*model.Lease, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, room_id, landlord_id, tenant_name, tenant_phone,
			start_date, end_date, rent_amount, deposit_amount, status, note, last_charge_at,
			created_at, updated_at
		FROM leases ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leases []*model.Lease
	for rows.Next() {
		lease := &model.Lease{}
		var statusStr string
		var lastChargeAt interface{}
		err := rows.Scan(
			&lease.ID, &lease.RoomID, &lease.LandlordID, &lease.TenantName, &lease.TenantPhone,
			&lease.StartDate, &lease.EndDate, &lease.RentAmount, &lease.DepositAmount, &statusStr, &lease.Note, &lastChargeAt,
			&lease.CreatedAt, &lease.UpdatedAt)
		if err != nil {
			return nil, err
		}

		lease.Status = model.LeaseStatus(statusStr)
		if lastChargeAt != nil {
			if t, ok := lastChargeAt.(time.Time); ok {
				lease.LastChargeAt = &t
			}
		}

		leases = append(leases, lease)
	}
	return leases, nil
}

// FindByRoomID 根据房间ID查找租约
func (r *LeaseRepository) FindByRoomID(roomID string) ([]*model.Lease, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, room_id, landlord_id, tenant_name, tenant_phone,
			start_date, end_date, rent_amount, deposit_amount, status, note, last_charge_at,
			created_at, updated_at
		FROM leases WHERE room_id = ? ORDER BY created_at DESC
	`, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leases []*model.Lease
	for rows.Next() {
		lease := &model.Lease{}
		var statusStr string
		var lastChargeAt interface{}
		err := rows.Scan(
			&lease.ID, &lease.RoomID, &lease.LandlordID, &lease.TenantName, &lease.TenantPhone,
			&lease.StartDate, &lease.EndDate, &lease.RentAmount, &lease.DepositAmount, &statusStr, &lease.Note, &lastChargeAt,
			&lease.CreatedAt, &lease.UpdatedAt)
		if err != nil {
			return nil, err
		}

		lease.Status = model.LeaseStatus(statusStr)
		if lastChargeAt != nil {
			if t, ok := lastChargeAt.(time.Time); ok {
				lease.LastChargeAt = &t
			}
		}

		leases = append(leases, lease)
	}
	return leases, nil
}

// FindByStatus 根据状态查找租约
func (r *LeaseRepository) FindByStatus(status model.LeaseStatus) ([]*model.Lease, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, room_id, landlord_id, tenant_name, tenant_phone,
			start_date, end_date, rent_amount, deposit_amount, status, note, last_charge_at,
			created_at, updated_at
		FROM leases WHERE status = ? ORDER BY created_at DESC
	`, string(status))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leases []*model.Lease
	for rows.Next() {
		lease := &model.Lease{}
		var statusStr string
		var lastChargeAt interface{}
		err := rows.Scan(
			&lease.ID, &lease.RoomID, &lease.LandlordID, &lease.TenantName, &lease.TenantPhone,
			&lease.StartDate, &lease.EndDate, &lease.RentAmount, &lease.DepositAmount, &statusStr, &lease.Note, &lastChargeAt,
			&lease.CreatedAt, &lease.UpdatedAt)
		if err != nil {
			return nil, err
		}

		lease.Status = model.LeaseStatus(statusStr)
		if lastChargeAt != nil {
			if t, ok := lastChargeAt.(time.Time); ok {
				lease.LastChargeAt = &t
			}
		}

		leases = append(leases, lease)
	}
	return leases, nil
}

// FindByRoomIDAndStatus 根据房间ID和状态查找租约
func (r *LeaseRepository) FindByRoomIDAndStatus(roomID string, status model.LeaseStatus) ([]*model.Lease, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, room_id, landlord_id, tenant_name, tenant_phone,
			start_date, end_date, rent_amount, deposit_amount, status, note, last_charge_at,
			created_at, updated_at
		FROM leases WHERE room_id = ? AND status = ? ORDER BY created_at DESC
	`, roomID, string(status))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leases []*model.Lease
	for rows.Next() {
		lease := &model.Lease{}
		var statusStr string
		var lastChargeAt interface{}
		err := rows.Scan(
			&lease.ID, &lease.RoomID, &lease.LandlordID, &lease.TenantName, &lease.TenantPhone,
			&lease.StartDate, &lease.EndDate, &lease.RentAmount, &lease.DepositAmount, &statusStr, &lease.Note, &lastChargeAt,
			&lease.CreatedAt, &lease.UpdatedAt)
		if err != nil {
			return nil, err
		}

		lease.Status = model.LeaseStatus(statusStr)
		if lastChargeAt != nil {
			if t, ok := lastChargeAt.(time.Time); ok {
				lease.LastChargeAt = &t
			}
		}

		leases = append(leases, lease)
	}
	return leases, nil
}

// FindActiveByRoomID 查找房间的生效中租约
func (r *LeaseRepository) FindActiveByRoomID(roomID string) (*model.Lease, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, room_id, landlord_id, tenant_name, tenant_phone,
			start_date, end_date, rent_amount, deposit_amount, status, note, last_charge_at,
			created_at, updated_at
		FROM leases WHERE room_id = ? AND status = 'active' ORDER BY created_at DESC LIMIT 1
	`, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		lease := &model.Lease{}
		var statusStr string
		var lastChargeAt interface{}
		err := rows.Scan(
			&lease.ID, &lease.RoomID, &lease.LandlordID, &lease.TenantName, &lease.TenantPhone,
			&lease.StartDate, &lease.EndDate, &lease.RentAmount, &lease.DepositAmount, &statusStr, &lease.Note, &lastChargeAt,
			&lease.CreatedAt, &lease.UpdatedAt)
		if err != nil {
			return nil, err
		}

		lease.Status = model.LeaseStatus(statusStr)
		if lastChargeAt != nil {
			if t, ok := lastChargeAt.(time.Time); ok {
				lease.LastChargeAt = &t
			}
		}

		return lease, nil
	}
	return nil, nil
}

// Delete 删除租约
func (r *LeaseRepository) Delete(id string) error {
	_, err := r.conn.DB().Exec("DELETE FROM leases WHERE id = ?", id)
	return err
}

// HasBills 检查是否有账单
func (r *LeaseRepository) HasBills(leaseID string) (bool, error) {
	var count int
	row := r.conn.DB().QueryRow("SELECT COUNT(*) FROM bills WHERE lease_id = ?", leaseID)
	err := row.Scan(&count)
	return count > 0, err
}

// HasDeposit 检查是否有押金
func (r *LeaseRepository) HasDeposit(leaseID string) (bool, error) {
	var count int
	row := r.conn.DB().QueryRow("SELECT COUNT(*) FROM deposits WHERE lease_id = ?", leaseID)
	err := row.Scan(&count)
	return count > 0, err
}

// FindByCriteria 按条件查找租约
func (r *LeaseRepository) FindByCriteria(criteria repository.LeaseCriteria, offset, limit int) ([]*model.Lease, error) {
	query := `
		SELECT id, room_id, landlord_id, tenant_name, tenant_phone,
			start_date, end_date, rent_amount, deposit_amount, status, note, last_charge_at,
			created_at, updated_at
		FROM leases
		WHERE 1 = 1
	`
	var args []interface{}

	if criteria.TenantName != "" {
		query += " AND tenant_name LIKE ?"
		args = append(args, "%"+criteria.TenantName+"%")
	}
	if criteria.TenantPhone != "" {
		query += " AND tenant_phone LIKE ?"
		args = append(args, "%"+criteria.TenantPhone+"%")
	}
	if criteria.Status != "" {
		query += " AND status = ?"
		args = append(args, criteria.Status)
	}
	if criteria.RoomID != "" {
		query += " AND room_id = ?"
		args = append(args, criteria.RoomID)
	}
	if criteria.StartDate != nil {
		query += " AND start_date >= ?"
		args = append(args, criteria.StartDate)
	}
	if criteria.EndDate != nil {
		query += " AND end_date <= ?"
		args = append(args, criteria.EndDate)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.conn.DB().Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leases []*model.Lease
	for rows.Next() {
		lease := &model.Lease{}
		var statusStr string
		var lastChargeAt interface{}
		err := rows.Scan(
			&lease.ID, &lease.RoomID, &lease.LandlordID, &lease.TenantName, &lease.TenantPhone,
			&lease.StartDate, &lease.EndDate, &lease.RentAmount, &lease.DepositAmount, &statusStr, &lease.Note, &lastChargeAt,
			&lease.CreatedAt, &lease.UpdatedAt)
		if err != nil {
			return nil, err
		}

		lease.Status = model.LeaseStatus(statusStr)
		if lastChargeAt != nil {
			if t, ok := lastChargeAt.(time.Time); ok {
				lease.LastChargeAt = &t
			}
		}

		leases = append(leases, lease)
	}
	return leases, nil
}

// CountByCriteria 按条件统计租约数量
func (r *LeaseRepository) CountByCriteria(criteria repository.LeaseCriteria) (int, error) {
	query := `
		SELECT COUNT(*) FROM leases
		WHERE 1 = 1
	`
	var args []interface{}

	if criteria.TenantName != "" {
		query += " AND tenant_name LIKE ?"
		args = append(args, "%"+criteria.TenantName+"%")
	}
	if criteria.TenantPhone != "" {
		query += " AND tenant_phone LIKE ?"
		args = append(args, "%"+criteria.TenantPhone+"%")
	}
	if criteria.Status != "" {
		query += " AND status = ?"
		args = append(args, criteria.Status)
	}
	if criteria.RoomID != "" {
		query += " AND room_id = ?"
		args = append(args, criteria.RoomID)
	}
	if criteria.StartDate != nil {
		query += " AND start_date >= ?"
		args = append(args, criteria.StartDate)
	}
	if criteria.EndDate != nil {
		query += " AND end_date <= ?"
		args = append(args, criteria.EndDate)
	}

	var count int
	row := r.conn.DB().QueryRow(query, args...)
	err := row.Scan(&count)
	return count, err
}
