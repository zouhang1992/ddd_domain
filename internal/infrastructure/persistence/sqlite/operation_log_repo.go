package sqlite

// OperationLogRepository is temporarily disabled
type OperationLogRepository struct {
	conn *Connection
}

// NewOperationLogRepository is temporarily disabled
func NewOperationLogRepository(conn *Connection) interface{} {
	return &OperationLogRepository{conn: conn}
}
