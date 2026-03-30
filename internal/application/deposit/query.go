package deposit

import (
	depositmodel "github.com/zouhang1992/ddd_domain/internal/domain/deposit/model"
)

// DepositQuery 押金查询接口
type DepositQuery interface {
	QueryName() string
}

// GetDepositQuery 获取押金查询
type GetDepositQuery struct {
	ID string
}

func (q GetDepositQuery) QueryName() string {
	return "get_deposit"
}

// ListDepositsQuery 列出押金查询
type ListDepositsQuery struct {
	LeaseID string
	Status  string
	Offset  int
	Limit   int
}

func (q ListDepositsQuery) QueryName() string {
	return "list_deposits"
}

// DepositQueryResult 押金查询结果
type DepositQueryResult struct {
	*depositmodel.Deposit
}

// DepositsQueryResult 押金列表查询结果
type DepositsQueryResult struct {
	Items []*depositmodel.Deposit `json:"items"`
	Total int                    `json:"total"`
	Page  int                    `json:"page"`
	Limit int                    `json:"limit"`
}
