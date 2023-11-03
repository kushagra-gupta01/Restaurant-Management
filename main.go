package main

import (
	"os"
	"github.com/gin-gonic/gin"
	"github.com/kushagra-gupta01/Restaurant-Management/database"
	"github.com/kushagra-gupta01/Restaurant-Management/middleware"
	"github.com/kushagra-gupta01/Restaurant-Management/routes"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main(){
	port :=os.Getenv("PORT")
	if port ==""{
		port = "8000"
	}
	router :=gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.InvoiceRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.TableRoutes(router)

	router.Run(":" +port)
}