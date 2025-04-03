package main

import (
	"fmt"
	"net"
	"net/http"
	"rest-api/internal/admin"
	"rest-api/internal/config"
	"rest-api/internal/user"
	client "rest-api/internal/user/db"
	"rest-api/pkg/logging"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	logger.Info("register user handler")

	mongo, err := client.NewMongoClient(cfg.Mongo.URI, logger)
	if err != nil {
		logger.Errorf("Can not connect to mongoDB %v", err)
	}
	NewMongoStorage := user.NewMongoStorage(mongo, cfg.Mongo.Database, cfg.Mongo.Collection, logger)
	userHandler := user.NewHandler(logger, NewMongoStorage)
	userHandler.Register(router)

	logger.Info("register admin handler")
	adminHandler := admin.NewHandler(logger, NewMongoStorage)
	adminHandler.Register(router)

	router.Handler("GET", "/metrics", promhttp.Handler())

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
