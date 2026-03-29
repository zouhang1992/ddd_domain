package print

import "errors"

// PrintBillCommand 打印账单命令
type PrintBillCommand struct {
	BillID string
}

// CommandName 实现 Command 接口
func (c PrintBillCommand) CommandName() string {
	return "print_bill"
}

// Validate 验证命令
func (c PrintBillCommand) Validate() error {
	if c.BillID == "" {
		return errors.New("bill id is required")
	}
	return nil
}

// PrintLeaseCommand 打印租约命令
type PrintLeaseCommand struct {
	LeaseID string
}

// CommandName 实现 Command 接口
func (c PrintLeaseCommand) CommandName() string {
	return "print_lease"
}

// Validate 验证命令
func (c PrintLeaseCommand) Validate() error {
	if c.LeaseID == "" {
		return errors.New("lease id is required")
	}
	return nil
}

// PrintInvoiceCommand 打印发票命令
type PrintInvoiceCommand struct {
	BillID string
}

// CommandName 实现 Command 接口
func (c PrintInvoiceCommand) CommandName() string {
	return "print_invoice"
}

// Validate 验证命令
func (c PrintInvoiceCommand) Validate() error {
	if c.BillID == "" {
		return errors.New("bill id is required")
	}
	return nil
}
