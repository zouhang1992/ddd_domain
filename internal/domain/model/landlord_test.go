package model

import (
	"testing"
	"time"
)

func TestNewLandlord(t *testing.T) {
	// 测试创建新的房东
	landlord := NewLandlord("123", "张三", "13812345678", "测试用户")

	if landlord.ID != "123" {
		t.Errorf("Expected ID '123', got '%s'", landlord.ID)
	}

	if landlord.Name != "张三" {
		t.Errorf("Expected name '张三', got '%s'", landlord.Name)
	}

	if landlord.Phone != "13812345678" {
		t.Errorf("Expected phone '13812345678', got '%s'", landlord.Phone)
	}

	if landlord.Note != "测试用户" {
		t.Errorf("Expected note '测试用户', got '%s'", landlord.Note)
	}

	// 检查创建时间是否合理
	if landlord.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if landlord.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}

	if landlord.CreatedAt != landlord.UpdatedAt {
		t.Error("Expected CreatedAt and UpdatedAt to be the same initially")
	}
}

func TestUpdate(t *testing.T) {
	// 测试更新房东信息
	landlord := NewLandlord("123", "张三", "13812345678", "测试用户")
	originalUpdatedAt := landlord.UpdatedAt

	time.Sleep(10 * time.Millisecond)
	landlord.Update("李四", "13912345678", "测试用户2")

	if landlord.Name != "李四" {
		t.Errorf("Expected name '李四', got '%s'", landlord.Name)
	}

	if landlord.Phone != "13912345678" {
		t.Errorf("Expected phone '13912345678', got '%s'", landlord.Phone)
	}

	if landlord.Note != "测试用户2" {
		t.Errorf("Expected note '测试用户2', got '%s'", landlord.Note)
	}

	if landlord.UpdatedAt.Before(originalUpdatedAt) {
		t.Error("Expected UpdatedAt to be after the original time")
	}
}

func TestLandlordEquality(t *testing.T) {
	// 测试相等性
	landlord1 := NewLandlord("123", "张三", "13812345678", "测试用户")
	landlord2 := NewLandlord("123", "张三", "13812345678", "测试用户")

	if landlord1.ID != landlord2.ID {
		t.Error("IDs should be the same")
	}
}
