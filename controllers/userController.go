package controller

import (
	"context"
	"fmt"
	"golang_api_build/database"
	"golang_api_build/models"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Userdata struct {
	Usr models.User `json:"user"`
	Token string `json:"token"`
}

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type Claims struct {
	Username *string `json:"username"`
	jwt.StandardClaims
}

var jwtKey = []byte("secret_key")
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userD := c.Param("user")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"username": userD}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		resp := Response{Status: "success", Data: user}
		c.JSON(http.StatusOK, resp)
	}
}

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		defer cancel()
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		count, err := userCollection.CountDocuments(ctx, bson.M{"username": user.Username})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the username"})
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password
		count, err = userCollection.CountDocuments(ctx, bson.M{"username": user.Username})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the username"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this username already exists"})
			return
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User item was not created"})
			return
		}
		defer cancel()
		resp := Response{Status: "success", Data: resultInsertionNumber}
		c.JSON(http.StatusOK, resp)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.UserLogin
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := userCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Username == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		}
		err = userCollection.FindOne(ctx, bson.M{"username": foundUser.Username}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		expirationTime := time.Now().Add(time.Minute * 5)
		claims := &Claims{
			Username: foundUser.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		resp := Response{Status: "success", Data: Userdata{Token:tokenString, Usr:foundUser}}
		c.JSON(http.StatusOK, resp)
	}
}

func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		filter := bson.M{"username": username}
		update := bson.M{"$set": bson.M{"name": user.Name, "description": user.Description, "address": user.Address, "dob": user.DOB}}
		result, err := userCollection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("modified count: ", result.ModifiedCount)
		err = userCollection.FindOne(ctx, bson.M{"username": foundUser.Username}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		resp := Response{Status: "success", Data: user}
		c.JSON(http.StatusOK, resp)
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.UserLogin

		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := userCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		filter := bson.M{"username": foundUser.Username}
		result, err := userCollection.DeleteOne(context.Background(), filter)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("delete count: ", result)
		resp := Response{Status: "data deleted successfully", Data: result.DeletedCount}
		c.JSON(http.StatusOK, resp)
	}
}


func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := Response{Status: "success", Data: getAllUsers()}
		c.JSON(http.StatusOK, resp)
	}
}


func getAllUsers() []models.User{
	cur, err := userCollection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var users []models.User

	for cur.Next(context.Background()) {
		var user models.User
		err := cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	defer cur.Close(context.Background())
	return users
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "email of password is incorrect"
		check = false
	}
	return check, msg
}
