package sqlite

import (
	"database/sql"
	"fmt"
)

// Migration 迁移接口
type Migration interface {
	Version() string
	Up(db *sql.Tx) error
	Down(db *sql.Tx) error
}

// migrations 所有迁移
var migrations = []Migration{
	&BaseMigration{},
	&AddDepositAmountToLeasesMigration{},
	&AddOperationLogsTableMigration{},
	&AddStatusAndNoteToRoomsMigration{},
	&AddDueDateToBillsMigration{},
	&AddRefundDepositAmountToBillsMigration{},
	&AddPrintJobsTableMigration{},
	&AddPrintJobDetailsMigration{},
	&AddBillPeriodToBillsMigration{},
}

// AddDepositAmountToLeasesMigration 为 leases 表添加 deposit_amount 字段
type AddDepositAmountToLeasesMigration struct{}

func (m *AddDepositAmountToLeasesMigration) Version() string {
	return "202603251330" // 格式：YYYYMMDDHHMM
}

func (m *AddDepositAmountToLeasesMigration) Up(tx *sql.Tx) error {
	// 检查列是否已存在
	rows, err := tx.Query("PRAGMA table_info(leases)")
	if err != nil {
		return fmt.Errorf("failed to check table info: %w", err)
	}
	defer rows.Close()

	colExists := false
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt_value sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt_value, &pk); err != nil {
			return fmt.Errorf("failed to scan table info: %w", err)
		}
		if name == "deposit_amount" {
			colExists = true
			break
		}
	}

	if !colExists {
		if _, err := tx.Exec("ALTER TABLE leases ADD COLUMN deposit_amount INTEGER NOT NULL DEFAULT 0"); err != nil {
			return fmt.Errorf("failed to add deposit_amount column: %w", err)
		}
	}
	return nil
}

func (m *AddDepositAmountToLeasesMigration) Down(tx *sql.Tx) error {
	// SQLite 不支持直接删除列，我们不处理回滚
	return nil
}

// AddOperationLogsTableMigration 创建 operation_logs 表
type AddOperationLogsTableMigration struct{}

func (m *AddOperationLogsTableMigration) Version() string {
	return "202603292100" // 格式：YYYYMMDDHHMM
}

