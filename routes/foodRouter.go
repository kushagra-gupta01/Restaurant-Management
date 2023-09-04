package routes

import(
	"github.com/gin-gonic/gin"
	"github.com/kushagra-gupta01/Restaurant-Management/controllers"
)

func FoodRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/foods",controlllers.GetFoods())
	incomingRoutes.GET("/foods/:food_id",controlllers.GetFood())
	incomingRoutes.POST("/foods",controlllers.CreateFood())
	incomingRoutes.PATCH("/food/:food_id",controlllers.UpdateFood())
}