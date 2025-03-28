package main

import (
	"net"
	"net/http"
	"rest-api/internal/admin"
	"rest-api/internal/user"
	"rest-api/pkg/logging"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {

	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()

	logger.Info("register user handler")
	userHandler := user.NewHandler(logger)
	userHandler.Register(router)

	logger.Info("register admin handler")
	adminHandler := admin.NewHandler()
	adminHandler.Register(router)

	start(router)

}

func start(router *httprouter.Router) {

	logger := logging.GetLogger()
	logger.Info("start application")
	listenet, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info("application is listening on http://0.0.0.0:8080")
	logger.Fatal(server.Serve(listenet))
}
