package main

import (
	"fmt"
	"net"
	"net/http"
	"rest-api/internal/admin"
	"rest-api/internal/config"
	"rest-api/internal/user"
	"rest-api/pkg/logging"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {

	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	logger.Info("register user handler")
	userHandler := user.NewHandler(logger)
	userHandler.Register(router)

	logger.Info("register admin handler")
	adminHandler := admin.NewHandler()
	adminHandler.Register(router)

	start(router, cfg)

}

func start(router *httprouter.Router, cfg *config.Config) {

	logger := logging.GetLogger()
	logger.Info("start application")
	listenaddress := fmt.Sprintf("%s:%s", cfg.Listen.Address, cfg.Listen.Port)
	listenet, err := net.Listen(cfg.Listen.Type, listenaddress)
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Infof("application is listening on %s", listenaddress)
	logger.Fatal(server.Serve(listenet))
}
