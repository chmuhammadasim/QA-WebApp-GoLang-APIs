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

	// Admin routes (only accessible by admin role)
	adminRoutes := api.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(middleware.RoleMiddleware("admin"))

	adminRoutes.HandleFunc("/questions", controllers.CreateQuestion).Methods("POST")
	adminRoutes.HandleFunc("/questions/{id}", controllers.UpdateQuestion).Methods("PUT")
	adminRoutes.HandleFunc("/questions/{id}", controllers.DeleteQuestion).Methods("DELETE")

	// User routes (accessible by user role)
	userRoutes := api.PathPrefix("/user").Subrouter()
	userRoutes.Use(middleware.RoleMiddleware("user"))

	userRoutes.HandleFunc("/questions", controllers.GetAllQuestions).Methods("GET")

	return router
}
