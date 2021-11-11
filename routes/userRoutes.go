package routes

import (
	"github.com/gin-gonic/gin"
	controller "golang_api_build/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/user/:user", controller.GetUser())
	incomingRoutes.PUT("/user", controller.CreateUser())
	incomingRoutes.POST("/user", controller.Login())
	incomingRoutes.POST("/update/:username", controller.UpdateUser())
	incomingRoutes.POST("/delete", controller.DeleteUser())
}