func (m *AddOperationLogsTableMigration) Up(tx *sql.Tx) error {
	// 创建 operation_logs 表
	if _, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS operation_logs (
		id TEXT PRIMARY KEY,
		timestamp DATETIME NOT NULL,
		event_name TEXT NOT NULL,
		domain_type TEXT NOT NULL,
		aggregate_id TEXT NOT NULL,
		operator_id TEXT,
		action TEXT,
		details TEXT,
		metadata TEXT,
		created_at DATETIME NOT NULL
	)
	`); err != nil {
		return fmt.Errorf("failed to create operation_logs table: %w", err)
	}

	// 创建索引
	if _, err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_operation_logs_timestamp ON operation_logs(timestamp)`); err != nil {
		return fmt.Errorf("failed to create idx_operation_logs_timestamp index: %w", err)
	}

	if _, err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_operation_logs_event_name ON operation_logs(event_name)`); err != nil {
		return fmt.Errorf("failed to create idx_operation_logs_event_name index: %w", err)
	}

	if _, err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_operation_logs_domain_type ON operation_logs(domain_type)`); err != nil {
		return fmt.Errorf("failed to create idx_operation_logs_domain_type index: %w", err)
	}

	if _, err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_operation_logs_aggregate_id ON operation_logs(aggregate_id)`); err != nil {
		return fmt.Errorf("failed to create idx_operation_logs_aggregate_id index: %w", err)
	}

	if _, err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_operation_logs_operator_id ON operation_logs(operator_id)`); err != nil {
		return fmt.Errorf("failed to create idx_operation_logs_operator_id index: %w", err)
	}

	if _, err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_operation_logs_created_at ON operation_logs(created_at)`); err != nil {
		return fmt.Errorf("failed to create idx_operation_logs_created_at index: %w", err)
	}

	return nil
}

func (m *AddOperationLogsTableMigration) Down(tx *sql.Tx) error {
	// 删除索引
	if _, err := tx.Exec("DROP INDEX IF EXISTS idx_operation_logs_timestamp"); err != nil {
		return err
	}
	if _, err := tx.Exec("DROP INDEX IF EXISTS idx_operation_logs_event_name"); err != nil {
		return err
	}
	if _, err := tx.Exec("DROP INDEX IF EXISTS idx_operation_logs_domain_type"); err != nil {
		return err
	}
	if _, err := tx.Exec("DROP INDEX IF EXISTS idx_operation_logs_aggregate_id"); err != nil {
		return err
	}
	if _, err := tx.Exec("DROP INDEX IF EXISTS idx_operation_logs_operator_id"); err != nil {
		return err
	}
	if _, err := tx.Exec("DROP INDEX IF EXISTS idx_operation_logs_created_at"); err != nil {
		return err
	}

	// 删除表
	if _, err := tx.Exec("DROP TABLE IF EXISTS operation_logs"); err != nil {
		return err
	}

	return nil
}

// AddStatusAndNoteToRoomsMigration 为 rooms 表添加 status 和 note 字段
type AddStatusAndNoteToRoomsMigration struct{}

func (m *AddStatusAndNoteToRoomsMigration) Version() string {
	return "202603292200" // 格式：YYYYMMDDHHMM
}

func (m *AddStatusAndNoteToRoomsMigration) Up(tx *sql.Tx) error {
	// 检查列是否已存在
	rows, err := tx.Query("PRAGMA table_info(rooms)")
	if err != nil {
		return fmt.Errorf("failed to check table info: %w", err)
	}
	defer rows.Close()

	statusExists := false
	noteExists := false
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt_value sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt_value, &pk); err != nil {
			return fmt.Errorf("failed to scan table info: %w", err)
		}
		if name == "status" {
			statusExists = true
		}
		if name == "note" {
			noteExists = true
		}
	}

	if !statusExists {
		if _, err := tx.Exec("ALTER TABLE rooms ADD COLUMN status TEXT NOT NULL DEFAULT 'available'"); err != nil {
			return fmt.Errorf("failed to add status column: %w", err)
		}
	}

	if !noteExists {
		if _, err := tx.Exec("ALTER TABLE rooms ADD COLUMN note TEXT"); err != nil {
			return fmt.Errorf("failed to add note column: %w", err)
		}
	}

	return nil
}

func (m *AddStatusAndNoteToRoomsMigration) Down(tx *sql.Tx) error {
	// SQLite 不支持直接删除列，我们不处理回滚
	return nil
}

// AddDueDateToBillsMigration 为 bills 表添加 due_date 字段
type AddDueDateToBillsMigration struct{}

func (m *AddDueDateToBillsMigration) Version() string {
	return "202603292300" // 格式：YYYYMMDDHHMM
}

func (m *AddDueDateToBillsMigration) Up(tx *sql.Tx) error {
	// 检查 due_date 列是否已存在
	rows, err := tx.Query("PRAGMA table_info(bills)")
	if err != nil {
		return fmt.Errorf("failed to check table info: %w", err)
	}
	defer rows.Close()

	dueDateExists := false
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt_value sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt_value, &pk); err != nil {
			return fmt.Errorf("failed to scan table info: %w", err)
		}
		if name == "due_date" {
			dueDateExists = true
			break
		}
	}

	if !dueDateExists {
		if _, err := tx.Exec("ALTER TABLE bills ADD COLUMN due_date DATETIME"); err != nil {
			return fmt.Errorf("failed to add due_date column: %w", err)
		}
	}
	return nil
}

func (m *AddDueDateToBillsMigration) Down(tx *sql.Tx) error {
	// SQLite 不支持直接删除列，我们不处理回滚
	return nil
}

// AddRefundDepositAmountToBillsMigration 为 bills 表添加 refund_deposit_amount 字段
type AddRefundDepositAmountToBillsMigration struct{}

func (m *AddRefundDepositAmountToBillsMigration) Version() string {
	return "202603301200" // 格式：YYYYMMDDHHMM
}

func (m *AddRefundDepositAmountToBillsMigration) Up(tx *sql.Tx) error {
	// 检查列是否已存在
	rows, err := tx.Query("PRAGMA table_info(bills)")
	if err != nil {
		return fmt.Errorf("failed to check table info: %w", err)
	}
	defer rows.Close()

	colExists := false
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt_value sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt_value, &pk); err != nil {
			return fmt.Errorf("failed to scan table info: %w", err)
		}
		if name == "refund_deposit_amount" {
			colExists = true
			break
		}
	}

	if !colExists {
		if _, err := tx.Exec("ALTER TABLE bills ADD COLUMN refund_deposit_amount INTEGER NOT NULL DEFAULT 0"); err != nil {
			return fmt.Errorf("failed to add refund_deposit_amount column: %w", err)
		}
	}
	return nil
}

func (m *AddRefundDepositAmountToBillsMigration) Down(tx *sql.Tx) error {
	// SQLite 不支持直接删除列，我们不处理回滚
	return nil
}

// AddPrintJobsTableMigration 创建 print_jobs 表
type AddPrintJobsTableMigration struct{}

func (m *AddPrintJobsTableMigration) Version() string {
	return "202603311400" // 格式：YYYYMMDDHHMM
}

func (m *AddPrintJobsTableMigration) Up(tx *sql.Tx) error {
	// 创建 print_jobs 表
	if _, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS print_jobs (
		id TEXT PRIMARY KEY,
		type TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'pending',
		reference_id TEXT NOT NULL,
		tenant_name TEXT,
		tenant_phone TEXT,
		room_id TEXT,
		room_number TEXT,
		address TEXT,
		landlord_name TEXT,
		amount INTEGER NOT NULL DEFAULT 0,
		error_msg TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		completed_at DATETIME
	)
	`); err != nil {
		return fmt.Errorf("failed to create print_jobs table: %w", err)
	}

	// 创建索引
	if _, err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_print_jobs_status ON print_jobs(status)`); err != nil {
		return fmt.Errorf("failed to create idx_print_jobs_status index: %w", err)
	}

	if _, err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_print_jobs_type ON print_jobs(type)`); err != nil {
		return fmt.Errorf("failed to create idx_print_jobs_type index: %w", err)
	}

	if _, err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_print_jobs_created_at ON print_jobs(created_at)`); err != nil {
		return fmt.Errorf("failed to create idx_print_jobs_created_at index: %w", err)
	}

	if _, err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_print_jobs_reference_id ON print_jobs(reference_id)`); err != nil {
		return fmt.Errorf("failed to create idx_print_jobs_reference_id index: %w", err)
	}

	return nil
}

