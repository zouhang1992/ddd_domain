package model

import (
	"testing"
	"time"
)

func TestNewLease(t *testing.T) {
	startDate, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	endDate, _ := time.Parse(time.RFC3339, "2023-12-31T23:59:59Z")

	lease := NewLease("123", "room-1", "landlord-1", "张三", "13812345678",
		startDate, endDate, 10000, 20000, "测试租约")

	if lease.ID != "123" {
		t.Errorf("Expected ID '123', got '%s'", lease.ID)
	}

	if lease.RoomID != "room-1" {
		t.Errorf("Expected RoomID 'room-1', got '%s'", lease.RoomID)
	}

	if lease.LandlordID != "landlord-1" {
		t.Errorf("Expected LandlordID 'landlord-1', got '%s'", lease.LandlordID)
	}

	if lease.TenantName != "张三" {
		t.Errorf("Expected TenantName '张三', got '%s'", lease.TenantName)
	}

	if lease.TenantPhone != "13812345678" {
		t.Errorf("Expected TenantPhone '13812345678', got '%s'", lease.TenantPhone)
	}

	if lease.StartDate != startDate {
		t.Errorf("Expected StartDate %v, got %v", startDate, lease.StartDate)
	}

	if lease.EndDate != endDate {
		t.Errorf("Expected EndDate %v, got %v", endDate, lease.EndDate)
	}

	if lease.RentAmount != 10000 {
		t.Errorf("Expected RentAmount 10000, got %d", lease.RentAmount)
	}

	if lease.Status != "pending" {
		t.Errorf("Expected Status 'pending', got '%s'", lease.Status)
	}

	if lease.Note != "测试租约" {
		t.Errorf("Expected Note '测试租约', got '%s'", lease.Note)
	}

	// 检查创建时间
	if lease.CreatedAt.IsZero() || lease.UpdatedAt.IsZero() {
		t.Error("Expected CreatedAt and UpdatedAt to be set")
	}
}

func TestActivate(t *testing.T) {
	startDate, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	endDate, _ := time.Parse(time.RFC3339, "2023-12-31T23:59:59Z")

	lease := NewLease("123", "room-1", "landlord-1", "张三", "13812345678",
		startDate, endDate, 10000, 20000, "测试租约")
	originalUpdatedAt := lease.UpdatedAt

	time.Sleep(10 * time.Millisecond)
	lease.Activate()

	if lease.Status != "active" {
		t.Errorf("Expected Status 'active', got '%s'", lease.Status)
	}

	if lease.UpdatedAt.Before(originalUpdatedAt) {
		t.Error("Expected UpdatedAt to be after activation")
	}
}

func TestCheckout(t *testing.T) {
	startDate, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	endDate, _ := time.Parse(time.RFC3339, "2023-12-31T23:59:59Z")

	lease := NewLease("123", "room-1", "landlord-1", "张三", "13812345678",
		startDate, endDate, 10000, 20000, "测试租约")
	lease.Activate()
	originalUpdatedAt := lease.UpdatedAt

	time.Sleep(10 * time.Millisecond)
	lease.Checkout()

	if lease.Status != "checkout" {
		t.Errorf("Expected Status 'checkout', got '%s'", lease.Status)
	}

	if lease.UpdatedAt.Before(originalUpdatedAt) {
		t.Error("Expected UpdatedAt to be after checkout")
	}
}

func TestUpdateLastChargeAt(t *testing.T) {
	startDate, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	endDate, _ := time.Parse(time.RFC3339, "2023-12-31T23:59:59Z")

	lease := NewLease("123", "room-1", "landlord-1", "张三", "13812345678",
		startDate, endDate, 10000, 20000, "测试租约")

	chargeDate, _ := time.Parse(time.RFC3339, "2023-06-01T10:00:00Z")

	lease.UpdateLastChargeAt(chargeDate)

	if lease.LastChargeAt == nil {
		t.Error("Expected LastChargeAt to be set")
	}

	if lease.LastChargeAt.Format(time.RFC3339) != chargeDate.Format(time.RFC3339) {
		t.Errorf("Expected LastChargeAt %v, got %v", chargeDate, lease.LastChargeAt)
	}
}

func TestIsActive(t *testing.T) {
	startDate, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	endDate, _ := time.Parse(time.RFC3339, "2023-12-31T23:59:59Z")

	lease := NewLease("123", "room-1", "landlord-1", "张三", "13812345678",
		startDate, endDate, 10000, 20000, "测试租约")

	if lease.IsActive() {
		t.Error("Expected lease to be inactive initially")
	}

	lease.Activate()

	if !lease.IsActive() {
		t.Error("Expected lease to be active after activation")
	}

	lease.Checkout()

	if lease.IsActive() {
		t.Error("Expected lease to be inactive after checkout")
	}
}
