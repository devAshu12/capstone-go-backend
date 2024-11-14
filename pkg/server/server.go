package server

import (
	"github/devAshu12/learning_platform_GO_backend/internal/middlewares"
	"github/devAshu12/learning_platform_GO_backend/pkg/routes"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	http.Server
}

func NewServer() *Server {
	mux := mux.NewRouter()
	mux.Use(middlewares.ErrorHandlingMiddleware)
	routes.RegisterRoutes(mux)
	return &Server{
		Server: http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
	}
}
