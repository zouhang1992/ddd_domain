package sqlite

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	operationlogmodel "github.com/zouhang1992/ddd_domain/internal/domain/operationlog/model"
	operationlogrepo "github.com/zouhang1992/ddd_domain/internal/domain/operationlog/repository"
)

// OperationLogRepository SQLite 操作日志仓储实现
type OperationLogRepository struct {
	conn *Connection
}

// NewOperationLogRepository 创建操作日志仓储
func NewOperationLogRepository(conn *Connection) operationlogrepo.OperationLogRepository {
	return &OperationLogRepository{conn: conn}
}

// tempOperationLog is a temporary struct for scanning
type tempOperationLog struct {
	ID          string
	Timestamp   time.Time
	EventName   string
	DomainType  string
	AggregateID string
	OperatorID  sql.NullString
	Action      sql.NullString
	Details     sql.NullString
	Metadata    sql.NullString
	CreatedAt   time.Time
}

// Save 保存操作日志
func (r *OperationLogRepository) Save(log *operationlogmodel.OperationLog) error {
	if log.ID == "" {
		log.ID = uuid.NewString()
	}

	// 检查是否已存在相同的日志（去重）
	// 通过 event_name, aggregate_id, timestamp 判断是否是同一条日志
	var exists bool
	err := r.conn.DB().QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM operation_logs
			WHERE event_name = ?
			AND aggregate_id = ?
			AND timestamp = ?
		)`,
		log.EventName, log.AggregateID, log.Timestamp).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		// 已存在相同日志，跳过保存
		return nil
	}

	detailsJSON, err := json.Marshal(log.Details)
	if err != nil {
		detailsJSON = []byte("{}")
	}

	metadataJSON, err := json.Marshal(log.Metadata)
	if err != nil {
		metadataJSON = []byte("{}")
	}

	_, err = r.conn.DB().Exec(`
		INSERT INTO operation_logs (
			id, timestamp, event_name, domain_type, aggregate_id,
			operator_id, action, details, metadata, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
		log.ID, log.Timestamp, log.EventName, log.DomainType, log.AggregateID,
		log.OperatorID, log.Action, string(detailsJSON), string(metadataJSON), log.CreatedAt)
	return err
}

