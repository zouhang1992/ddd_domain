package mysql

import (
	"database/sql"
	"time"

	leasemodel "github.com/zouhang1992/ddd_domain/internal/domain/lease/model"
	leaserepo "github.com/zouhang1992/ddd_domain/internal/domain/lease/repository"
)

// LeaseRepository MySQL 租约仓储实现
type LeaseRepository struct {
	conn *Connection
}

// NewLeaseRepository 创建租约仓储
func NewLeaseRepository(conn *Connection) leaserepo.LeaseRepository {
	return &LeaseRepository{conn: conn}
}

// tempLease is a temporary struct for scanning
type tempLease struct {
	ID            string
	RoomID        string
	LandlordID    sql.NullString
	TenantName    string
	TenantPhone   sql.NullString
	StartDate     time.Time
	EndDate       time.Time
	RentAmount    int64
	DepositAmount int64
	Status        string
	Note          sql.NullString
	LastChargeAt  *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Save 保存租约
func (r *LeaseRepository) Save(lease *leasemodel.Lease) error {
	var lastChargeAt interface{}
	if lease.LastChargeAt != nil {
		lastChargeAt = *lease.LastChargeAt
	}

	_, err := r.conn.DB().Exec(`
		INSERT INTO leases (
			id, room_id, landlord_id, tenant_name, tenant_phone,
			start_date, end_date, rent_amount, deposit_amount, status, note, last_charge_at,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			room_id = VALUES(room_id),
			landlord_id = VALUES(landlord_id),
			tenant_name = VALUES(tenant_name),
			tenant_phone = VALUES(tenant_phone),
			start_date = VALUES(start_date),
			end_date = VALUES(end_date),
			rent_amount = VALUES(rent_amount),
			deposit_amount = VALUES(deposit_amount),
			status = VALUES(status),
			note = VALUES(note),
			last_charge_at = VALUES(last_charge_at),
			updated_at = VALUES(updated_at)
		`,
		lease.IDField, lease.RoomID, lease.LandlordID, lease.TenantName, lease.TenantPhone,
		lease.StartDate, lease.EndDate, lease.RentAmount, lease.DepositAmount, string(lease.Status), lease.Note, lastChargeAt,
		lease.CreatedAt, lease.UpdatedAt)
	return err
}

// FindByID 根据ID查找租约
func (r *LeaseRepository) FindByID(id string) (*leasemodel.Lease, error) {
	row := r.conn.DB().QueryRow(`
		SELECT id, room_id, landlord_id, tenant_name, tenant_phone,
			start_date, end_date, rent_amount, deposit_amount, status, note, last_charge_at,
			created_at, updated_at
		FROM leases WHERE id = ?
		`, id)

	var temp tempLease
	err := row.Scan(
		&temp.ID, &temp.RoomID, &temp.LandlordID, &temp.TenantName, &temp.TenantPhone,
		&temp.StartDate, &temp.EndDate, &temp.RentAmount, &temp.DepositAmount, &temp.Status, &temp.Note, &temp.LastChargeAt,
		&temp.CreatedAt, &temp.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Now construct the lease using NewLease
	landlordID := ""
	if temp.LandlordID.Valid {
		landlordID = temp.LandlordID.String
	}
	tenantPhone := ""
	if temp.TenantPhone.Valid {
		tenantPhone = temp.TenantPhone.String
	}
	note := ""
	if temp.Note.Valid {
		note = temp.Note.String
	}
	lease := leasemodel.NewLease(temp.ID, temp.RoomID, landlordID, temp.TenantName, tenantPhone,
		temp.StartDate, temp.EndDate, temp.RentAmount, temp.DepositAmount, note)
	lease.Status = leasemodel.LeaseStatus(temp.Status)
	lease.LastChargeAt = temp.LastChargeAt
	lease.CreatedAt = temp.CreatedAt
	lease.UpdatedAt = temp.UpdatedAt
	// Clear events that were automatically added by NewLease
	lease.ClearEvents()

	return lease, nil
}

// FindAll 查找所有租约
func (r *LeaseRepository) FindAll() ([]*leasemodel.Lease, error) {
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

	var leases []*leasemodel.Lease
	for rows.Next() {
		var temp tempLease
		err := rows.Scan(
			&temp.ID, &temp.RoomID, &temp.LandlordID, &temp.TenantName, &temp.TenantPhone,
			&temp.StartDate, &temp.EndDate, &temp.RentAmount, &temp.DepositAmount, &temp.Status, &temp.Note, &temp.LastChargeAt,
			&temp.CreatedAt, &temp.UpdatedAt)
		if err != nil {
			return nil, err
		}

		landlordID := ""
		if temp.LandlordID.Valid {
			landlordID = temp.LandlordID.String
		}
		tenantPhone := ""
		if temp.TenantPhone.Valid {
			tenantPhone = temp.TenantPhone.String
		}
		note := ""
		if temp.Note.Valid {
			note = temp.Note.String
		}
		lease := leasemodel.NewLease(temp.ID, temp.RoomID, landlordID, temp.TenantName, tenantPhone,
			temp.StartDate, temp.EndDate, temp.RentAmount, temp.DepositAmount, note)
		lease.Status = leasemodel.LeaseStatus(temp.Status)
		lease.LastChargeAt = temp.LastChargeAt
		lease.CreatedAt = temp.CreatedAt
		lease.UpdatedAt = temp.UpdatedAt
		lease.ClearEvents()

		leases = append(leases, lease)
	}
	return leases, nil
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

// FindActiveLeasesExpiringBefore 查找即将过期的生效租约
func (r *LeaseRepository) FindActiveLeasesExpiringBefore(expireTime time.Time) ([]*leasemodel.Lease, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, room_id, landlord_id, tenant_name, tenant_phone,
			start_date, end_date, rent_amount, deposit_amount, status, note, last_charge_at,
			created_at, updated_at
		FROM leases WHERE status = 'active' AND end_date <= ? ORDER BY end_date ASC
		`, expireTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leases []*leasemodel.Lease
	for rows.Next() {
		var temp tempLease
		err := rows.Scan(
			&temp.ID, &temp.RoomID, &temp.LandlordID, &temp.TenantName, &temp.TenantPhone,
			&temp.StartDate, &temp.EndDate, &temp.RentAmount, &temp.DepositAmount, &temp.Status, &temp.Note, &temp.LastChargeAt,
			&temp.CreatedAt, &temp.UpdatedAt)
		if err != nil {
			return nil, err
		}

		landlordID := ""
		if temp.LandlordID.Valid {
			landlordID = temp.LandlordID.String
		}
		tenantPhone := ""
		if temp.TenantPhone.Valid {
			tenantPhone = temp.TenantPhone.String
		}
		note := ""
		if temp.Note.Valid {
			note = temp.Note.String
		}
		lease := leasemodel.NewLease(temp.ID, temp.RoomID, landlordID, temp.TenantName, tenantPhone,
			temp.StartDate, temp.EndDate, temp.RentAmount, temp.DepositAmount, note)
		lease.Status = leasemodel.LeaseStatus(temp.Status)
		lease.LastChargeAt = temp.LastChargeAt
		lease.CreatedAt = temp.CreatedAt
		lease.UpdatedAt = temp.UpdatedAt
		lease.ClearEvents()

		leases = append(leases, lease)
	}
	return leases, nil
}

// HasActiveLeaseForRoom 检查房间是否有活跃的租约（pending 或 active 状态）
func (r *LeaseRepository) HasActiveLeaseForRoom(roomID string) (bool, error) {
	var count int
	row := r.conn.DB().QueryRow(`
		SELECT COUNT(*) FROM leases
		WHERE room_id = ? AND status IN ('pending', 'active')
		`, roomID)
	err := row.Scan(&count)
	return count > 0, err
}
