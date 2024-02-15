package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SignUpInput struct {
	Name     string `json:"name" binding:"required,min=2,max=64"`
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type SignInInput struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type User struct {
	Id       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
}
