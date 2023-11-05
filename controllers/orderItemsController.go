package controllers

import (
	"context"
	"log"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kushagra-gupta01/Restaurant-Management/database"
	"github.com/kushagra-gupta01/Restaurant-Management/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		
		orderItemsToBeInserted := []interface{}{}
		order.Table_id = orderItemPack.Table_id
		order_id := OrderItemOrderCreator(order)

		for _,orderItem := range orderItemPack.Order_items{
			orderItem.Order_id = order_id

			ValidationErr := validate.Struct(orderItem)
			if ValidationErr !=nil{
				c.JSON(http.StatusBadRequest,gin.H{"error":ValidationErr.Error()})
				return
			}

			orderItem.ID = primitive.NewObjectID()
			orderItem.Created_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
			orderItem.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
			orderItem.Order_item_id = orderItem.ID.Hex()

			var num = toFixed(*orderItem.Unit_price,2)
			orderItem.Unit_price = &num
			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
		}

		insertedOrderItems, err:= orderItemCollection.InsertMany(ctx,orderItemsToBeInserted)

		if err!=nil{
			log.Fatal(err)
		}
		defer cancel()

		c.JSON(http.StatusOK,insertedOrderItems)
	}
}

func UpdateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)
		
		var orderItem model.OrderItem
		orderItemId := c.Param("order_item_id")
		fliter := bson.M{"order_item_id":orderItemId}

		var updateObj primitive.D
		if orderItem.Unit_price !=nil{
			updateObj = append(updateObj, bson.E{"unit_price",&*orderItem.Unit_price})
		}
		orderItem.Quantity !=nil{
			updateObj = append(updateObj , bson.E{"quantity",*orderItem.Quantity})
		}
		orderItem.Food_id!=nil{
			updateObj = append(updateObj, bson.E{"food_id", *orderItem.Food_id})
		}

		orderItem.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", orderItem.Updated_at})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result,err := orderItemCollection.UpdateOne(
			ctx,
			fliter,
			bson.D{
				{"$set", updateObj}
			},
			&opt,
		)

		if err!=nil{
			msg := "order item update failed"
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK,result)
	}
}

func ItemsByOrder(id string)(OrderItems []primitive.M,err error){
	
}

func GetOrderItemsByOrder()  	gin.HandlerFunc{
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