package server

import (
	"context"
	"lamoda/internal/model"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type serviceInterface interface {
	MakeReservation(ctx context.Context, products model.ReservedProducts) error
	DeleteReservation(ctx context.Context, ReservedProducts model.ReservedProducts) error
	SelectAvailableQty(ctx context.Context, stockID int) ([]model.Products, error)
}

type Server struct {
	service serviceInterface
	log     *logrus.Logger
}

func ErrorResponse(err error, rw http.ResponseWriter, r *http.Request, status int) {
	response := model.ErrorResponse{
		Error: err.Error(),
	}
	render.Status(r, status)
	render.JSON(rw, r, response)

}

// @Summary	Резервирование товаров
// @Tags		product
// @Produce	json
// @Param		object	body	model.ReservedProducts	true	"Товары для резервирования"
// @Success	201		"Successful response"
// @Failure	400		{object}	model.ErrorResponse
// @Failure	500		{object}	model.ErrorResponse
//
// @Router		/product [post]
func (server *Server) handlerMakeReservation(rw http.ResponseWriter, r *http.Request) {
	var reservedProducts model.ReservedProducts
	if err := render.DecodeJSON(r.Body, &reservedProducts); err != nil {
		server.log.Error(err.Error())
		ErrorResponse(err, rw, r, http.StatusBadRequest)
		return
	}

	if err := validator.New().Struct(reservedProducts); err != nil {
		server.log.Error(err.Error())
		ErrorResponse(err, rw, r, http.StatusBadRequest)
		return
	}
	if err := server.service.MakeReservation(r.Context(), reservedProducts); err != nil {
		ErrorResponse(err, rw, r, http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

// @Summary	Освобождение товаров из резерва
// @Tags		product
// @Produce	json
// @Param		object	body	model.ReservedProducts	true	"Товары для освобождения"
// @Success	200		"Successful response"
// @Failure	400		{object}	model.ErrorResponse
// @Failure	500		{object}	model.ErrorResponse
//
// @Router		/product [delete]
func (server *Server) handlerDeleteReservation(rw http.ResponseWriter, r *http.Request) {
	var reservedProducts model.ReservedProducts
	if err := render.DecodeJSON(r.Body, &reservedProducts); err != nil {
		server.log.Error(err.Error())
		ErrorResponse(err, rw, r, http.StatusBadRequest)
		return
	}

	if err := validator.New().Struct(reservedProducts); err != nil {
		server.log.Error(err.Error())
		ErrorResponse(err, rw, r, http.StatusBadRequest)
		return
	}
	if err := server.service.DeleteReservation(r.Context(), reservedProducts); err != nil {
		ErrorResponse(err, rw, r, http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

// @Summary	Получение кол-ва оставшихся товаров на складе
// @Tags		stock
// @Produce	json
// @Param		object	body		model.Stock	true	"ID склада"
// @Success	200		{object}	model.Products
// @Failure	400		{object}	model.ErrorResponse
// @Failure	500		{object}	model.ErrorResponse
//
// @Router		/stock [get]
func (server *Server) handlerGetAvailableQty(rw http.ResponseWriter, r *http.Request) {
	var stock model.Stock
	if err := render.DecodeJSON(r.Body, &stock); err != nil {
		server.log.Error(err.Error())
		ErrorResponse(err, rw, r, http.StatusBadRequest)
		return
	}

	if err := validator.New().Struct(stock); err != nil {
		server.log.Error(err.Error())
		ErrorResponse(err, rw, r, http.StatusBadRequest)
		return
	}
	availableProducts, err := server.service.SelectAvailableQty(r.Context(), stock.StockID)
	if err != nil {
		ErrorResponse(err, rw, r, http.StatusInternalServerError)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(rw, r, availableProducts)
}
