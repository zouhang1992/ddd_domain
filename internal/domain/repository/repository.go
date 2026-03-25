package repository

// Repository 基础仓储接口
type Repository[T any, ID any] interface {
	Save(entity T) error
	FindByID(id ID) (T, error)
	FindAll() ([]T, error)
	Delete(id ID) error
}
