package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/muhafs/go-restaurant-management/controllers"
)

func InvoiceRoutes(r *gin.Engine) {
	r.POST("/invoices", controllers.CreateInvoice)
	r.GET("/invoices", controllers.FindInvoices)
	r.GET("/invoices/:invoice_id", controllers.FindInvoice)
	r.PATCH("/invoices/:invoice_id", controllers.UpdateInvoice)
}
