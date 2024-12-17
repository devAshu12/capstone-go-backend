package server

import (
	"fmt"
	"github/devAshu12/learning_platform_GO_backend/internal/middlewares"
	"github/devAshu12/learning_platform_GO_backend/pkg/routes"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	http.Server
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, CORS is enabled with rs/cors!")
}

func NewServer() *Server {
	mux := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Replace with your allowed origin(s)
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	mux.Use(middlewares.ErrorHandlingMiddleware)
	routes.RegisterRoutes(mux)

	handler := c.Handler(mux)
	return &Server{
		Server: http.Server{
			Addr:    ":8080",
			Handler: handler,
		},
	}
}
