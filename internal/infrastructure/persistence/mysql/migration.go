package mysql

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
	&AddSessionsTableMigration{},
	&AddStatusAndNoteToRoomsMigration{},
	&AddDueDateToBillsMigration{},
	&AddRefundDepositAmountToBillsMigration{},
	&AddPrintJobsTableMigration{},
	&AddPrintJobDetailsMigration{},
	&AddBillPeriodToBillsMigration{},
}

// columnExists 检查 MySQL 列是否存在
func columnExists(tx *sql.Tx, tableName, columnName string) (bool, error) {
	var count int
	err := tx.QueryRow(`
		SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		AND TABLE_NAME = ?
		AND COLUMN_NAME = ?
	`, tableName, columnName).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check column exists: %w", err)
	}
	return count > 0, nil
}

// tableExists 检查 MySQL 表是否存在
func tableExists(tx *sql.Tx, tableName string) (bool, error) {
	var count int
	err := tx.QueryRow(`
		SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES
		WHERE TABLE_SCHEMA = DATABASE()
		AND TABLE_NAME = ?
	`, tableName).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check table exists: %w", err)
	}
	return count > 0, nil
}

// AddDepositAmountToLeasesMigration 为 leases 表添加 deposit_amount 字段
type AddDepositAmountToLeasesMigration struct{}

func (m *AddDepositAmountToLeasesMigration) Version() string {
	return "202603251330"
}

func (m *AddDepositAmountToLeasesMigration) Up(tx *sql.Tx) error {
	exists, err := columnExists(tx, "leases", "deposit_amount")
	if err != nil {
		return err
	}

	if !exists {
		if _, err := tx.Exec("ALTER TABLE leases ADD COLUMN deposit_amount INT NOT NULL DEFAULT 0"); err != nil {
			return fmt.Errorf("failed to add deposit_amount column: %w", err)
		}
	}
	return nil
}

func (m *AddDepositAmountToLeasesMigration) Down(tx *sql.Tx) error {
	exists, err := columnExists(tx, "leases", "deposit_amount")
	if err != nil {
		return err
	}
	if exists {
		if _, err := tx.Exec("ALTER TABLE leases DROP COLUMN deposit_amount"); err != nil {
			return err
		}
	}
	return nil
}

// AddOperationLogsTableMigration 创建 operation_logs 表
type AddOperationLogsTableMigration struct{}

func (m *AddOperationLogsTableMigration) Version() string {
	return "202603292100"
}

func (m *AddOperationLogsTableMigration) Up(tx *sql.Tx) error {
	exists, err := tableExists(tx, "operation_logs")
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	if _, err := tx.Exec(`
	CREATE TABLE operation_logs (
		id VARCHAR(36) PRIMARY KEY,
		timestamp DATETIME NOT NULL,
		event_name VARCHAR(255) NOT NULL,
		domain_type VARCHAR(255) NOT NULL,
		aggregate_id VARCHAR(36) NOT NULL,
		operator_id VARCHAR(36),
		action TEXT,
		details TEXT,
		metadata TEXT,
		created_at DATETIME NOT NULL,
		INDEX idx_operation_logs_timestamp (timestamp),
		INDEX idx_operation_logs_event_name (event_name),
		INDEX idx_operation_logs_domain_type (domain_type),
		INDEX idx_operation_logs_aggregate_id (aggregate_id),
		INDEX idx_operation_logs_operator_id (operator_id),
		INDEX idx_operation_logs_created_at (created_at)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`); err != nil {
		return fmt.Errorf("failed to create operation_logs table: %w", err)
	}

	return nil
}

func (m *AddOperationLogsTableMigration) Down(tx *sql.Tx) error {
	exists, err := tableExists(tx, "operation_logs")
	if err != nil {
		return err
	}
	if exists {
		if _, err := tx.Exec("DROP TABLE operation_logs"); err != nil {
			return err
		}
	}
	return nil
}

// AddSessionsTableMigration 创建 sessions 表
type AddSessionsTableMigration struct{}

func (m *AddSessionsTableMigration) Version() string {
	return "202604011200"
}

