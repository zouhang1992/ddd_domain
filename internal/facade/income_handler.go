package facade

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
)

// IncomeHandler 收入汇总 HTTP 处理器
type IncomeHandler struct {
	billRepo    repository.BillRepository
	depositRepo repository.DepositRepository
}

// NewIncomeHandler 创建收入汇总处理器
func NewIncomeHandler(billRepo repository.BillRepository, depositRepo repository.DepositRepository) *IncomeHandler {
	return &IncomeHandler{billRepo: billRepo, depositRepo: depositRepo}
}

// RegisterRoutes 注册路由
func (h *IncomeHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /income", h.GetIncome)
}

// IncomeReport 收入报告
type IncomeReport struct {
	Year                  int    `json:"year"`
	Month                 int    `json:"month"`
	RentTotal             int64  `json:"rent_total"`
	WaterTotal            int64  `json:"water_total"`
	ElectricTotal         int64  `json:"electric_total"`
	OtherTotal            int64  `json:"other_total"`
	DepositIncome         int64  `json:"deposit_income"`
	DepositExpense        int64  `json:"deposit_expense"`
	Total                 int64  `json:"total"`
	TotalFormatted        string `json:"total_formatted"`
	RentFormatted         string `json:"rent_formatted"`
	WaterFormatted        string `json:"water_formatted"`
	ElectricFormatted     string `json:"electric_formatted"`
	OtherFormatted        string `json:"other_formatted"`
	DepositIncomeFormatted string `json:"deposit_income_formatted"`
	DepositExpenseFormatted string `json:"deposit_expense_formatted"`
}

// formatMoney 格式化金额
func formatMoney(amount int64) string {
	return stringFormat("%.2f", float64(amount)/100)
}

// stringFormat 字符串格式化（简单实现）
func stringFormat(format string, a ...interface{}) string {
	var buf [64]byte
	n := 0
	for _, v := range a {
		switch val := v.(type) {
		case float64:
			n = copy(buf[n:], formatFloat64(val))
		}
	}
	return string(buf[:n])
}

// formatFloat64 格式化浮点数
func formatFloat64(f float64) string {
	var s string
	i := int(f)
	d := int((f - float64(i)) * 100)
	if d < 0 {
		d = -d
	}
	s = stringFormatInt(i) + "." + stringFormatInt2(d)
	return s
}

// stringFormatInt 格式化整数
func stringFormatInt(i int) string {
	return fmt.Sprintf("%d", i) // 使用 fmt.Sprintf 正确格式化整数
}

// stringFormatInt2 格式化两位数
func stringFormatInt2(i int) string {
	if i < 10 {
		return "0" + stringFormatInt(i)
	}
	return stringFormatInt(i)
}

// GetIncome 获取收入汇总
func (h *IncomeHandler) GetIncome(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	monthStr := query.Get("month")

	var year int
	var mon time.Month
	var err error

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

	bills, err := h.billRepo.FindByMonth(year, mon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deposits, err := h.depositRepo.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var report IncomeReport
	report.Year = year
	report.Month = int(mon)

	for _, bill := range bills {
		if bill.Status == "paid" && bill.PaidAt != nil {
			report.RentTotal += bill.RentAmount
			report.WaterTotal += bill.WaterAmount
			report.ElectricTotal += bill.ElectricAmount
			report.OtherTotal += bill.OtherAmount

			// 总金额计算：租金+水电+其他
			report.Total += bill.RentAmount + bill.WaterAmount + bill.ElectricAmount + bill.OtherAmount
		}
	}

	// 处理押金收入和支出
	for _, deposit := range deposits {
		// 检查押金是否在指定月份内创建（收入）
		if deposit.CreatedAt.Year() == year && deposit.CreatedAt.Month() == mon && deposit.Status == model.DepositStatusCollected {
			report.DepositIncome += deposit.Amount
		}

		// 检查押金是否在指定月份内退还（支出）
		if deposit.RefundedAt != nil && deposit.RefundedAt.Year() == year && deposit.RefundedAt.Month() == mon && deposit.Status == model.DepositStatusRefunded {
			report.DepositExpense += deposit.Amount
		}
	}

	// 调整总金额，加入押金净收入
	report.Total += report.DepositIncome - report.DepositExpense

	report.TotalFormatted = formatMoney(report.Total)
	report.RentFormatted = formatMoney(report.RentTotal)
	report.WaterFormatted = formatMoney(report.WaterTotal)
	report.ElectricFormatted = formatMoney(report.ElectricTotal)
	report.OtherFormatted = formatMoney(report.OtherTotal)
	report.DepositIncomeFormatted = formatMoney(report.DepositIncome)
	report.DepositExpenseFormatted = formatMoney(report.DepositExpense)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(report)
}