func (m *AddPrintJobsTableMigration) Down(tx *sql.Tx) error {
	// 删除索引
	if _, err := tx.Exec("DROP INDEX IF EXISTS idx_print_jobs_status"); err != nil {
		return err
	}
	if _, err := tx.Exec("DROP INDEX IF EXISTS idx_print_jobs_type"); err != nil {
		return err
	}
	if _, err := tx.Exec("DROP INDEX IF EXISTS idx_print_jobs_created_at"); err != nil {
		return err
	}
	if _, err := tx.Exec("DROP INDEX IF EXISTS idx_print_jobs_reference_id"); err != nil {
		return err
	}

	// 删除表
	if _, err := tx.Exec("DROP TABLE IF EXISTS print_jobs"); err != nil {
		return err
	}

	return nil
}

// AddPrintJobDetailsMigration 为 print_jobs 表添加详细信息字段
type AddPrintJobDetailsMigration struct{}

func (m *AddPrintJobDetailsMigration) Version() string {
	return "202603311500" // 格式：YYYYMMDDHHMM
}

func (m *AddPrintJobDetailsMigration) Up(tx *sql.Tx) error {
	// 检查列是否已存在
	rows, err := tx.Query("PRAGMA table_info(print_jobs)")
	if err != nil {
		return fmt.Errorf("failed to check table info: %w", err)
	}
	defer rows.Close()

	columns := make(map[string]bool)
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt_value sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt_value, &pk); err != nil {
			return fmt.Errorf("failed to scan table info: %w", err)
		}
		columns[name] = true
	}

	// 添加缺失的列
	if !columns["tenant_phone"] {
		if _, err := tx.Exec("ALTER TABLE print_jobs ADD COLUMN tenant_phone TEXT"); err != nil {
			return fmt.Errorf("failed to add tenant_phone column: %w", err)
		}
	}
	if !columns["room_id"] {
		if _, err := tx.Exec("ALTER TABLE print_jobs ADD COLUMN room_id TEXT"); err != nil {
			return fmt.Errorf("failed to add room_id column: %w", err)
		}
	}
	if !columns["room_number"] {
		if _, err := tx.Exec("ALTER TABLE print_jobs ADD COLUMN room_number TEXT"); err != nil {
			return fmt.Errorf("failed to add room_number column: %w", err)
		}
	}
	if !columns["address"] {
		if _, err := tx.Exec("ALTER TABLE print_jobs ADD COLUMN address TEXT"); err != nil {
			return fmt.Errorf("failed to add address column: %w", err)
		}
	}
	if !columns["landlord_name"] {
		if _, err := tx.Exec("ALTER TABLE print_jobs ADD COLUMN landlord_name TEXT"); err != nil {
			return fmt.Errorf("failed to add landlord_name column: %w", err)
		}
	}

	return nil
}

func (m *AddPrintJobDetailsMigration) Down(tx *sql.Tx) error {
	// SQLite 不支持直接删除列，我们不处理回滚
	return nil
}

