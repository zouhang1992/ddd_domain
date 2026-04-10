package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	leasemodel "github.com/zouhang1992/ddd_domain/internal/domain/lease/model"
	leaserepo "github.com/zouhang1992/ddd_domain/internal/domain/lease/repository"
	depositmodel "github.com/zouhang1992/ddd_domain/internal/domain/deposit/model"
	depositrepo "github.com/zouhang1992/ddd_domain/internal/domain/deposit/repository"
	roomrepo "github.com/zouhang1992/ddd_domain/internal/domain/room/repository"
	roommodel "github.com/zouhang1992/ddd_domain/internal/domain/room/model"
	domerrors "github.com/zouhang1992/ddd_domain/internal/domain/common/errors"
)

// LeaseService 租约领域服务
type LeaseService struct {
	leaseRepo   leaserepo.LeaseRepository
	depositRepo depositrepo.DepositRepository
	roomRepo    roomrepo.RoomRepository
}

// NewLeaseService 创建租约领域服务
func NewLeaseService(leaseRepo leaserepo.LeaseRepository, depositRepo depositrepo.DepositRepository, roomRepo roomrepo.RoomRepository) *LeaseService {
	return &LeaseService{
		leaseRepo:   leaseRepo,
		depositRepo: depositRepo,
		roomRepo:    roomRepo,
	}
}

// ValidateRoomForLease 校验房间是否可用于租约
func (s *LeaseService) ValidateRoomForLease(room *roommodel.Room) error {
	if room == nil {
		return domerrors.ErrRoomNotFound
	}

	if room.Status != roommodel.RoomStatusAvailable {
		return domerrors.ErrRoomNotAvailable
	}

	return nil
}

// CreateLeaseResult 创建租约结果
type CreateLeaseResult struct {
	Lease   *leasemodel.Lease
	Deposit *depositmodel.Deposit
}

// CreateLease 创建租约（含押金）
func (s *LeaseService) CreateLease(roomID, landlordID, tenantName, tenantPhone string,
	startDate, endDate time.Time, rentAmount, depositAmount int64, note, depositNote string) (*CreateLeaseResult, error) {

	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		return nil, err
	}
	if err := s.ValidateRoomForLease(room); err != nil {
		return nil, err
	}

	// 检查房间是否已有活跃的租约（pending 或 active 状态）
	hasActiveLease, err := s.leaseRepo.HasActiveLeaseForRoom(roomID)
	if err != nil {
		return nil, err
	}
	if hasActiveLease {
		return nil, domerrors.ErrRoomNotAvailable
	}

	leaseID := uuid.NewString()
	lease := leasemodel.NewLease(leaseID, roomID, landlordID, tenantName, tenantPhone, startDate, endDate, rentAmount, depositAmount, note)

	var deposit *depositmodel.Deposit
	if depositAmount > 0 {
		depositID := uuid.NewString()
		deposit = depositmodel.NewDeposit(depositID, leaseID, depositAmount, depositNote)
	}

	return &CreateLeaseResult{
		Lease:   lease,
		Deposit: deposit,
	}, nil
}

// ValidateDelete 校验租约是否可删除
func (s *LeaseService) ValidateDelete(leaseID string) error {
	hasBills, err := s.leaseRepo.HasBills(leaseID)
	if err != nil {
		return err
	}
	if hasBills {
		return domerrors.ErrCannotDelete
	}

	hasDeposit, err := s.leaseRepo.HasDeposit(leaseID)
	if err != nil {
		return err
	}
	if hasDeposit {
		return domerrors.ErrCannotDelete
	}

	return nil
}

// ValidateActivate 校验租约是否可激活
func (s *LeaseService) ValidateActivate(lease *leasemodel.Lease, room *roommodel.Room) error {
	if lease.Status != leasemodel.LeaseStatusPending {
		return errors.New("租约状态不是待生效，当前状态: " + string(lease.Status))
	}

	// 使用 Asia/Shanghai 时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		// 如果加载时区失败，使用本地时间
		loc = time.Local
	}

	// 允许激活开始日期在今天或之前的租约（使用 Asia/Shanghai 时区）
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	leaseStartDate := time.Date(lease.StartDate.Year(), lease.StartDate.Month(), lease.StartDate.Day(), 0, 0, 0, 0, loc)

	if leaseStartDate.After(today) {
		return errors.New("租约开始日期晚于今天，开始日期: " + lease.StartDate.Format("2006-01-02") + ", 今天: " + today.Format("2006-01-02"))
	}

	if err := s.ValidateRoomForLease(room); err != nil {
		return err
	}

	return nil
}

// CheckAndExpireLeases 检查并处理到期租约
func (s *LeaseService) CheckAndExpireLeases() (int, error) {
	now := time.Now()

	leases, err := s.leaseRepo.FindActiveLeasesExpiringBefore(now)
	if err != nil {
		return 0, err
	}

	processed := 0
	for _, lease := range leases {
		if lease.Status != leasemodel.LeaseStatusActive {
			continue
		}

		lease.Expire()

		if err := s.leaseRepo.Save(lease); err != nil {
			continue
		}

		processed++
	}

	return processed, nil
}
