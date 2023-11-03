package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/kushagra-gupta01/Restaurant-Management/routes"
	"go.keploy.io/server/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
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
		c.JSON(http.StatusOK,allOrders[0])
	}
}

func GetOrder() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)
		orderId := c.Param("order_id")
		var order models.Order

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
		var order models.Order
		var table models.Table
	}
}

func UpdateOrder() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)
		var order models.Order
		var table models.Table
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