// AddBillPeriodToBillsMigration 为 bills 表添加 bill_start 和 bill_end 字段
type AddBillPeriodToBillsMigration struct{}

func (m *AddBillPeriodToBillsMigration) Version() string {
	return "202603311600" // 格式：YYYYMMDDHHMM
}

func (m *AddBillPeriodToBillsMigration) Up(tx *sql.Tx) error {
	// 检查列是否已存在
	rows, err := tx.Query("PRAGMA table_info(bills)")
	if err != nil {
		return fmt.Errorf("failed to check table info: %w", err)
	}
	defer rows.Close()

	billStartExists := false
	billEndExists := false
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt_value sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt_value, &pk); err != nil {
			return fmt.Errorf("failed to scan table info: %w", err)
		}
		if name == "bill_start" {
			billStartExists = true
		}
		if name == "bill_end" {
			billEndExists = true
		}
	}

	if !billStartExists {
		if _, err := tx.Exec("ALTER TABLE bills ADD COLUMN bill_start DATETIME"); err != nil {
			return fmt.Errorf("failed to add bill_start column: %w", err)
		}
	}

	if !billEndExists {
		if _, err := tx.Exec("ALTER TABLE bills ADD COLUMN bill_end DATETIME"); err != nil {
			return fmt.Errorf("failed to add bill_end column: %w", err)
		}
	}

	return nil
}

func (m *AddBillPeriodToBillsMigration) Down(tx *sql.Tx) error {
	// SQLite 不支持直接删除列，我们不处理回滚
	return nil
}

// BaseMigration 基础数据模型迁移（创建所有表）
type BaseMigration struct{}

func (m *BaseMigration) Version() string {
	return "202603241430" // 格式：YYYYMMDDHHMM
}

