package main

import (
	"fmt"
	"lamoda/internal/logger"
	"lamoda/internal/server"
	"lamoda/internal/service"
	"lamoda/internal/storage"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// @title			API для управления товарами на складах
// @version		1.0
// @description	С помощью API можно резервировать товары, снимать резерв и получать информацию о доступных на складе товарах
//
// @host			localhost:3030
// @BasePath		/
func main() {

	godotenv.Load()
	log := logger.InitLogger(os.Getenv("LOG_LEVEL"))
	storage, err := storage.NewStorage(os.Getenv("DSN_STRING"), log)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer storage.Close()

	service := service.NewService(storage)

	router := server.NewRouter(service, log, os.Getenv("RUN_ADDR"))
	srv := &http.Server{
		Addr:    os.Getenv("RUN_ADDR"),
		Handler: router,
	}
	fmt.Println(os.Getenv("RUN_ADDR"))
	log.Info("Server started")
	if err := srv.ListenAndServe(); err != nil {
		log.Error(err.Error())
	}
}
