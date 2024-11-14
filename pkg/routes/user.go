package routes

import (
	"github/devAshu12/learning_platform_GO_backend/pkg/handlers"

	"github.com/gorilla/mux"
)

func UserRouters(router *mux.Router) {

	userRouter := router.PathPrefix("/user").Subrouter()

	// Apply JWT Middleware
	// todoRouter.Use(middlewares.JWTMiddleware)

	userRouter.HandleFunc("/track-progress", handlers.UpdateProgress).Methods("POST")
	userRouter.HandleFunc("/register", handlers.Register).Methods("POST")
	userRouter.HandleFunc("/login", handlers.Login).Methods("POST")
}
