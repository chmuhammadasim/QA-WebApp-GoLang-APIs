package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Answer struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	QuestionID primitive.ObjectID `json:"question_id"`
	Content    string             `json:"content"`
	CreatedAt  primitive.DateTime `json:"created_at"`
	UpdatedAt  primitive.DateTime `json:"updated_at"`
}
