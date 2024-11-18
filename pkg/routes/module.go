package routes

import (
	"github/devAshu12/learning_platform_GO_backend/internal/middlewares"
	"github/devAshu12/learning_platform_GO_backend/pkg/handlers"

	"github.com/gorilla/mux"
)

func ModuleRoute(router *mux.Router) {

	courseRouter := router.PathPrefix("/module").Subrouter()

	courseRouter.Use(middlewares.AuthMiddleware)
	courseRouter.HandleFunc("", handlers.CreateModule).Methods("POST")
	courseRouter.HandleFunc("", handlers.GetModules).Methods("GET")
}
