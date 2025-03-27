package main

import (
	"log"
	"net"
	"net/http"
	"rest-api/internal/admin"
	"rest-api/internal/user"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {

	log.Println("create router")
	router := httprouter.New()

	log.Println("register user handler")
	userHandler := user.NewHandler()
	userHandler.Register(router)

	log.Println("register admin handler")
	adminHandler := admin.NewHandler()
	adminHandler.Register(router)

	start(router)

}

func start(router *httprouter.Router) {

	log.Println("start application")
	listenet, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("application is listening on http://0.0.0.0:8080")
	log.Fatal(server.Serve(listenet))
}
