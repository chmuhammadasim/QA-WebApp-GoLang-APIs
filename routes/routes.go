package routes

import (
	"qa-app/controllers"
	"qa-app/middleware"

	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	// Auth routes
	router.HandleFunc("/signup", controllers.Signup).Methods("POST")
	router.HandleFunc("/login", controllers.Login).Methods("POST")

	// Question routes (protected)
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	api.HandleFunc("/questions", controllers.CreateQuestion).Methods("POST")
	api.HandleFunc("/questions", controllers.GetAllQuestions).Methods("GET")
	api.HandleFunc("/questions/{id}", controllers.UpdateQuestion).Methods("PUT")
	api.HandleFunc("/questions/{id}", controllers.DeleteQuestion).Methods("DELETE")

	return router
}
