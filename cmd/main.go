package main

import (
	"context"
	"lamoda/internal/logger"
	"lamoda/internal/server"
	"lamoda/internal/service"
	"lamoda/internal/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
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

	ctx, cancel := context.WithCancel(context.Background())
	conn, err := pgx.Connect(ctx, os.Getenv("DSN_STRING"))
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	log.Info("Connected to postgres")
	defer conn.Close(ctx)

	storage, err := storage.NewStorage(conn, log)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	service := service.NewService(storage)

	router := server.NewRouter(service, log, os.Getenv("RUN_ADDR"))
	srv := &http.Server{
		Addr:    os.Getenv("RUN_ADDR"),
		Handler: router,
	}
	var exit chan bool
	go gracefulShutdown(ctx, log, srv, conn, exit)

	log.Info("Server started")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}

	cancel()
	<-exit
}

func gracefulShutdown(ctx context.Context, log *logrus.Logger, server *http.Server, conn *pgx.Conn, exit chan bool) {
	// SIGTERM for docker container default signal
	signalCtx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// wait for parent or signal context to cancel
	<-signalCtx.Done()
	log.Info("shutting down http server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info("closing pgx connection...")
	conn.Close(shutdownCtx)

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error(
			"error shutting down http server",
			"error", err,
		)
	}
	exit <- true
}
