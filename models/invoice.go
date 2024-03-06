package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID             primitive.ObjectID `bson:"_id"`
	InvoiceID      string             `json:"invoice_id"`
	OrderID        string             `json:"order_id"`
	PaymentMethod  *string            `json:"payment_method" validate:"eq=CARD|eq=CASH|eq="`
	PaymentStatus  *string            `json:"payment_status" validate:"required,eq=PENDING|eq=PAID"`
	PaymentDueDate time.Time          `json:"payment_due_date"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

type InvoiceResponse struct {
	InvoiceID      string    `json:"invoice_id"`
	OrderID        string    `json:"order_id"`
	PaymentDue     any       `json:"payment_due"`
	PaymentDueDate time.Time `json:"payment_due_date"`
	PaymentMethod  string    `json:"payment_method" validate:"eq=CARD|eq=CASH|eq="`
	PaymentStatus  *string   `json:"payment_status" validate:"required,eq=PENDING|eq=PAID"`
	TableNumber    any       `json:"table_number"`
	OrderDetail    any       `json:"order_detail"`
}
