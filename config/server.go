package config

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func InitServer() (*httprouter.Router, *http.Server) {
	router := httprouter.New()
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	return router, &server
}
