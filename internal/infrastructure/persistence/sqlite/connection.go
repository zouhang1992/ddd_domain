package sqlite

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

// Config SQLite 配置
type Config struct {
	DSN string // Data Source Name, e.g., "data/ddd.db"
}

// Connection SQLite 连接
type Connection struct {
	db *sql.DB
}

// NewConnection 创建 SQLite 连接
func NewConnection(cfg Config) (*Connection, error) {
	db, err := sql.Open("sqlite", cfg.DSN)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	conn := &Connection{db: db}
	if err := conn.initSchema(); err != nil {
		return nil, err
	}

	return conn, nil
}

// DB 获取原始数据库连接
func (c *Connection) DB() *sql.DB {
	return c.db
}

// Close 关闭连接
func (c *Connection) Close() error {
	return c.db.Close()
}

// BeginTx 开始事务
func (c *Connection) BeginTx() (*Tx, error) {
	tx, err := c.db.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx: tx}, nil
}

// initSchema 初始化数据库表
func (c *Connection) initSchema() error {
	schema := `
CREATE TABLE IF NOT EXISTS sagas (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	state TEXT NOT NULL,
	current_step INTEGER NOT NULL DEFAULT 0,
	error TEXT,
	data BLOB,
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS locations (
	id TEXT PRIMARY KEY,
	short_name TEXT NOT NULL,
	detail TEXT,
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS rooms (
	id TEXT PRIMARY KEY,
	location_id TEXT NOT NULL,
	room_number TEXT NOT NULL,
	tags TEXT,
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL,
	FOREIGN KEY (location_id) REFERENCES locations(id),
	UNIQUE(location_id, room_number)
);

CREATE TABLE IF NOT EXISTS landlords (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	phone TEXT,
	note TEXT,
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL
);

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
);

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
);

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
);

CREATE TABLE IF NOT EXISTS operation_logs (
	id TEXT PRIMARY KEY,
	timestamp DATETIME NOT NULL,
	event_name TEXT NOT NULL,
	domain_type TEXT NOT NULL,
	aggregate_id TEXT,
	operator_id TEXT,
	action TEXT NOT NULL,
	details TEXT,
	metadata TEXT,
	created_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_operation_logs_timestamp ON operation_logs(timestamp);
CREATE INDEX IF NOT EXISTS idx_operation_logs_event_name ON operation_logs(event_name);
CREATE INDEX IF NOT EXISTS idx_operation_logs_domain_type ON operation_logs(domain_type);
CREATE INDEX IF NOT EXISTS idx_operation_logs_aggregate_id ON operation_logs(aggregate_id);
CREATE INDEX IF NOT EXISTS idx_operation_logs_operator_id ON operation_logs(operator_id);
CREATE INDEX IF NOT EXISTS idx_operation_logs_created_at ON operation_logs(created_at);
`
	_, err := c.db.Exec(schema)
	return err
}

// Tx 事务包装
type Tx struct {
	tx *sql.Tx
}

// Commit 提交事务
func (t *Tx) Commit() error {
	return t.tx.Commit()
}

// Rollback 回滚事务
func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}

// Exec 执行 SQL
func (t *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.tx.Exec(query, args...)
}

// Query 查询 SQL
func (t *Tx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.Query(query, args...)
}

// QueryRow 查询单行
func (t *Tx) QueryRow(query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRow(query, args...)
}
