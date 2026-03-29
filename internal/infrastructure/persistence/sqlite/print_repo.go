package sqlite

// PrintJobRepository is temporarily disabled
type PrintJobRepository struct {
	conn *Connection
}

// NewPrintJobRepository is temporarily disabled
func NewPrintJobRepository(conn *Connection) interface{} {
	return &PrintJobRepository{conn: conn}
}
