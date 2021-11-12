package controller

import (
	"context"
	"fmt"
	"golang_api_build/database"
	"golang_api_build/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)
// struct to send response - status and data
type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}
// struct for token 
type Claims struct {
	Username *string `json:"username"`
	jwt.StandardClaims
}
// struct to get longitude and latitude data
type LongLat struct {
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}

var jwtKey = []byte("secret_key")
//setting up collection
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

// contoller for getting single user
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userD := c.Param("user")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		//getting userdata wrt userD provided from DB
		defer cancel()
		err := userCollection.FindOne(ctx, bson.M{"username": userD}).Decode(&user)
		//error handling - while querying userdata
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		resp := Response{Status: "success", Data: user}
		c.JSON(http.StatusOK, resp)
	}
}
// controller for creating user
func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		defer cancel()
		//binds incoming request data to user struct
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//validate data based upon validation specified
		validationErr := validate.Struct(user)
		//handle validation error
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		//hashing password
		password := HashPassword(user.Password)
		user.Password = password
		//check whether user with given user data is present in db 
		count, err := userCollection.CountDocuments(ctx, bson.M{"username": user.Username})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the username"})
			return
		}
		//if count is greater than zero means user with given data is present and throw error
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this username already exists"})
			return
		}
		// declaring createdAt and objectID for user
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		// inserting user to db
		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		//handling error for insertion in db
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User item was not created"})
			return
		}
		defer cancel()
		resp := Response{Status: "success", Data: resultInsertionNumber}
		c.JSON(http.StatusOK, resp)
	}
}
//contoller for loggin in user
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.UserLogin
		var foundUser models.User
		//binds incoming request data to user struct
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//getting user data from db, for checking it with request data
		err := userCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&foundUser)
		defer cancel()
		//handling error for invalid user
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username or password is incorrect"})
			return
		}
		// checking password - password in user and password in db
		passwordIsValid, msg := VerifyPassword(user.Password, foundUser.Password)
		defer cancel()
		// handling error for invalid password
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		//error handling
		if foundUser.Username == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		}
		err = userCollection.FindOne(ctx, bson.M{"username": foundUser.Username}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// setting and creating JWT token 
		expirationTime := time.Now().Add(time.Minute * 5000)
		claims := &Claims{
			Username: &foundUser.Username,
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
		//sending token and user data
		resp := Response{Status: "success", Data: models.Userdata{Token: tokenString, Usr: foundUser}}
		c.JSON(http.StatusOK, resp)
	}
}
// controller for updating user
func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User
		//binds incoming data to struct data
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// checking whether user is present in db
		err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		// password checking for updating user
		passwordIsValid, msg := VerifyPassword(user.Password, foundUser.Password)
		defer cancel()
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		// updating user data
		filter := bson.M{"username": username}
		update := bson.M{"$set": bson.M{"name": user.Name, "description": user.Description, "address": user.Address, "dob": user.DOB,"longitude":user.Longitude,"latitude":user.Latitude}}
		result, err := userCollection.UpdateOne(context.Background(), filter, update)
		//handling error - any error caused during updating user
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("modified count: ", result.ModifiedCount)
		// fetaching updated user data
		err = userCollection.FindOne(ctx, bson.M{"username": foundUser.Username}).Decode(&user)
		defer cancel()
		//handling error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		//response -  user data
		resp := Response{Status: "success", Data: user}
		c.JSON(http.StatusOK, resp)
	}
}
// controller for deleting user
func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.UserLogin

		var foundUser models.User
		// binds incoming user data to struct user
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// checking whether the user is present or not
		err := userCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&foundUser)
		defer cancel()
		//handling error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username or password is incorrect"})
			return
		}
		// validating password in db
		passwordIsValid, msg := VerifyPassword(user.Password, foundUser.Password)
		defer cancel()
		// handling password error 
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		// filtering data to be deleted
		filter := bson.M{"username": foundUser.Username}
		result, err := userCollection.DeleteOne(context.Background(), filter)
		// handling updation error
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("delete count: ", result)
		// sending user response back
		resp := Response{Status: "data deleted successfully", Data: result.DeletedCount}
		c.JSON(http.StatusOK, resp)
	}
}
// controller for getting all users
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		// getAllUsers fetchs all users from db
		resp := Response{Status: "success", Data: getAllUsers()}
		c.JSON(http.StatusOK, resp)
	}
}
// controller to get user by coordinates
func FilterUserViaLongLat() gin.HandlerFunc {
	return func(c *gin.Context) {
		var usr []models.User
		user := getAllUsers()
		var coordinate LongLat
		// binds incoming data to LongLat struct
		if err := c.BindJSON(&coordinate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
		}
		// converting incoming string data to int
		longD, _ := strconv.ParseFloat(coordinate.Longitude, 32)
		latD, _ := strconv.ParseFloat(coordinate.Latitude, 32)
		longDT := int(longD)
		latDT := int(latD)
		// mapping and checking for user with given coordinates
		for i := range user {
			lat, _ := strconv.ParseFloat(user[i].Latitude, 32)
			long, _ := (strconv.ParseFloat(user[i].Longitude, 32))
			longT := int(long)
			latT := int(lat)
			if longDT == longT && latDT == latT {
				usr = append(usr, user[i])
			}
		}
		// handling not users found error
		if len(usr) == 0 {
			resp := Response{Status: "No user found", Data: []int{}}
			c.JSON(http.StatusOK, resp)
			return
		}
		// sending user fetched with given coordinates
		resp := Response{Status: "success", Data: usr}
		c.JSON(http.StatusOK, resp)
	}
}
// function to fetch all users
func getAllUsers() []models.User {
	cur, err := userCollection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var users []models.User
	// converting bson data to array
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
// function to hashpassword for storing it to db securly
func HashPassword(password string) string {
	// hashing password to store securly
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}
// function to verfy hashed password and password sent via user 
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	// comparing given and hashed password from db
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""
	// handling error 
	if err != nil {
		msg = "email of password is incorrect"
		check = false
	}
	return check, msg
}
