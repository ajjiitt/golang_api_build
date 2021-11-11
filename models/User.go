package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id"`
	Username    *string            `json:"username" validate:"required,min=4",max=100`
	Name        *string            `json:"name" validate:"required,min=2,max=100"`
	DOB         string             `json:"dob" validate:"required"`
	Address     *string            `json:"address" validate:"required,max=150"`
	Password    *string            `json:"password" validate:"required,min=6"`
	Description *string            `json:"description" validate:"required"`
	Created_at  time.Time          `json:"created_at"`
}

type UserLogin struct {
	Username    *string            `json:"username" validate:"required,min=4",max=100`
	Password    *string            `json:"password" validate:"required,min=6"`
}