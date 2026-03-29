package facade

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	billrepo "github.com/zouhang1992/ddd_domain/internal/domain/bill/repository"
	depositrepo "github.com/zouhang1992/ddd_domain/internal/domain/deposit/repository"
)

// IncomeHandler 收入汇总 HTTP 处理器
type IncomeHandler struct {
	billRepo    billrepo.BillRepository
	depositRepo depositrepo.DepositRepository
}

// NewIncomeHandler 创建收入汇总处理器
func NewIncomeHandler(billRepo billrepo.BillRepository, depositRepo depositrepo.DepositRepository) *IncomeHandler {
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

	bills, err := h.billRepo.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.depositRepo.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var report IncomeReport
	report.Year = year
	report.Month = int(mon)

	for _, bill := range bills {
		if bill.PaidAt != nil {
			report.Total += bill.Amount
		}
	}

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
