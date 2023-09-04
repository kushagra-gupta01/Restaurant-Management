package routes

import(
	"github.com/gin-gonic/gin"
	"github.com/kushagra-gupta01/Restaurant-Management/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/users",controllers.GetUsers())
	incomingRoutes.GET("/users/:user_id",controllers.GetUser())
	incomingRoutes.POST("/users/singup",controllers.SignUp())
	incomingRoutes.POST("/users/login",controllers.Login())
}