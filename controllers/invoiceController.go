package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kushagra-gupta01/Restaurant-Management/database"
	"github.com/kushagra-gupta01/Restaurant-Management/model"
	"github.com/kushagra-gupta01/Restaurant-Management/routes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InvoiceViewFormat struct{
	Invoice_id				string
	Payment_method		string
	Order_id					string
	Payment_status		*string
	Payment_due				interface{}
	Table_number			interface{}
	Payment_due_date	time.Time
	Order_details			interface{}
}

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")


func GetInvoices() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx,cancel := context.WithTimeout(context.Background(), 100*time.Second)
		result,err := invoiceCollection.Find(context.TODO(), bson.M{})
		defer cancel()

		if err!=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error ocurred while listing invoice items"})
		}
	var allInvoices []bson.M

	if err := result.All(ctx, &allInvoices); err !=nil{
		log.Fatal(err)
	}
	c.JSON(http.StatusOK,allInvoices)
	}
}

func GetInvoice() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		invoiceId := c.Param("invoice_id")

		var invoice model.Invoice
		err := invoiceCollection.FindOne(ctx,gin.H{"invoice_id":invoiceId}).Decode(&invoice)
		defer cancel()
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while listing invoice item"})
		}

		var invoiceView InvoiceViewFormat

		allOrderItems,err := ItemsByOrder(invoice.Order_id)
		invoiceView.Order_id = invoice.Order_id
		invoiceView.Payment_due_date = invoice.Payment_due_date
		
		invoiceView.Payment_method = "NULL"
		if invoiceView.Payment_method != nil{
			invoiceView.Payment_method = *invoice.Payment_method
		}

		invoiceView.Invoice_id = invoice.Invoice_id
		invoiceView.Payment_status = *&invoice.Payment_status
		invoiceView.Payment_due = allOrderItems[0]["payment_due"]
		invoiceView.Table_number = allOrderItems[0]["table_number"]
		invoiceView.Order_details = allOrderItems[0]["order_items"]

		c.JSON(http.StatusOK,invoiceView)
	}
}

func CreateInvoice() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func UpdateInvoice() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		
		var invoice model.Invoice
		invoiceId := c.Param("invoice_id")

		if err := c.BindJSON(&invoice); err !=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		filter := bson.M{"invoice_id":invoiceId}

		var updateObj primitive.D
	}
}