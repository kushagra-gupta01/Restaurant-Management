package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/kushagra-gupta01/Restaurant-Management/model"
	"github.com/kushagra-gupta01/Restaurant-Management/routes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var OrderCollection *mongo.Collection = database.OpenCollection(database.Client,"order")

func GetOrders() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result,err :=OrderCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err !=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error while listing order items"})
		}
		var allOrders[]bson.M
		if err = result.All(ctx, &allOrders);err!=nil{
			log.Fatal(err)
		}
		c.JSON(http.StatusOK,allOrders)
	}
}

func GetOrder() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)
		orderId := c.Param("order_id")
		var order model.Order

		err:=OrderCollection.FindOne(ctx,bson.M{"order_id":orderId}).Decode(&order)
		defer cancel()
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error while fetching the orders"})
		}
		c.JSON(http.StatusOK,order)
	}
}

func CreateOrder() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		var order model.Order
		var table model.Table

		if err := c.BindJSON(&order); err!= nil{
			c.JSON(http.StatusBadRequest,gin.H{"Error":err.Error()})
			return
		}

		validationErr := validate.Struct(order)

		if validationErr!= nil{
			c.JSON(http.StatusBadRequest,gin.H{"Error":validationErr.Error()})
			return
		}

		if order.Table_id != nil{
			err:= TableCollection.FindOne(ctx,bson.M{"table_id":order.Table_id}).Decode(&table)
			defer cancel()

			if err !=nil{
				msg := fmt.Sprintf("message table not found")
				c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
				return 
			}
		}

		order.Created_at,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Updated_at,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		order.ID = primitive.NewObjectID()
		order.Order_id = order.ID.Hex()

		result,InsertErr := OrderCollection.InsertOne(ctx,order)
		if InsertErr != nil{
			msg := fmt.Sprintf("order item was not created")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return 
		}
		defer cancel()
		c.JSON(http.StatusOK,result)
	}
}

func UpdateOrder() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)
		var order model.Order
		var table model.Table
		orderId := c.Param("order_id")

		if err := c.BindJSON(&food);err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		var updateObj primitive.D
		if order.Table_id!=nil{
			err := OrderCollection.FindOne(ctx,bson.M{"table_id":order.Table_id}).Decode(&order)
			defer cancel()
			if err != nil{
				msg := fmt.Sprintf("message:order not found")
				c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
				return
			}
			updateObj = append(updateObj,bson.E{"table_id":order.Table_id})
		}
		order.Updated_at, _ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		updateObj = append(updateObj,bson.E{"updated_at",order.Updated_at})

		upsert := true
		filter := bson.M{"order_id":orderId}
		opt :=  options.UpdateOptions{
			Upsert: &upsert,
		}

		result,err := OrderCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set",updateObj},
			},
			&opt,
		)
		if err!=nil{
			msg := fmt.Sprintf("order update failed")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK,result)
	}
}

func OrderItemOrderCreator(order model.Order) string{
	
	order.Created_at,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))

	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()
	
	OrderCollection.InsertOne(ctx,order)
	defer cancel()

	return order.Order_id
}