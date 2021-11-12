package middleware

import (
	"fmt"
	controller "golang_api_build/controllers"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("secret_key")

func Authuser() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
		}
		fmt.Println("Middleware running")
		validateToken(c,clientToken)
		fmt.Println("Middleware ran")
	}
}

func validateToken(c *gin.Context, tokenStr string) {
	claims := &controller.Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access"})
			c.Abort()
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "something went wrong"})
		c.Abort()
		return
	}

	if !tkn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access"})
		c.Abort()
		return
	}
	c.Next()
}
