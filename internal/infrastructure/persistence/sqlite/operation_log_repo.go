package sqlite

import (
	"database/sql"
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
)

// OperationLogRepository SQLite 操作日志仓储实现
type OperationLogRepository struct {
	conn *Connection
}

// NewOperationLogRepository 创建操作日志仓储
func NewOperationLogRepository(conn *Connection) repository.OperationLogRepository {
	return &OperationLogRepository{conn: conn}
}

// Save 保存操作日志
func (r *OperationLogRepository) Save(log *model.OperationLog) error {
	detailsJSON, err := log.MarshalDetails()
	if err != nil {
		return err
	}

	metadataJSON, err := log.MarshalMetadata()
	if err != nil {
		return err
	}

	_, err = r.conn.DB().Exec(`
		INSERT OR REPLACE INTO operation_logs (
			id, timestamp, event_name, domain_type, aggregate_id,
			operator_id, action, details, metadata, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, log.ID(), log.Timestamp(), log.EventName(), log.DomainType(), log.AggregateID(),
		log.OperatorID(), log.Action(), detailsJSON, metadataJSON, log.CreatedAt())
	return err
}

// FindByID 根据ID查找操作日志
func (r *OperationLogRepository) FindByID(id string) (*model.OperationLog, error) {
	row := r.conn.DB().QueryRow(`
		SELECT id, timestamp, event_name, domain_type, aggregate_id,
			operator_id, action, details, metadata, created_at
		FROM operation_logs WHERE id = ?
	`, id)

	var idStr, eventName, domainType, aggregateID, operatorID, action string
	var timestamp, createdAt time.Time
	var detailsJSON, metadataJSON sql.NullString

	err := row.Scan(&idStr, &timestamp, &eventName, &domainType, &aggregateID,
		&operatorID, &action, &detailsJSON, &metadataJSON, &createdAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	log := model.NewOperationLog(idStr, timestamp, eventName, domainType, aggregateID,
		operatorID, action, make(map[string]interface{}), make(map[string]interface{}))

	if detailsJSON.Valid {
		if err := log.UnmarshalDetails(detailsJSON.String); err != nil {
			return nil, err
		}
	}

	if metadataJSON.Valid {
		if err := log.UnmarshalMetadata(metadataJSON.String); err != nil {
			return nil, err
		}
	}

	return log, nil
}

// FindAll 查找所有操作日志
func (r *OperationLogRepository) FindAll(offset, limit int) ([]*model.OperationLog, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, timestamp, event_name, domain_type, aggregate_id,
			operator_id, action, details, metadata, created_at
		FROM operation_logs ORDER BY timestamp DESC LIMIT ? OFFSET ?
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

// FindByDomainType 按领域类型查找操作日志
func (r *OperationLogRepository) FindByDomainType(domainType string, offset, limit int) ([]*model.OperationLog, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, timestamp, event_name, domain_type, aggregate_id,
			operator_id, action, details, metadata, created_at
		FROM operation_logs WHERE domain_type = ? ORDER BY timestamp DESC LIMIT ? OFFSET ?
	`, domainType, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

// FindByTimeRange 按时间范围查找操作日志
func (r *OperationLogRepository) FindByTimeRange(startTime, endTime time.Time, offset, limit int) ([]*model.OperationLog, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, timestamp, event_name, domain_type, aggregate_id,
			operator_id, action, details, metadata, created_at
		FROM operation_logs WHERE timestamp BETWEEN ? AND ?
		ORDER BY timestamp DESC LIMIT ? OFFSET ?
	`, startTime, endTime, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

// FindByAggregateID 按聚合ID查找操作日志
func (r *OperationLogRepository) FindByAggregateID(aggregateID string, offset, limit int) ([]*model.OperationLog, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, timestamp, event_name, domain_type, aggregate_id,
			operator_id, action, details, metadata, created_at
		FROM operation_logs WHERE aggregate_id = ? ORDER BY timestamp DESC LIMIT ? OFFSET ?
	`, aggregateID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

// FindByOperatorID 按操作人ID查找操作日志
func (r *OperationLogRepository) FindByOperatorID(operatorID string, offset, limit int) ([]*model.OperationLog, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, timestamp, event_name, domain_type, aggregate_id,
			operator_id, action, details, metadata, created_at
		FROM operation_logs WHERE operator_id = ? ORDER BY timestamp DESC LIMIT ? OFFSET ?
	`, operatorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

// FindByCriteria 按条件查找操作日志
func (r *OperationLogRepository) FindByCriteria(criteria repository.OperationLogCriteria, offset, limit int) ([]*model.OperationLog, error) {
	query := `
		SELECT id, timestamp, event_name, domain_type, aggregate_id,
			operator_id, action, details, metadata, created_at
		FROM operation_logs
		WHERE 1 = 1
	`
	var args []interface{}

	if criteria.DomainType != "" {
		query += " AND domain_type = ?"
		args = append(args, criteria.DomainType)
	}
	if criteria.EventName != "" {
		query += " AND event_name LIKE ?"
		args = append(args, "%"+criteria.EventName+"%")
	}
	if criteria.AggregateID != "" {
		query += " AND aggregate_id = ?"
		args = append(args, criteria.AggregateID)
	}
	if criteria.OperatorID != "" {
		query += " AND operator_id = ?"
		args = append(args, criteria.OperatorID)
	}
	if criteria.StartTime != nil {
		query += " AND timestamp >= ?"
		args = append(args, criteria.StartTime)
	}
	if criteria.EndTime != nil {
		query += " AND timestamp <= ?"
		args = append(args, criteria.EndTime)
	}

	query += " ORDER BY timestamp DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.conn.DB().Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

// CountByCriteria 按条件统计操作日志数量
func (r *OperationLogRepository) CountByCriteria(criteria repository.OperationLogCriteria) (int, error) {
	query := `
		SELECT COUNT(*) FROM operation_logs
		WHERE 1 = 1
	`
	var args []interface{}

	if criteria.DomainType != "" {
		query += " AND domain_type = ?"
		args = append(args, criteria.DomainType)
	}
	if criteria.EventName != "" {
		query += " AND event_name LIKE ?"
		args = append(args, "%"+criteria.EventName+"%")
	}
	if criteria.AggregateID != "" {
		query += " AND aggregate_id = ?"
		args = append(args, criteria.AggregateID)
	}
	if criteria.OperatorID != "" {
		query += " AND operator_id = ?"
		args = append(args, criteria.OperatorID)
	}
	if criteria.StartTime != nil {
		query += " AND timestamp >= ?"
		args = append(args, criteria.StartTime)
	}
	if criteria.EndTime != nil {
		query += " AND timestamp <= ?"
		args = append(args, criteria.EndTime)
	}

	var count int
	row := r.conn.DB().QueryRow(query, args...)
	err := row.Scan(&count)
	return count, err
}

// scanLogs 扫描查询结果并转换为 model.OperationLog 切片
func (r *OperationLogRepository) scanLogs(rows *sql.Rows) ([]*model.OperationLog, error) {
	var logs []*model.OperationLog

	for rows.Next() {
		var idStr, eventName, domainType, aggregateID, operatorID, action string
		var timestamp, createdAt time.Time
		var detailsJSON, metadataJSON sql.NullString

		err := rows.Scan(&idStr, &timestamp, &eventName, &domainType, &aggregateID,
			&operatorID, &action, &detailsJSON, &metadataJSON, &createdAt)
		if err != nil {
			return nil, err
		}

		log := model.NewOperationLog(idStr, timestamp, eventName, domainType, aggregateID,
			operatorID, action, make(map[string]interface{}), make(map[string]interface{}))

		if detailsJSON.Valid {
			if err := log.UnmarshalDetails(detailsJSON.String); err != nil {
				return nil, err
			}
		}

		if metadataJSON.Valid {
			if err := log.UnmarshalMetadata(metadataJSON.String); err != nil {
				return nil, err
			}
		}

		logs = append(logs, log)
	}

	return logs, nil
}
