package routes

import (
	controller "golang_api_build/controllers"
	middleware "golang_api_build/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	//this are unsecured routes - initially used for login and creating users
	incomingRoutes.PUT("/user", controller.CreateUser())
	incomingRoutes.POST("/user", controller.Login())
	//middleware check
	incomingRoutes.Use(middleware.Authuser())
	//all routes below middleware as protected using JWT token

	// to get single user data
	incomingRoutes.GET("/user/:user", controller.GetUser())
	//to get all users 
	incomingRoutes.GET("/users", controller.GetUsers())
	// update users data
	incomingRoutes.POST("/update/:username", controller.UpdateUser())
	//delete users data
	incomingRoutes.POST("/delete", controller.DeleteUser())
	// filter users based upon given coordinates
	incomingRoutes.POST("/filterbycoordinates",controller.FilterUserViaLongLat())
}