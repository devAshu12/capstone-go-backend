package routes

import (
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {

	UserRouters(router)

}
