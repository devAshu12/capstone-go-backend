package routes

import (
	"github/devAshu12/learning_platform_GO_backend/internal/middlewares"
	"github/devAshu12/learning_platform_GO_backend/pkg/handlers"

	"github.com/gorilla/mux"
)

func VideoRouters(router *mux.Router) {

	videoRouter := router.PathPrefix("/video").Subrouter()

	// Apply JWT Middleware
	videoRouter.Use(middlewares.AuthMiddleware)

	videoRouter.HandleFunc("", handlers.CreateVideo).Methods("POST")
	videoRouter.HandleFunc("", handlers.RemoveVideo).Methods("DELETE")
}
