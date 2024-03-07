package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionInvoice = "invoices"

type Invoice struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	InvoiceID      string             `json:"invoice_id" bson:"invoice_id"`
	OrderID        string             `json:"order_id" bson:"order_id"`
	PaymentMethod  string             `json:"payment_method" bson:"payment_method"`
	PaymentStatus  string             `json:"payment_status" bson:"payment_status"`
	PaymentDueDate time.Time          `json:"payment_due_date" bson:"payment_due_date"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}
