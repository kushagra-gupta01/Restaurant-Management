package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kushagra-gupta01/Restaurant-Management/controllers"
)

func OrderItemRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/orderItems",controllers.GetOrderItems())
	incomingRoutes.GET("/orderItems/:orderItem_id",controllers.GetOrderItem())
	incomingRoutes.GET("/orderItems-order/:order_id",controllers.GetOrderItemsByOrder())
	incomingRoutes.POST("/orderItems",controllers.CreateOrderItem())
	incomingRoutes.PATCH("/orderItems/:orderItem_id",controllers.UpdateOrderItem())
}