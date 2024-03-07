package model

import "time"

type InvoiceResponse struct {
	InvoiceID      string    `json:"invoice_id"`
	OrderID        string    `json:"order_id"`
	PaymentDue     any       `json:"payment_due"`
	PaymentDueDate time.Time `json:"payment_due_date"`
	PaymentMethod  string    `json:"payment_method"`
	PaymentStatus  *string   `json:"payment_status"`
	TableNumber    any       `json:"table_number"`
	OrderDetail    any       `json:"order_detail"`
}
