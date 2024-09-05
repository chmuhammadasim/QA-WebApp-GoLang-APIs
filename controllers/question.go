package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"qa-app/db"
	"qa-app/models"
	"qa-app/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var question models.Question
	err := json.NewDecoder(r.Body).Decode(&question)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	question.ID = primitive.NewObjectID()
	question.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	question.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = db.GetCollection("questions").InsertOne(context.Background(), question)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to create question")
		return
	}

	utils.SendResponse(w, http.StatusCreated, question)
}

func GetAllQuestions(w http.ResponseWriter, r *http.Request) {
	cursor, err := db.GetCollection("questions").Find(context.Background(), bson.M{})
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to fetch questions")
		return
	}
	defer cursor.Close(context.Background())

	var questions []models.Question
	for cursor.Next(context.Background()) {
		var question models.Question
		err := cursor.Decode(&question)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Error decoding question")
			return
		}
		questions = append(questions, question)
	}

	utils.SendResponse(w, http.StatusOK, questions)
}
