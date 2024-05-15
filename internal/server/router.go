package server

import (
	_ "lamoda/docs"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	swagger "github.com/swaggo/http-swagger"
)

func NewRouter(service serviceInterface, log *logrus.Logger, runAddr string) chi.Router {
	router := chi.NewRouter()

	server := &Server{
		service: service,
		log:     log,
	}

	router.Route("/product", func(r chi.Router) {
		r.Post("/", server.handlerMakeReservation)
		r.Delete("/", server.handlerDeleteReservation)
	})

	router.Get("/stock/{id}", server.handlerGetAvailableQty)
	router.Get("/swagger/*", swagger.Handler(
		swagger.URL("http://localhost"+runAddr+"/swagger/doc.json"),
	))

	return router
}
