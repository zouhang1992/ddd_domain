package model

import (
	"testing"
	"time"
)

func TestNewOperationLog(t *testing.T) {
	// 测试创建操作日志
	id := "test-id-123"
	timestamp := time.Now()
	eventName := "landlord.created"
	domainType := "landlord"
	aggregateID := "landlord-123"
	operatorID := "user-456"
	action := "created"

	details := map[string]interface{}{
		"name": "张三",
		"phone": "13812345678",
	}
	metadata := map[string]interface{}{
		"source": "test",
	}

	log := NewOperationLog(id, timestamp, eventName, domainType, aggregateID, operatorID, action, details, metadata)

	if log.ID() != id {
		t.Errorf("Expected ID '%s', got '%s'", id, log.ID())
	}

	if log.EventName() != eventName {
		t.Errorf("Expected EventName '%s', got '%s'", eventName, log.EventName())
	}

	if log.DomainType() != domainType {
		t.Errorf("Expected DomainType '%s', got '%s'", domainType, log.DomainType())
	}

	if log.AggregateID() != aggregateID {
		t.Errorf("Expected AggregateID '%s', got '%s'", aggregateID, log.AggregateID())
	}

	if log.OperatorID() != operatorID {
		t.Errorf("Expected OperatorID '%s', got '%s'", operatorID, log.OperatorID())
	}

	if log.Action() != action {
		t.Errorf("Expected Action '%s', got '%s'", action, log.Action())
	}

	if log.Details()["name"] != "张三" {
		t.Errorf("Expected Details.name '张三', got '%v'", log.Details()["name"])
	}

	if log.Metadata()["source"] != "test" {
		t.Errorf("Expected Metadata.source 'test', got '%v'", log.Metadata()["source"])
	}

	if log.CreatedAt().IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

func TestOperationLogDetailsJSON(t *testing.T) {
	// 测试详细数据的JSON序列化和反序列化
	log := NewOperationLog(
		"test-id",
		time.Now(),
		"test.event",
		"test",
		"agg-123",
		"user-456",
		"created",
		map[string]interface{}{"key": "value", "number": 123},
		nil,
	)

	detailsJSON, err := log.MarshalDetails()
	if err != nil {
		t.Fatalf("Failed to marshal details: %v", err)
	}

	if detailsJSON == "" {
		t.Error("Expected non-empty details JSON")
	}

	newLog := NewOperationLog(
		"test-id",
		time.Now(),
		"test.event",
		"test",
		"agg-123",
		"user-456",
		"created",
		nil,
		nil,
	)

	err = newLog.UnmarshalDetails(detailsJSON)
	if err != nil {
		t.Fatalf("Failed to unmarshal details: %v", err)
	}

	if newLog.Details()["key"] != "value" {
		t.Errorf("Expected unmarshaled key 'value', got '%v'", newLog.Details()["key"])
	}
}

func TestOperationLogSetters(t *testing.T) {
	// 测试设置方法
	log := NewOperationLog(
		"test-id",
		time.Now(),
		"test.event",
		"test",
		"agg-123",
		"user-456",
		"created",
		nil,
		nil,
	)

	newDetails := map[string]interface{}{"new": "details"}
	newMetadata := map[string]interface{}{"new": "metadata"}

	log.SetDetails(newDetails)
	log.SetMetadata(newMetadata)

	if log.Details()["new"] != "details" {
		t.Error("Failed to set details")
	}

	if log.Metadata()["new"] != "metadata" {
		t.Error("Failed to set metadata")
	}
}
