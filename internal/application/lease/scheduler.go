package lease

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	leaseservice "github.com/zouhang1992/ddd_domain/internal/domain/lease/service"
)

// LeaseExpirationScheduler 租约到期定时调度器
type LeaseExpirationScheduler struct {
	leaseService *leaseservice.LeaseService
	logger       *zap.Logger
	cron         *cron.Cron
	running      bool
}

// NewLeaseExpirationScheduler 创建租约到期定时调度器
func NewLeaseExpirationScheduler(leaseService *leaseservice.LeaseService, logger *zap.Logger) *LeaseExpirationScheduler {
	return &LeaseExpirationScheduler{
		leaseService: leaseService,
		logger:       logger,
	}
}

// Start 启动调度器
func (s *LeaseExpirationScheduler) Start() error {
	if s.running {
		s.logger.Warn("Lease expiration scheduler already running")
		return nil
	}

	s.logger.Info("Starting lease expiration scheduler")

	// 立即执行一次
	s.runCheck()

	// 创建 cron 调度器
	s.cron = cron.New()

	// 每小时执行一次（整点）
	_, err := s.cron.AddFunc("0 * * * *", func() {
		s.runCheck()
	})
	if err != nil {
		s.logger.Error("Failed to schedule lease expiration check", zap.Error(err))
		return err
	}

	// 启动 cron 调度器
	s.cron.Start()
	s.running = true

	s.logger.Info("Lease expiration scheduler started, will run hourly")
	return nil
}

// Stop 停止调度器
func (s *LeaseExpirationScheduler) Stop() {
	if !s.running {
		return
	}

	s.logger.Info("Stopping lease expiration scheduler")

	if s.cron != nil {
		ctx := s.cron.Stop()
		<-ctx.Done()
	}

	s.running = false
	s.logger.Info("Lease expiration scheduler stopped")
}

// runCheck 执行租约到期检查
func (s *LeaseExpirationScheduler) runCheck() {
	s.logger.Info("Running lease expiration check")

	count, err := s.leaseService.CheckAndExpireLeases()
	if err != nil {
		s.logger.Error("Failed to check and expire leases", zap.Error(err))
		return
	}

	s.logger.Info("Lease expiration check completed",
		zap.Int("processed_count", count))
}