func (m *AddSessionsTableMigration) Up(tx *sql.Tx) error {
	exists, err := tableExists(tx, "sessions")
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	if _, err := tx.Exec(`
	CREATE TABLE sessions (
		id VARCHAR(36) PRIMARY KEY,
		user_id VARCHAR(36) NOT NULL,
		access_token TEXT NOT NULL,
		refresh_token TEXT,
		id_token TEXT NOT NULL,
		claims TEXT NOT NULL,
		expires_at DATETIME NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		INDEX idx_sessions_user_id (user_id),
		INDEX idx_sessions_expires_at (expires_at)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`); err != nil {
		return fmt.Errorf("failed to create sessions table: %w", err)
	}

	return nil
}

func (m *AddSessionsTableMigration) Down(tx *sql.Tx) error {
	exists, err := tableExists(tx, "sessions")
	if err != nil {
		return err
	}
	if exists {
		if _, err := tx.Exec("DROP TABLE sessions"); err != nil {
			return err
		}
	}
	return nil
}

// AddStatusAndNoteToRoomsMigration 为 rooms 表添加 status 和 note 字段
type AddStatusAndNoteToRoomsMigration struct{}

func (m *AddStatusAndNoteToRoomsMigration) Version() string {
	return "202603292200"
}

func (m *AddStatusAndNoteToRoomsMigration) Up(tx *sql.Tx) error {
	statusExists, err := columnExists(tx, "rooms", "status")
	if err != nil {
		return err
	}

	noteExists, err := columnExists(tx, "rooms", "note")
	if err != nil {
		return err
	}

	if !statusExists {
		if _, err := tx.Exec("ALTER TABLE rooms ADD COLUMN status VARCHAR(50) NOT NULL DEFAULT 'available'"); err != nil {
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
	// MySQL 支持删除列
	statusExists, err := columnExists(tx, "rooms", "status")
	if err != nil {
		return err
	}
	if statusExists {
		if _, err := tx.Exec("ALTER TABLE rooms DROP COLUMN status"); err != nil {
			return err
		}
	}

	noteExists, err := columnExists(tx, "rooms", "note")
	if err != nil {
		return err
	}
	if noteExists {
		if _, err := tx.Exec("ALTER TABLE rooms DROP COLUMN note"); err != nil {
			return err
		}
	}

	return nil
}

// AddDueDateToBillsMigration 为 bills 表添加 due_date 字段
type AddDueDateToBillsMigration struct{}

func (m *AddDueDateToBillsMigration) Version() string {
	return "202603292300"
}

func (m *AddDueDateToBillsMigration) Up(tx *sql.Tx) error {
	exists, err := columnExists(tx, "bills", "due_date")
	if err != nil {
		return err
	}

	if !exists {
		if _, err := tx.Exec("ALTER TABLE bills ADD COLUMN due_date DATETIME"); err != nil {
			return fmt.Errorf("failed to add due_date column: %w", err)
		}
	}
	return nil
}

func (m *AddDueDateToBillsMigration) Down(tx *sql.Tx) error {
	exists, err := columnExists(tx, "bills", "due_date")
	if err != nil {
		return err
	}
	if exists {
		if _, err := tx.Exec("ALTER TABLE bills DROP COLUMN due_date"); err != nil {
			return err
		}
	}
	return nil
}

// AddRefundDepositAmountToBillsMigration 为 bills 表添加 refund_deposit_amount 字段
type AddRefundDepositAmountToBillsMigration struct{}

func (m *AddRefundDepositAmountToBillsMigration) Version() string {
	return "202603301200"
}

func (m *AddRefundDepositAmountToBillsMigration) Up(tx *sql.Tx) error {
	exists, err := columnExists(tx, "bills", "refund_deposit_amount")
	if err != nil {
		return err
	}

	if !exists {
		if _, err := tx.Exec("ALTER TABLE bills ADD COLUMN refund_deposit_amount INT NOT NULL DEFAULT 0"); err != nil {
			return fmt.Errorf("failed to add refund_deposit_amount column: %w", err)
		}
	}
	return nil
}

func (m *AddRefundDepositAmountToBillsMigration) Down(tx *sql.Tx) error {
	exists, err := columnExists(tx, "bills", "refund_deposit_amount")
	if err != nil {
		return err
	}
	if exists {
		if _, err := tx.Exec("ALTER TABLE bills DROP COLUMN refund_deposit_amount"); err != nil {
			return err
		}
	}
	return nil
}

// AddPrintJobsTableMigration 创建 print_jobs 表
type AddPrintJobsTableMigration struct{}

func (m *AddPrintJobsTableMigration) Version() string {
	return "202603311400"
}

func (m *AddPrintJobsTableMigration) Up(tx *sql.Tx) error {
	exists, err := tableExists(tx, "print_jobs")
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	if _, err := tx.Exec(`
	CREATE TABLE print_jobs (
		id VARCHAR(36) PRIMARY KEY,
		type VARCHAR(50) NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'pending',
		reference_id VARCHAR(36) NOT NULL,
		tenant_name VARCHAR(255),
		tenant_phone VARCHAR(50),
		room_id VARCHAR(36),
		room_number VARCHAR(50),
		address TEXT,
		landlord_name VARCHAR(255),
		amount INT NOT NULL DEFAULT 0,
		error_msg TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		completed_at DATETIME,
		INDEX idx_print_jobs_status (status),
		INDEX idx_print_jobs_type (type),
		INDEX idx_print_jobs_created_at (created_at),
		INDEX idx_print_jobs_reference_id (reference_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`); err != nil {
		return fmt.Errorf("failed to create print_jobs table: %w", err)
	}

	return nil
}

func (m *AddPrintJobsTableMigration) Down(tx *sql.Tx) error {
	exists, err := tableExists(tx, "print_jobs")
	if err != nil {
		return err
	}
	if exists {
		if _, err := tx.Exec("DROP TABLE print_jobs"); err != nil {
			return err
		}
	}
	return nil
}

// AddPrintJobDetailsMigration 为 print_jobs 表添加详细信息字段
type AddPrintJobDetailsMigration struct{}

func (m *AddPrintJobDetailsMigration) Version() string {
	return "202603311500"
}

func (m *AddPrintJobDetailsMigration) Up(tx *sql.Tx) error {
	columns := []string{"tenant_phone", "room_id", "room_number", "address", "landlord_name"}
	for _, col := range columns {
		exists, err := columnExists(tx, "print_jobs", col)
		if err != nil {
			return err
		}
		if !exists {
			var colType string
			if col == "address" {
				colType = "TEXT"
			} else {
				colType = "VARCHAR(255)"
			}
			if _, err := tx.Exec(fmt.Sprintf("ALTER TABLE print_jobs ADD COLUMN %s %s", col, colType)); err != nil {
				return fmt.Errorf("failed to add %s column: %w", col, err)
			}
		}
	}
	return nil
}

func (m *AddPrintJobDetailsMigration) Down(tx *sql.Tx) error {
	// MySQL 支持删除列，但我们不处理回滚
	return nil
}

// AddBillPeriodToBillsMigration 为 bills 表添加 bill_start 和 bill_end 字段
type AddBillPeriodToBillsMigration struct{}

func (m *AddBillPeriodToBillsMigration) Version() string {
	return "202603311600"
}

func (m *AddBillPeriodToBillsMigration) Up(tx *sql.Tx) error {
	billStartExists, err := columnExists(tx, "bills", "bill_start")
	if err != nil {
		return err
	}

	billEndExists, err := columnExists(tx, "bills", "bill_end")
	if err != nil {
		return err
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
	// MySQL 支持删除列，但我们不处理回滚
	return nil
}

// BaseMigration 基础数据模型迁移（创建所有表）
type BaseMigration struct{}

func (m *BaseMigration) Version() string {
	return "202603241430"
}

func (m *BaseMigration) Up(tx *sql.Tx) error {
	// 创建 sagas 表
	sagasExists, err := tableExists(tx, "sagas")
	if err != nil {
		return err
	}
	if !sagasExists {
		if _, err := tx.Exec(`
		CREATE TABLE sagas (
			id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			state VARCHAR(50) NOT NULL,
			current_step INT NOT NULL DEFAULT 0,
			error TEXT,
			data JSON,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
		`); err != nil {
			return fmt.Errorf("failed to create sagas table: %w", err)
		}
	}

	// 创建 locations 表
	locationsExists, err := tableExists(tx, "locations")
	if err != nil {
		return err
	}
	if !locationsExists {
		if _, err := tx.Exec(`
		CREATE TABLE locations (
			id VARCHAR(36) PRIMARY KEY,
			short_name VARCHAR(255) NOT NULL,
			detail TEXT,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
		`); err != nil {
			return fmt.Errorf("failed to create locations table: %w", err)
		}
	}

	// 创建 rooms 表
	roomsExists, err := tableExists(tx, "rooms")
	if err != nil {
		return err
	}
	if !roomsExists {
		if _, err := tx.Exec(`
		CREATE TABLE rooms (
			id VARCHAR(36) PRIMARY KEY,
			location_id VARCHAR(36) NOT NULL,
			room_number VARCHAR(50) NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'available',
			tags TEXT,
			note TEXT,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			FOREIGN KEY (location_id) REFERENCES locations(id),
			UNIQUE KEY uk_location_room (location_id, room_number)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
		`); err != nil {
			return fmt.Errorf("failed to create rooms table: %w", err)
		}
	}

	// 创建 landlords 表
	landlordsExists, err := tableExists(tx, "landlords")
	if err != nil {
		return err
	}
	if !landlordsExists {
		if _, err := tx.Exec(`
		CREATE TABLE landlords (
			id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			phone VARCHAR(50),
			note TEXT,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
		`); err != nil {
			return fmt.Errorf("failed to create landlords table: %w", err)
		}
	}

	// 创建 leases 表
	leasesExists, err := tableExists(tx, "leases")
	if err != nil {
		return err
	}
	if !leasesExists {
		if _, err := tx.Exec(`
		CREATE TABLE leases (
			id VARCHAR(36) PRIMARY KEY,
			room_id VARCHAR(36) NOT NULL,
			landlord_id VARCHAR(36),
			tenant_name VARCHAR(255) NOT NULL,
			tenant_phone VARCHAR(50),
			start_date DATETIME NOT NULL,
			end_date DATETIME NOT NULL,
			rent_amount INT NOT NULL DEFAULT 0,
			deposit_amount INT NOT NULL DEFAULT 0,
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			note TEXT,
			last_charge_at DATETIME,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			FOREIGN KEY (room_id) REFERENCES rooms(id),
			FOREIGN KEY (landlord_id) REFERENCES landlords(id),
			INDEX idx_leases_room_id (room_id),
			INDEX idx_leases_room_id_status (room_id, status),
			INDEX idx_leases_status (status)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
		`); err != nil {
			return fmt.Errorf("failed to create leases table: %w", err)
		}
	}

	// 创建 bills 表
	billsExists, err := tableExists(tx, "bills")
	if err != nil {
		return err
	}
	if !billsExists {
		if _, err := tx.Exec(`
		CREATE TABLE bills (
			id VARCHAR(36) PRIMARY KEY,
			lease_id VARCHAR(36) NOT NULL,
			type VARCHAR(50) NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			amount INT NOT NULL DEFAULT 0,
			rent_amount INT NOT NULL DEFAULT 0,
			water_amount INT NOT NULL DEFAULT 0,
			electric_amount INT NOT NULL DEFAULT 0,
			other_amount INT NOT NULL DEFAULT 0,
			refund_deposit_amount INT NOT NULL DEFAULT 0,
			bill_start DATETIME,
			bill_end DATETIME,
			due_date DATETIME,
			paid_at DATETIME,
			note TEXT,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			FOREIGN KEY (lease_id) REFERENCES leases(id),
			INDEX idx_bills_lease_id (lease_id),
			INDEX idx_bills_paid_at (paid_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
		`); err != nil {
			return fmt.Errorf("failed to create bills table: %w", err)
		}
	}

	// 创建 deposits 表
	depositsExists, err := tableExists(tx, "deposits")
	if err != nil {
		return err
	}
	if !depositsExists {
		if _, err := tx.Exec(`
		CREATE TABLE deposits (
			id VARCHAR(36) PRIMARY KEY,
			lease_id VARCHAR(36) NOT NULL,
			amount INT NOT NULL DEFAULT 0,
			status VARCHAR(50) NOT NULL DEFAULT 'collected',
			refunded_at DATETIME,
			deducted_at DATETIME,
			note TEXT,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			FOREIGN KEY (lease_id) REFERENCES leases(id),
			INDEX idx_deposits_lease_id (lease_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
		`); err != nil {
			return fmt.Errorf("failed to create deposits table: %w", err)
		}
	}

	return nil
}

func (m *BaseMigration) Down(tx *sql.Tx) error {
	// 删除表（按依赖关系的相反顺序删除）
	tables := []string{"deposits", "bills", "leases", "landlords", "rooms", "locations", "sagas"}
	for _, table := range tables {
		exists, err := tableExists(tx, table)
		if err != nil {
			return err
		}
		if exists {
			if _, err := tx.Exec(fmt.Sprintf("DROP TABLE %s", table)); err != nil {
				return err
			}
		}
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
