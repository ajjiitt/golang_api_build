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
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)

	router.GET("/api/status", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "API sent successfully"})
	})

	router.Run(":" + port)
}
