package sqlite

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	printmodel "github.com/zouhang1992/ddd_domain/internal/domain/print/model"
	printrepo "github.com/zouhang1992/ddd_domain/internal/domain/print/repository"
)

// PrintJobRepository SQLite 打印作业仓储实现
type PrintJobRepository struct {
	conn *Connection
}

// NewPrintJobRepository 创建打印作业仓储
func NewPrintJobRepository(conn *Connection) printrepo.PrintJobRepository {
	return &PrintJobRepository{conn: conn}
}

// tempPrintJob is a temporary struct for scanning
type tempPrintJob struct {
	ID           string
	Type         string
	Status       string
	ReferenceID  string
	TenantName   sql.NullString
	TenantPhone  sql.NullString
	RoomID       sql.NullString
	RoomNumber   sql.NullString
	Address      sql.NullString
	LandlordName sql.NullString
	Amount       int64
	ErrorMsg     sql.NullString
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CompletedAt  sql.NullTime
}

// Save 保存打印作业
func (r *PrintJobRepository) Save(job *printmodel.PrintJob) error {
	if job.ID() == "" {
		job.IDField = uuid.NewString()
	}

	// 检查是否已存在
	var exists bool
	err := r.conn.DB().QueryRow("SELECT EXISTS(SELECT 1 FROM print_jobs WHERE id = ?)", job.ID()).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		// 更新
		_, err = r.conn.DB().Exec(`
			UPDATE print_jobs SET
				type = ?, status = ?, reference_id = ?, tenant_name = ?,
				tenant_phone = ?, room_id = ?, room_number = ?, address = ?, landlord_name = ?,
				amount = ?, error_msg = ?, updated_at = ?, completed_at = ?
			WHERE id = ?
			`,
			job.Type, job.Status, job.ReferenceID, job.TenantName,
			job.TenantPhone, job.RoomID, job.RoomNumber, job.Address, job.LandlordName,
			job.Amount, job.ErrorMsg, job.UpdatedAt, job.CompletedAt,
			job.ID())
		return err
	}

	// 插入
	_, err = r.conn.DB().Exec(`
		INSERT INTO print_jobs (
			id, type, status, reference_id, tenant_name, tenant_phone,
			room_id, room_number, address, landlord_name, amount,
			error_msg, created_at, updated_at, completed_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
		job.ID(), job.Type, job.Status, job.ReferenceID, job.TenantName,
		job.TenantPhone, job.RoomID, job.RoomNumber, job.Address, job.LandlordName,
		job.Amount, job.ErrorMsg, job.CreatedAt, job.UpdatedAt, job.CompletedAt)
	return err
}

// FindByID 根据ID查找打印作业
func (r *PrintJobRepository) FindByID(id string) (*printmodel.PrintJob, error) {
	row := r.conn.DB().QueryRow(`
		SELECT id, type, status, reference_id, tenant_name, tenant_phone,
			room_id, room_number, address, landlord_name, amount,
			error_msg, created_at, updated_at, completed_at
		FROM print_jobs WHERE id = ?
		`, id)

	var temp tempPrintJob
	err := row.Scan(
		&temp.ID, &temp.Type, &temp.Status, &temp.ReferenceID, &temp.TenantName, &temp.TenantPhone,
		&temp.RoomID, &temp.RoomNumber, &temp.Address, &temp.LandlordName, &temp.Amount,
		&temp.ErrorMsg, &temp.CreatedAt, &temp.UpdatedAt, &temp.CompletedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return tempToPrintJobModel(&temp)
}

// FindAll 分页查找所有打印作业
func (r *PrintJobRepository) FindAll(offset, limit int) ([]*printmodel.PrintJob, int, error) {
	// 获取总数
	var total int
	row := r.conn.DB().QueryRow(`SELECT COUNT(*) FROM print_jobs`)
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}

	// 分页查询
	rows, err := r.conn.DB().Query(`
		SELECT id, type, status, reference_id, tenant_name, tenant_phone,
			room_id, room_number, address, landlord_name, amount,
			error_msg, created_at, updated_at, completed_at
		FROM print_jobs ORDER BY created_at DESC LIMIT ? OFFSET ?
		`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var jobs []*printmodel.PrintJob
	for rows.Next() {
		var temp tempPrintJob
		err := rows.Scan(
			&temp.ID, &temp.Type, &temp.Status, &temp.ReferenceID, &temp.TenantName, &temp.TenantPhone,
			&temp.RoomID, &temp.RoomNumber, &temp.Address, &temp.LandlordName, &temp.Amount,
			&temp.ErrorMsg, &temp.CreatedAt, &temp.UpdatedAt, &temp.CompletedAt)
		if err != nil {
			return nil, 0, err
		}

		job, err := tempToPrintJobModel(&temp)
		if err != nil {
			return nil, 0, err
		}
		jobs = append(jobs, job)
	}
	return jobs, total, nil
}

// FindByStatus 根据状态分页查找打印作业
func (r *PrintJobRepository) FindByStatus(status printmodel.PrintJobStatus, offset, limit int) ([]*printmodel.PrintJob, int, error) {
	// 获取总数
	var total int
	row := r.conn.DB().QueryRow(`SELECT COUNT(*) FROM print_jobs WHERE status = ?`, status)
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}

	// 分页查询
	rows, err := r.conn.DB().Query(`
		SELECT id, type, status, reference_id, tenant_name, tenant_phone,
			room_id, room_number, address, landlord_name, amount,
			error_msg, created_at, updated_at, completed_at
		FROM print_jobs WHERE status = ? ORDER BY created_at DESC LIMIT ? OFFSET ?
		`, status, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var jobs []*printmodel.PrintJob
	for rows.Next() {
		var temp tempPrintJob
		err := rows.Scan(
			&temp.ID, &temp.Type, &temp.Status, &temp.ReferenceID, &temp.TenantName, &temp.TenantPhone,
			&temp.RoomID, &temp.RoomNumber, &temp.Address, &temp.LandlordName, &temp.Amount,
			&temp.ErrorMsg, &temp.CreatedAt, &temp.UpdatedAt, &temp.CompletedAt)
		if err != nil {
			return nil, 0, err
		}

		job, err := tempToPrintJobModel(&temp)
		if err != nil {
			return nil, 0, err
		}
		jobs = append(jobs, job)
	}
	return jobs, total, nil
}

// FindByType 根据类型分页查找打印作业
func (r *PrintJobRepository) FindByType(jobType printmodel.PrintJobType, offset, limit int) ([]*printmodel.PrintJob, int, error) {
	// 获取总数
	var total int
	row := r.conn.DB().QueryRow(`SELECT COUNT(*) FROM print_jobs WHERE type = ?`, jobType)
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}

	// 分页查询
	rows, err := r.conn.DB().Query(`
		SELECT id, type, status, reference_id, tenant_name, tenant_phone,
			room_id, room_number, address, landlord_name, amount,
			error_msg, created_at, updated_at, completed_at
		FROM print_jobs WHERE type = ? ORDER BY created_at DESC LIMIT ? OFFSET ?
		`, jobType, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var jobs []*printmodel.PrintJob
	for rows.Next() {
		var temp tempPrintJob
		err := rows.Scan(
			&temp.ID, &temp.Type, &temp.Status, &temp.ReferenceID, &temp.TenantName, &temp.TenantPhone,
			&temp.RoomID, &temp.RoomNumber, &temp.Address, &temp.LandlordName, &temp.Amount,
			&temp.ErrorMsg, &temp.CreatedAt, &temp.UpdatedAt, &temp.CompletedAt)
		if err != nil {
			return nil, 0, err
		}

		job, err := tempToPrintJobModel(&temp)
		if err != nil {
			return nil, 0, err
		}
		jobs = append(jobs, job)
	}
	return jobs, total, nil
}

// FindByTimeRange 根据时间范围分页查找打印作业
func (r *PrintJobRepository) FindByTimeRange(start, end time.Time, offset, limit int) ([]*printmodel.PrintJob, int, error) {
	// 获取总数
	var total int
	row := r.conn.DB().QueryRow(`SELECT COUNT(*) FROM print_jobs WHERE created_at >= ? AND created_at <= ?`, start, end)
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}

	// 分页查询
	rows, err := r.conn.DB().Query(`
		SELECT id, type, status, reference_id, tenant_name, tenant_phone,
			room_id, room_number, address, landlord_name, amount,
			error_msg, created_at, updated_at, completed_at
		FROM print_jobs WHERE created_at >= ? AND created_at <= ? ORDER BY created_at DESC LIMIT ? OFFSET ?
		`, start, end, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var jobs []*printmodel.PrintJob
	for rows.Next() {
		var temp tempPrintJob
		err := rows.Scan(
			&temp.ID, &temp.Type, &temp.Status, &temp.ReferenceID, &temp.TenantName, &temp.TenantPhone,
			&temp.RoomID, &temp.RoomNumber, &temp.Address, &temp.LandlordName, &temp.Amount,
			&temp.ErrorMsg, &temp.CreatedAt, &temp.UpdatedAt, &temp.CompletedAt)
		if err != nil {
			return nil, 0, err
		}

		job, err := tempToPrintJobModel(&temp)
		if err != nil {
			return nil, 0, err
		}
		jobs = append(jobs, job)
	}
	return jobs, total, nil
}

// FindByFilters 根据多种筛选条件分页查找打印作业
func (r *PrintJobRepository) FindByFilters(status printmodel.PrintJobStatus, jobType printmodel.PrintJobType, start, end *time.Time, offset, limit int) ([]*printmodel.PrintJob, int, error) {
	// 构建查询条件
	args := []interface{}{}
	whereClauses := []string{}

	if status != "" {
		whereClauses = append(whereClauses, "status = ?")
		args = append(args, status)
	}
	if jobType != "" {
		whereClauses = append(whereClauses, "type = ?")
		args = append(args, jobType)
	}
	if start != nil {
		whereClauses = append(whereClauses, "created_at >= ?")
		args = append(args, *start)
	}
	if end != nil {
		whereClauses = append(whereClauses, "created_at <= ?")
		args = append(args, *end)
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = " WHERE "
		for i, clause := range whereClauses {
			if i > 0 {
				whereClause += " AND "
			}
			whereClause += clause
		}
	}

	// 获取总数
	countQuery := "SELECT COUNT(*) FROM print_jobs" + whereClause
	var total int
	row := r.conn.DB().QueryRow(countQuery, args...)
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}

	// 分页查询
	dataQuery := `
		SELECT id, type, status, reference_id, tenant_name, tenant_phone,
			room_id, room_number, address, landlord_name, amount,
			error_msg, created_at, updated_at, completed_at
		FROM print_jobs` + whereClause + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.conn.DB().Query(dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var jobs []*printmodel.PrintJob
	for rows.Next() {
		var temp tempPrintJob
		err := rows.Scan(
			&temp.ID, &temp.Type, &temp.Status, &temp.ReferenceID, &temp.TenantName, &temp.TenantPhone,
			&temp.RoomID, &temp.RoomNumber, &temp.Address, &temp.LandlordName, &temp.Amount,
			&temp.ErrorMsg, &temp.CreatedAt, &temp.UpdatedAt, &temp.CompletedAt)
		if err != nil {
			return nil, 0, err
		}

		job, err := tempToPrintJobModel(&temp)
		if err != nil {
			return nil, 0, err
		}
		jobs = append(jobs, job)
	}
	return jobs, total, nil
}

// Delete 删除打印作业
func (r *PrintJobRepository) Delete(id string) error {
	_, err := r.conn.DB().Exec("DELETE FROM print_jobs WHERE id = ?", id)
	return err
}

// tempToPrintJobModel 将临时结构转换为领域模型
func tempToPrintJobModel(temp *tempPrintJob) (*printmodel.PrintJob, error) {
	tenantName := ""
	if temp.TenantName.Valid {
		tenantName = temp.TenantName.String
	}

	tenantPhone := ""
	if temp.TenantPhone.Valid {
		tenantPhone = temp.TenantPhone.String
	}

	roomID := ""
	if temp.RoomID.Valid {
		roomID = temp.RoomID.String
	}

	roomNumber := ""
	if temp.RoomNumber.Valid {
		roomNumber = temp.RoomNumber.String
	}

	address := ""
	if temp.Address.Valid {
		address = temp.Address.String
	}

	landlordName := ""
	if temp.LandlordName.Valid {
		landlordName = temp.LandlordName.String
	}

	errorMsg := ""
	if temp.ErrorMsg.Valid {
		errorMsg = temp.ErrorMsg.String
	}

	var completedAt *time.Time
	if temp.CompletedAt.Valid {
		completedAt = &temp.CompletedAt.Time
	}

	job := printmodel.NewPrintJob(
		temp.ID,
		printmodel.PrintJobType(temp.Type),
		temp.ReferenceID,
		tenantName,
		tenantPhone,
		roomID,
		roomNumber,
		address,
		landlordName,
		temp.Amount,
	)
	job.Status = printmodel.PrintJobStatus(temp.Status)
	job.ErrorMsg = errorMsg
	job.CreatedAt = temp.CreatedAt
	job.UpdatedAt = temp.UpdatedAt
	job.CompletedAt = completedAt

	return job, nil
}
