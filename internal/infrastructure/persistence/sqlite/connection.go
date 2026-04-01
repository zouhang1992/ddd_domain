package sqlite

import (
	"database/sql"
	_ "modernc.org/sqlite"

	"go.uber.org/zap"

	"github.com/zouhang1992/ddd_domain/internal/application/config"
)

// Config SQLite 配置（别名以保持向后兼容）
type Config = config.DatabaseConfig

// Connection SQLite 连接
type Connection struct {
	db  *sql.DB
	log *zap.Logger
}

// NewConnection 创建 SQLite 连接
func NewConnection(cfg config.DatabaseConfig, logger *zap.Logger) (*Connection, error) {
	logger.Info("Opening database connection", zap.String("dsn", cfg.DSN))

	// Add SQLite 配置详解：
	// _journal=WAL: 使用 Write-Ahead Logging 模式，提高并发性能
	// _busy_timeout=5000: 设置 5 秒的锁等待超时
	// _txlock=immediate: 立即获取写锁，减少锁竞争
	dsnWithParams := cfg.DSN + "?_journal=WAL&_busy_timeout=5000&_txlock=immediate"

	db, err := sql.Open("sqlite", dsnWithParams)
	if err != nil {
		logger.Error("Failed to open database", zap.Error(err))
		return nil, err
	}

	// 设置连接池参数以避免连接泄漏
	db.SetMaxOpenConns(1) // SQLite 只需要一个写连接
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0) // 连接不过期

	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		return nil, err
	}

	// 启用 WAL 模式并设置同步模式
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		logger.Warn("Failed to set WAL mode", zap.Error(err))
		// 继续执行，WAL 不是必需的
	}

	// 设置同步模式为 NORMAL，在数据安全和性能之间取得平衡
	_, err = db.Exec("PRAGMA synchronous=NORMAL;")
	if err != nil {
		logger.Warn("Failed to set synchronous mode", zap.Error(err))
	}

	// 设置缓存大小
	_, err = db.Exec("PRAGMA cache_size=-64000;") // ~64MB cache
	if err != nil {
		logger.Warn("Failed to set cache size", zap.Error(err))
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
