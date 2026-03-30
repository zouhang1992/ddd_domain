package sqlite

import (
	"database/sql"
	_ "modernc.org/sqlite"

	"go.uber.org/zap"
)

// Config SQLite 配置
type Config struct {
	DSN string // Data Source Name, e.g., "data/ddd.db"
}

// Connection SQLite 连接
type Connection struct {
	db  *sql.DB
	log *zap.Logger
}

// NewConnection 创建 SQLite 连接
func NewConnection(cfg Config, logger *zap.Logger) (*Connection, error) {
	logger.Info("Opening database connection", zap.String("dsn", cfg.DSN))

	db, err := sql.Open("sqlite", cfg.DSN)
	if err != nil {
		logger.Error("Failed to open database", zap.Error(err))
		return nil, err
	}

	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		return nil, err
	}

	logger.Info("Database connection established successfully")

	conn := &Connection{db: db, log: logger}

	// Run migrations (this will create/update all tables)
	if err := RunMigrations(db); err != nil {
		logger.Error("Failed to run database migrations", zap.Error(err))
		return nil, err
	}

	logger.Info("Database schema initialized successfully")
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
