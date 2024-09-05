package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"qa-app/db"
	"qa-app/models"
	"qa-app/utils"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateQuestion creates a new question
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

// GetAllQuestions retrieves all questions
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

	if err := cursor.Err(); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Cursor error")
		return
	}

	utils.SendResponse(w, http.StatusOK, questions)
}

// UpdateQuestion updates an existing question by its ID
func UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		utils.SendError(w, http.StatusBadRequest, "Missing question ID")
		return
	}

	var updatedQuestion models.Question
	err := json.NewDecoder(r.Body).Decode(&updatedQuestion)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid question ID")
		return
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"title":     updatedQuestion.Title,
			"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	result, err := db.GetCollection("questions").UpdateOne(context.Background(), filter, update)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to update question")
		return
	}

	if result.MatchedCount == 0 {
		utils.SendError(w, http.StatusNotFound, "Question not found")
		return
	}

	utils.SendResponse(w, http.StatusOK, "Question updated successfully")
}

// DeleteQuestion deletes a question by its ID
func DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		utils.SendError(w, http.StatusBadRequest, "Missing question ID")
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid question ID")
		return
	}

	filter := bson.M{"_id": objectID}

	result, err := db.GetCollection("questions").DeleteOne(context.Background(), filter)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to delete question")
		return
	}

	if result.DeletedCount == 0 {
		utils.SendError(w, http.StatusNotFound, "Question not found")
		return
	}

	utils.SendResponse(w, http.StatusOK, "Question deleted successfully")
}
