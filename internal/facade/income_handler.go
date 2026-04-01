package facade

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/zouhang1992/ddd_domain/internal/infrastructure/middleware"
	billrepo "github.com/zouhang1992/ddd_domain/internal/domain/bill/repository"
	depositrepo "github.com/zouhang1992/ddd_domain/internal/domain/deposit/repository"
	leaserepo "github.com/zouhang1992/ddd_domain/internal/domain/lease/repository"
	roomrepo "github.com/zouhang1992/ddd_domain/internal/domain/room/repository"
)

// IncomeHandler 收入汇总 HTTP 处理器
type IncomeHandler struct {
	billRepo        billrepo.BillRepository
	depositRepo     depositrepo.DepositRepository
	leaseRepo       leaserepo.LeaseRepository
	roomRepo        roomrepo.RoomRepository
	authMiddleware  *middleware.AuthMiddleware
}

// NewIncomeHandler 创建收入汇总处理器
func NewIncomeHandler(
	billRepo billrepo.BillRepository, 
	depositRepo depositrepo.DepositRepository, 
	leaseRepo leaserepo.LeaseRepository, 
	roomRepo roomrepo.RoomRepository,
	authMiddleware *middleware.AuthMiddleware,
) *IncomeHandler {
	return &IncomeHandler{
		billRepo:       billRepo, 
		depositRepo:    depositRepo, 
		leaseRepo:      leaseRepo, 
		roomRepo:       roomRepo,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes 注册路由
func (h *IncomeHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /income", h.authMiddleware.RequireAuth(h.GetIncome))
}

// IncomeReport 收入报告
type IncomeReport struct {
	Year                  int    `json:"year"`
	Month                 int    `json:"month"`

	// 收入部分
	RentIncome            int64  `json:"rent_income"`
	WaterIncome           int64  `json:"water_income"`
	ElectricIncome        int64  `json:"electric_income"`
	OtherIncome           int64  `json:"other_income"`
	DepositIncome         int64  `json:"deposit_income"`

	// 支出部分
	RentExpense           int64  `json:"rent_expense"`
	DepositExpense        int64  `json:"deposit_expense"`

	// 计算结果
	TotalIncome           int64  `json:"total_income"`
	TotalExpense          int64  `json:"total_expense"`
	NetIncome             int64  `json:"net_income"`

	// 格式化字符串
	RentIncomeFormatted     string `json:"rent_income_formatted"`
	WaterIncomeFormatted    string `json:"water_income_formatted"`
	ElectricIncomeFormatted string `json:"electric_income_formatted"`
	OtherIncomeFormatted    string `json:"other_income_formatted"`
	DepositIncomeFormatted  string `json:"deposit_income_formatted"`
	RentExpenseFormatted    string `json:"rent_expense_formatted"`
	DepositExpenseFormatted string `json:"deposit_expense_formatted"`
	TotalIncomeFormatted    string `json:"total_income_formatted"`
	TotalExpenseFormatted   string `json:"total_expense_formatted"`
	NetIncomeFormatted      string `json:"net_income_formatted"`
}

// formatMoney 格式化金额
func formatMoney(amount int64) string {
	return fmt.Sprintf("%.2f", float64(amount)/100)
}

// isSameMonth 检查时间是否在指定年月
func isSameMonth(t time.Time, year int, month time.Month) bool {
	return t.Year() == year && t.Month() == month
}

// GetIncome 获取收入汇总
func (h *IncomeHandler) GetIncome(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	monthStr := query.Get("month")
	locationID := query.Get("location_id")

	var year int
	var mon time.Month

	if monthStr == "" {
		// 默认当前月份
		now := time.Now()
		year = now.Year()
		mon = now.Month()
	} else {
		t, parseErr := time.Parse("2006-01", monthStr)
		if parseErr != nil {
			http.Error(w, "invalid month format (should be YYYY-MM)", http.StatusBadRequest)
			return
		}
		year = t.Year()
		mon = t.Month()
	}

	bills, err := h.billRepo.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deposits, err := h.depositRepo.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 构建 leaseID -> locationID 的映射（即使没有位置筛选也需要，用于押金筛选）
	leaseToLocation := make(map[string]string)
	leases, err := h.leaseRepo.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rooms, err := h.roomRepo.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 构建 roomID -> locationID 映射
	roomToLocation := make(map[string]string)
	for _, room := range rooms {
		roomToLocation[room.ID()] = room.LocationID
	}

	// 构建 leaseID -> locationID 映射
	for _, lease := range leases {
		if locID, ok := roomToLocation[lease.RoomID]; ok {
			leaseToLocation[lease.ID()] = locID
		}
	}

	var report IncomeReport
	report.Year = year
	report.Month = int(mon)

	// 计算账单收入和支出（只统计已支付且在指定月份的账单）
	for _, bill := range bills {
		// 位置筛选
		if locationID != "" {
			billLocationID := leaseToLocation[bill.LeaseID]
			if billLocationID != locationID {
				continue
			}
		}

		if bill.PaidAt != nil && isSameMonth(*bill.PaidAt, year, mon) {
			if bill.Type == "checkout" {
				// 退租结算账单：分别处理租金的收入和支出
				if bill.RentAmount < 0 {
					// 负数表示退还租金，算作支出
					report.RentExpense += -bill.RentAmount
				} else {
					// 正数表示收取租金，算作收入
					report.RentIncome += bill.RentAmount
				}
				// 水费、电费、其他费用始终算作收入
				report.WaterIncome += bill.WaterAmount
				report.ElectricIncome += bill.ElectricAmount
				report.OtherIncome += bill.OtherAmount
			} else {
				// 普通账单：所有金额算作收入
				report.RentIncome += bill.RentAmount
				report.WaterIncome += bill.WaterAmount
				report.ElectricIncome += bill.ElectricAmount
				report.OtherIncome += bill.OtherAmount
			}
		}
	}

	// 计算押金收入和支出
	for _, deposit := range deposits {
		// 位置筛选
		if locationID != "" {
			depositLocationID := leaseToLocation[deposit.LeaseID]
			if depositLocationID != locationID {
				continue
			}
		}

		// 押金收入：created_at 在指定月份
		if isSameMonth(deposit.CreatedAt, year, mon) {
			report.DepositIncome += deposit.Amount
		}
		// 押金支出：refunded_at 在指定月份
		if deposit.RefundedAt != nil && isSameMonth(*deposit.RefundedAt, year, mon) {
			report.DepositExpense += deposit.Amount
		}
	}

	// 计算总计
	report.TotalIncome = report.RentIncome + report.WaterIncome + report.ElectricIncome + report.OtherIncome + report.DepositIncome
	report.TotalExpense = report.RentExpense + report.DepositExpense
	report.NetIncome = report.TotalIncome - report.TotalExpense

	// 格式化金额
	report.RentIncomeFormatted = formatMoney(report.RentIncome)
	report.WaterIncomeFormatted = formatMoney(report.WaterIncome)
	report.ElectricIncomeFormatted = formatMoney(report.ElectricIncome)
	report.OtherIncomeFormatted = formatMoney(report.OtherIncome)
	report.DepositIncomeFormatted = formatMoney(report.DepositIncome)
	report.RentExpenseFormatted = formatMoney(report.RentExpense)
	report.DepositExpenseFormatted = formatMoney(report.DepositExpense)
	report.TotalIncomeFormatted = formatMoney(report.TotalIncome)
	report.TotalExpenseFormatted = formatMoney(report.TotalExpense)
	report.NetIncomeFormatted = formatMoney(report.NetIncome)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(report)
}