// FindByID 根据ID查找操作日志
func (r *OperationLogRepository) FindByID(id string) (*operationlogmodel.OperationLog, error) {
	row := r.conn.DB().QueryRow(`
		SELECT id, timestamp, event_name, domain_type, aggregate_id,
			operator_id, action, details, metadata, created_at
		FROM operation_logs WHERE id = ?
		`, id)

	var temp tempOperationLog
	err := row.Scan(
		&temp.ID, &temp.Timestamp, &temp.EventName, &temp.DomainType, &temp.AggregateID,
		&temp.OperatorID, &temp.Action, &temp.Details, &temp.Metadata, &temp.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return tempToModel(&temp)
}

// FindByAggregateID 根据聚合ID查找操作日志
func (r *OperationLogRepository) FindByAggregateID(aggregateID string) ([]*operationlogmodel.OperationLog, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, timestamp, event_name, domain_type, aggregate_id,
			operator_id, action, details, metadata, created_at
		FROM operation_logs WHERE aggregate_id = ? ORDER BY timestamp DESC
		`, aggregateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*operationlogmodel.OperationLog
	for rows.Next() {
		var temp tempOperationLog
		err := rows.Scan(
			&temp.ID, &temp.Timestamp, &temp.EventName, &temp.DomainType, &temp.AggregateID,
			&temp.OperatorID, &temp.Action, &temp.Details, &temp.Metadata, &temp.CreatedAt)
		if err != nil {
			return nil, err
		}

		log, err := tempToModel(&temp)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

// FindByDomainType 根据领域类型分页查找操作日志
func (r *OperationLogRepository) FindByDomainType(domainType string, offset, limit int) ([]*operationlogmodel.OperationLog, int, error) {
	// 获取总数
	var total int
	row := r.conn.DB().QueryRow(`SELECT COUNT(*) FROM operation_logs WHERE domain_type = ?`, domainType)
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}

	// 分页查询
	rows, err := r.conn.DB().Query(`
		SELECT id, timestamp, event_name, domain_type, aggregate_id,
			operator_id, action, details, metadata, created_at
		FROM operation_logs WHERE domain_type = ? ORDER BY timestamp DESC LIMIT ? OFFSET ?
		`, domainType, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []*operationlogmodel.OperationLog
	for rows.Next() {
		var temp tempOperationLog
		err := rows.Scan(
			&temp.ID, &temp.Timestamp, &temp.EventName, &temp.DomainType, &temp.AggregateID,
			&temp.OperatorID, &temp.Action, &temp.Details, &temp.Metadata, &temp.CreatedAt)
		if err != nil {
			return nil, 0, err
		}

		log, err := tempToModel(&temp)
		if err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}
	return logs, total, nil
}

// FindByDomainTypeAndAggregateID 根据领域类型和聚合ID分页查找操作日志
func (r *OperationLogRepository) FindByDomainTypeAndAggregateID(domainType, aggregateID string, offset, limit int) ([]*operationlogmodel.OperationLog, int, error) {
	// 获取总数
	var total int
	row := r.conn.DB().QueryRow(`SELECT COUNT(*) FROM operation_logs WHERE domain_type = ? AND aggregate_id = ?`, domainType, aggregateID)
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}

	// 分页查询
	rows, err := r.conn.DB().Query(`
		SELECT id, timestamp, event_name, domain_type, aggregate_id,
			operator_id, action, details, metadata, created_at
		FROM operation_logs WHERE domain_type = ? AND aggregate_id = ? ORDER BY timestamp DESC LIMIT ? OFFSET ?
		`, domainType, aggregateID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []*operationlogmodel.OperationLog
	for rows.Next() {
		var temp tempOperationLog
		err := rows.Scan(
			&temp.ID, &temp.Timestamp, &temp.EventName, &temp.DomainType, &temp.AggregateID,
			&temp.OperatorID, &temp.Action, &temp.Details, &temp.Metadata, &temp.CreatedAt)
		if err != nil {
			return nil, 0, err
		}

		log, err := tempToModel(&temp)
		if err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}
	return logs, total, nil
}

// FindByTimeRange 根据时间范围分页查找操作日志
func (r *OperationLogRepository) FindByTimeRange(start, end time.Time, offset, limit int) ([]*operationlogmodel.OperationLog, int, error) {
	// 获取总数
	var total int
	row := r.conn.DB().QueryRow(`SELECT COUNT(*) FROM operation_logs WHERE timestamp >= ? AND timestamp <= ?`, start, end)
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}

	// 分页查询
	rows, err := r.conn.DB().Query(`
		SELECT id, timestamp, event_name, domain_type, aggregate_id,
			operator_id, action, details, metadata, created_at
		FROM operation_logs WHERE timestamp >= ? AND timestamp <= ? ORDER BY timestamp DESC LIMIT ? OFFSET ?
		`, start, end, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []*operationlogmodel.OperationLog
	for rows.Next() {
		var temp tempOperationLog
		err := rows.Scan(
			&temp.ID, &temp.Timestamp, &temp.EventName, &temp.DomainType, &temp.AggregateID,
			&temp.OperatorID, &temp.Action, &temp.Details, &temp.Metadata, &temp.CreatedAt)
		if err != nil {
			return nil, 0, err
		}

		log, err := tempToModel(&temp)
		if err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}
	return logs, total, nil
}

// tempToModel 将临时结构转换为领域模型
func tempToModel(temp *tempOperationLog) (*operationlogmodel.OperationLog, error) {
	var details map[string]any
	if temp.Details.Valid {
		if err := json.Unmarshal([]byte(temp.Details.String), &details); err != nil {
			details = make(map[string]any)
		}
	} else {
		details = make(map[string]any)
	}

	var metadata map[string]any
	if temp.Metadata.Valid {
		if err := json.Unmarshal([]byte(temp.Metadata.String), &metadata); err != nil {
			metadata = make(map[string]any)
		}
	} else {
		metadata = make(map[string]any)
	}

	operatorID := ""
	if temp.OperatorID.Valid {
		operatorID = temp.OperatorID.String
	}

	action := ""
	if temp.Action.Valid {
		action = temp.Action.String
	}

	return operationlogmodel.NewOperationLog(
		temp.ID,
		temp.Timestamp,
		temp.EventName,
		temp.DomainType,
		temp.AggregateID,
		operatorID,
		action,
		details,
		metadata,
	), nil
}
