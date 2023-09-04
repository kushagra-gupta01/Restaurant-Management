package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/kushagra-gupta01/Restaurant-Management/routes"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetOrderItems() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func GetOrderItem() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func CreateOrderItem() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func UpdateOrderItem() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		
	}
}

func ItemsByOrder(id string)(OrderItems []primitive.M,err error){
	
}

func GetOrderItemsByOrder() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}