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
}

// AddDepositAmountToLeasesMigration 为 leases 表添加 deposit_amount 字段
type AddDepositAmountToLeasesMigration struct{}

func (m *AddDepositAmountToLeasesMigration) Version() string {
	return "202603251330" // 格式：YYYYMMDDHHMM
}

func (m *AddDepositAmountToLeasesMigration) Up(tx *sql.Tx) error {
	// 检查列是否已存在
	var colCount int
	row := tx.QueryRow("PRAGMA table_info(leases) WHERE name = 'deposit_amount'")
	if err := row.Scan(&colCount, &colCount, &colCount, &colCount, &colCount); err != nil && err != sql.ErrNoRows {
		// 列不存在，添加它
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

// BaseMigration 基础数据模型迁移（创建所有表）
type BaseMigration struct{}

func (m *BaseMigration) Version() string {
	return "202603241430" // 格式：YYYYMMDDHHMM
}

func (m *BaseMigration) Up(tx *sql.Tx) error {
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

	// 创建 rooms 表
	if _, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS rooms (
		id TEXT PRIMARY KEY,
		location_id TEXT NOT NULL,
		room_number TEXT NOT NULL,
		tags TEXT,
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

	// 创建 bills 表
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
