package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Question struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `json:"title"`
	Content   string             `json:"content"`
	CreatedAt primitive.DateTime `json:"created_at"`
	UpdatedAt primitive.DateTime `json:"updated_at"`
}
