package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"qa-app/db"
	"qa-app/models"
	"qa-app/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	var userData models.User
	err := json.NewDecoder(r.Body).Decode(&userData)
	log.Println(userData)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Set hashed password
	userData.Password = string(hashedPassword)
	userData.Role = "user" // Default role

	_, err = db.GetCollection("users").InsertOne(context.Background(), userData)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	utils.SendResponse(w, http.StatusCreated, "User created successfully")
}

func Login(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var user models.User
	err = db.GetCollection("users").FindOne(context.Background(), bson.M{"username": loginData.Username}).Decode(&user)
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.ID.Hex(),
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(utils.GetSecretKey()))
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to create token")
		return
	}

	utils.SendResponse(w, http.StatusOK, map[string]interface{}{
		"token": tokenString,
	})
}
