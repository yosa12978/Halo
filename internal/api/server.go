package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/yosa12978/halo/internal/api/routers"
)

func Run() {
	go func(router *mux.Router) {
		http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router)
	}(routers.InitRouter())
}
