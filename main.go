package main

import (
	"fmt"
	routes "golang_api_build/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("<-- API -->")
	//loading env variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	//setting application port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router := gin.New()
	//gin logger for clean UI and logging of api calls
	router.Use(gin.Logger())
	// routing user-routes
	routes.UserRoutes(router)
	//test api 
	router.GET("/api/status", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "API sent successfully"})
	})

	router.Run(":" + port)
}
