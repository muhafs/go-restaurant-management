package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhafs/go-restaurant-management/database"
	"github.com/muhafs/go-restaurant-management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

func FindInvoices(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	result, err := invoiceCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var invoices []bson.M
	if err := result.All(ctx, &invoices); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, invoices)
}

func FindInvoice(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	invoiceID := c.Param("invoice_id")
	var invoice models.Invoice

	if err := invoiceCollection.FindOne(ctx, bson.M{"invoiceid": invoiceID}).Decode(&invoice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items, err := ItemsByOrder(c, invoice.OrderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "items not found"})
		return
	}

	var invoiceResponse models.InvoiceResponse

	invoiceResponse.OrderID = invoice.OrderID
	invoiceResponse.PaymentDueDate = invoice.PaymentDueDate

	invoiceResponse.PaymentMethod = "null"
	if invoice.PaymentMethod != nil {
		invoiceResponse.PaymentMethod = *invoice.PaymentMethod
	}

	invoiceResponse.InvoiceID = invoice.InvoiceID
	invoiceResponse.PaymentStatus = invoice.PaymentStatus

	invoiceResponse.PaymentDue = items[0]["payment_due"]
	invoiceResponse.TableNumber = items[0]["table_number"]
	invoiceResponse.OrderDetail = items[0]["items"]

	c.JSON(http.StatusOK, invoiceResponse)
}

func CreateInvoice(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	var order models.Order
	var invoice models.Invoice

	if err := c.BindJSON(&invoice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if vErr := validate.Struct(invoice); vErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": vErr.Error()})
		return
	}

	if err := orderCollection.FindOne(ctx, bson.M{"order_id": invoice.OrderID}).Decode(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "order not found"})
		return
	}

	if invoice.PaymentStatus == nil {
		pending := "PENDING"
		invoice.PaymentStatus = &pending
	}

	invoice.PaymentDueDate, _ = time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
	invoice.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	invoice.ID = primitive.NewObjectID()
	invoice.InvoiceID = invoice.ID.Hex()

	result, err := invoiceCollection.InsertOne(ctx, invoice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invoice item failed to create"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func UpdateInvoice(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	invoiceID := c.Param("invoice_id")

	var invoice models.Invoice
	if err := c.BindJSON(&invoice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var updateObj bson.D

	if invoice.PaymentMethod != nil {
		updateObj = append(updateObj, bson.E{Key: "payment_method", Value: invoice.PaymentMethod})
	}

	if invoice.PaymentStatus != nil {
		updateObj = append(updateObj, bson.E{Key: "payment_status", Value: invoice.PaymentStatus})
	}

	invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: invoice.UpdatedAt})

	if invoice.PaymentStatus == nil {
		pending := "PENDING"
		invoice.PaymentStatus = &pending
	}

	result, err := invoiceCollection.UpdateOne(
		ctx,
		bson.M{"invoice_id": invoiceID},
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invoice item failed to update"})
		return
	}

	c.JSON(http.StatusOK, result)
}
