package routers

import "github.com/gorilla/mux"

func InitRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	return router
}
