package routes

import (
	controller "golang_api_build/controllers"
	middleware "golang_api_build/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.PUT("/user", controller.CreateUser())
	incomingRoutes.POST("/user", controller.Login())
	incomingRoutes.Use(middleware.Authuser())
	incomingRoutes.GET("/user/:user", controller.GetUser())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.POST("/update/:username", controller.UpdateUser())
	incomingRoutes.POST("/delete", controller.DeleteUser())
}