package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Employee struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	Name     string             `json:"name,omitempty" validate:"required"`
	Location string             `json:"location,omitempty" validate:"required"`
	Title    string             `json:"title,omitempty" validate:"required"`
	Email    string             `json:"email,omitempty" validate:"required"`
        
}
