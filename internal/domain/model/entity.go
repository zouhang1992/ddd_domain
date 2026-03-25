package model

// Entity 基础实体接口
type Entity interface {
	ID() string
}

// BaseEntity 基础实体实现
type BaseEntity struct {
	id string
}

// NewBaseEntity 创建基础实体
func NewBaseEntity(id string) BaseEntity {
	return BaseEntity{id: id}
}

// ID 实现 Entity 接口
func (e BaseEntity) ID() string {
	return e.id
}

// ValueObject 值对象接口
type ValueObject interface {
	Equals(other ValueObject) bool
}
