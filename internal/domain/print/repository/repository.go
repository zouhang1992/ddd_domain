package repository

import (
	"time"

	printmodel "github.com/zouhang1992/ddd_domain/internal/domain/print/model"
)

// PrintJobRepository 打印作业仓储接口
type PrintJobRepository interface {
	Save(job *printmodel.PrintJob) error
	FindByID(id string) (*printmodel.PrintJob, error)
	FindAll(offset, limit int) ([]*printmodel.PrintJob, int, error)
	FindByStatus(status printmodel.PrintJobStatus, offset, limit int) ([]*printmodel.PrintJob, int, error)
	FindByType(jobType printmodel.PrintJobType, offset, limit int) ([]*printmodel.PrintJob, int, error)
	FindByTimeRange(start, end time.Time, offset, limit int) ([]*printmodel.PrintJob, int, error)
	FindByFilters(status printmodel.PrintJobStatus, jobType printmodel.PrintJobType, start, end *time.Time, offset, limit int) ([]*printmodel.PrintJob, int, error)
	Delete(id string) error
}
