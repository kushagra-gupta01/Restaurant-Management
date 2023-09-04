package routes

import(
	"github.com/gin-gonic/gin"
	"github.com/kushagra-gupta01/Restaurant-Management/controlllers"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/users",controlllers.GetUsers())
	incomingRoutes.GET("/users/:user_id",controlllers.GetUser())
	incomingRoutes.POST("/users/singup",controlllers.SignUp())
	incomingRoutes.POST("/users/login",controlllers.Login())
}