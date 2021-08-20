package routers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yosa12978/halo/internal/api/handlers"
)

func InitRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	api := router.PathPrefix("/api").Subrouter()

	ah := handlers.NewAuthHandler()

	api.HandleFunc("/login", ah.LogIn).Methods(http.MethodPost)
	api.HandleFunc("/signup", ah.SignUp).Methods(http.MethodPost)

	return router
}