func (m *BaseMigration) Up(tx *sql.Tx) error {
	// 创建 sagas 表
	if _, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS sagas (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		state TEXT NOT NULL,
		current_step INTEGER NOT NULL DEFAULT 0,
		error TEXT,
		data BLOB,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	)
	`); err != nil {
		return fmt.Errorf("failed to create sagas table: %w", err)
	}

	// 创建 locations 表
	if _, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS locations (
		id TEXT PRIMARY KEY,
		short_name TEXT NOT NULL,
		detail TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	)
	`); err != nil {
		return fmt.Errorf("failed to create locations table: %w", err)
	}

	// 创建 rooms 表 (with status and note from beginning)
	if _, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS rooms (
		id TEXT PRIMARY KEY,
		location_id TEXT NOT NULL,
		room_number TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'available',
		tags TEXT,
		note TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (location_id) REFERENCES locations(id),
		UNIQUE(location_id, room_number)
	)
	`); err != nil {
		return fmt.Errorf("failed to create rooms table: %w", err)
	}

	// 创建 landlords 表
	if _, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS landlords (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		phone TEXT,
		note TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	)
	`); err != nil {
		return fmt.Errorf("failed to create landlords table: %w", err)
	}

	// 创建 leases 表
	if _, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS leases (
		id TEXT PRIMARY KEY,
		room_id TEXT NOT NULL,
		landlord_id TEXT,
		tenant_name TEXT NOT NULL,
		tenant_phone TEXT,
		start_date DATETIME NOT NULL,
		end_date DATETIME NOT NULL,
		rent_amount INTEGER NOT NULL DEFAULT 0,
		deposit_amount INTEGER NOT NULL DEFAULT 0,
		status TEXT NOT NULL DEFAULT 'pending',
		note TEXT,
		last_charge_at DATETIME,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (room_id) REFERENCES rooms(id),
		FOREIGN KEY (landlord_id) REFERENCES landlords(id)
	)
	`); err != nil {
		return fmt.Errorf("failed to create leases table: %w", err)
	}

	// 创建 bills 表 (with due_date, refund_deposit_amount, bill_start, bill_end from beginning)
	if _, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS bills (
		id TEXT PRIMARY KEY,
		lease_id TEXT NOT NULL,
		type TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'pending',
		amount INTEGER NOT NULL DEFAULT 0,
		rent_amount INTEGER NOT NULL DEFAULT 0,
		water_amount INTEGER NOT NULL DEFAULT 0,
		electric_amount INTEGER NOT NULL DEFAULT 0,
		other_amount INTEGER NOT NULL DEFAULT 0,
		refund_deposit_amount INTEGER NOT NULL DEFAULT 0,
		bill_start DATETIME,
		bill_end DATETIME,
		due_date DATETIME,
		paid_at DATETIME,
		note TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (lease_id) REFERENCES leases(id)
	)
	`); err != nil {
		return fmt.Errorf("failed to create bills table: %w", err)
	}

	// 创建 deposits 表
	if _, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS deposits (
		id TEXT PRIMARY KEY,
		lease_id TEXT NOT NULL,
		amount INTEGER NOT NULL DEFAULT 0,
		status TEXT NOT NULL DEFAULT 'collected',
		refunded_at DATETIME,
		deducted_at DATETIME,
		note TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (lease_id) REFERENCES leases(id)
	)
	`); err != nil {
		return fmt.Errorf("failed to create deposits table: %w", err)
	}

	// 创建 indices 表
	if _, err := tx.Exec(`
	CREATE INDEX IF NOT EXISTS idx_leases_room_id ON leases(room_id)
	`); err != nil {
		return fmt.Errorf("failed to create idx_leases_room_id index: %w", err)
	}

	if _, err := tx.Exec(`
	CREATE INDEX IF NOT EXISTS idx_leases_room_id_status ON leases(room_id, status)
	`); err != nil {
		return fmt.Errorf("failed to create idx_leases_room_id_status index: %w", err)
	}

	if _, err := tx.Exec(`
	CREATE INDEX IF NOT EXISTS idx_leases_status ON leases(status)
	`); err != nil {
		return fmt.Errorf("failed to create idx_leases_status index: %w", err)
	}

	if _, err := tx.Exec(`
	CREATE INDEX IF NOT EXISTS idx_bills_lease_id ON bills(lease_id)
	`); err != nil {
		return fmt.Errorf("failed to create idx_bills_lease_id index: %w", err)
	}

	if _, err := tx.Exec(`
	CREATE INDEX IF NOT EXISTS idx_bills_paid_at ON bills(paid_at)
	`); err != nil {
		return fmt.Errorf("failed to create idx_bills_paid_at index: %w", err)
	}

	if _, err := tx.Exec(`
	CREATE INDEX IF NOT EXISTS idx_deposits_lease_id ON deposits(lease_id)
	`); err != nil {
		return fmt.Errorf("failed to create idx_deposits_lease_id index: %w", err)
	}

	return nil
}

func (m *BaseMigration) Down(tx *sql.Tx) error {
	// 删除索引
	if _, err := tx.Exec("DROP INDEX IF EXISTS idx_leases_room_id"); err != nil {
		return err
	}

	if _, err := tx.Exec("DROP INDEX IF EXISTS idx_leases_room_id_status"); err != nil {
		return err
	}

	if _, err := tx.Exec("DROP INDEX IF EXISTS idx_leases_status"); err != nil {
		return err
	}

	if _, err := tx.Exec("DROP INDEX IF EXISTS idx_bills_lease_id"); err != nil {
		return err
	}

	if _, err := tx.Exec("DROP INDEX IF EXISTS idx_bills_paid_at"); err != nil {
		return err
	}

	if _, err := tx.Exec("DROP INDEX IF EXISTS idx_deposits_lease_id"); err != nil {
		return err
	}

	// 删除表（按依赖关系的相反顺序删除）
	if _, err := tx.Exec("DROP TABLE IF EXISTS deposits"); err != nil {
		return err
	}

	if _, err := tx.Exec("DROP TABLE IF EXISTS bills"); err != nil {
		return err
	}

	if _, err := tx.Exec("DROP TABLE IF EXISTS leases"); err != nil {
		return err
	}

	if _, err := tx.Exec("DROP TABLE IF EXISTS landlords"); err != nil {
		return err
	}

	if _, err := tx.Exec("DROP TABLE IF EXISTS rooms"); err != nil {
		return err
	}

	if _, err := tx.Exec("DROP TABLE IF EXISTS locations"); err != nil {
		return err
	}

	return nil
}

// RunMigrations 运行迁移
func RunMigrations(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, m := range migrations {
		fmt.Printf("Running migration %s...\n", m.Version())
		if err := m.Up(tx); err != nil {
			return fmt.Errorf("migration %s failed: %w", m.Version(), err)
		}
		fmt.Printf("Migration %s completed.\n", m.Version())
	}

	return tx.Commit()
}

// RunRollback 回滚迁移
func RunRollback(db *sql.DB, targetVersion string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i := len(migrations) - 1; i >= 0; i-- {
		m := migrations[i]
		if m.Version() <= targetVersion {
			fmt.Printf("Rolling back migration %s...\n", m.Version())
			if err := m.Down(tx); err != nil {
				return fmt.Errorf("rollback migration %s failed: %w", m.Version(), err)
			}
			fmt.Printf("Rollback migration %s completed.\n", m.Version())
		}
	}

	return tx.Commit()
}
