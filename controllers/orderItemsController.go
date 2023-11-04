package controllers

import (
	"context"
	"log"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/kushagra-gupta01/Restaurant-Management/database"
	"github.com/kushagra-gupta01/Restaurant-Management/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderItemPack struct{
	Table_id *string
	Order_items []*model.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client,"orderItem") 

func GetOrderItems() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)
		result,err := orderItemCollection.Find(context.TODO(),bson.M{})
		defer cancel()

		if err !=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error while listing ordered items"})
			return
		}
		var allOrderItems []bson.M
		if err := result.All(ctx, &allOrderItems);err!=nil{
			log.Fatal(err)
			return
		}
		c.JSON(http.StatusOK,allOrderItems)
	}
}

func GetOrderItem() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(),100*time.Second)
		orderItemId := c.Param("order_item_id")

		var orderItem model.OrderItem
		err := orderItemCollection.FindOne(ctx,bson.M{"orderItem_id":orderItemId}).Decode(&orderItem)
		defer cancel()
		if err !=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while listing order item"})
			return
		}
		c.JSON(http.StatusOK,orderItem)
	}
}

func CreateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)

		var orderItemPack OrderItemPack
		var order model.Order

		if err:=c.BindJSON(&orderItemPack);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		order.Order_Date,_ =time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		
	}
}

func UpdateOrderItem() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		
	}
}

func ItemsByOrder(id string)(OrderItems []primitive.M,err error){
	
}

func GetOrderItemsByOrder() gin.HandlerFunc{
	return func(c *gin.Context) {
		orderId := c.Param("order_id")
		allOrderItems,err := ItemsByOrder(orderId)
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error while listing order items by order"})
			return
		}
		c.JSON(http.StatusOK,allOrderItems)
	}
}