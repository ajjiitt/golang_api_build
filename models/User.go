package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)
//struct to maintain user data
type User struct {
	ID          primitive.ObjectID `bson:"_id"`
	Username    string            `json:"username" validate:"required,min=4",max=100`
	Name        string            `json:"name" validate:"required,min=2,max=100"`
	DOB         string             `json:"dob" validate:"required"`
	Address     string            `json:"address" validate:"required,max=150"`
	Password    string            `json:"password" validate:"required,min=6"`
	Description string            `json:"description" validate:"required"`
	Longitude   string            `json:"long"`
	Latitude    string            `json:"lat"`
	Created_at  time.Time          `json:"created_at"`
}
// struct to get data while login in user
type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password" validate:"required,min=6"`
}
// struct to send user data with token after successful login
type Userdata struct {
	Usr   User `json:"user"`
	Token string      `json:"token"`
}