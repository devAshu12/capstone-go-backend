package routes

import (
	"github/devAshu12/learning_platform_GO_backend/internal/middlewares"
	"github/devAshu12/learning_platform_GO_backend/pkg/handlers"

	"github.com/gorilla/mux"
)

func CourseRoute(router *mux.Router) {

	courseRouter := router.PathPrefix("/course").Subrouter()

	courseRouter.Use(middlewares.AuthMiddleware)
	courseRouter.HandleFunc("", handlers.GetCourses).Methods("GET")
	courseRouter.HandleFunc("", handlers.CreateCourse).Methods("POST")
}
