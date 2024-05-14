package model

import "errors"

var (
	ErrStockIsNotAvailable   = errors.New("Stock is not available")
	ErrNotEnoughAvailableQty = errors.New("Available quantity on stock is not enough for reservation")
	ErrProductNotFound       = errors.New("Product not found on stock")
	ErrStockNotFound         = errors.New("Stock is not found")
)

type Products struct {
	ProductID int   `json:"product_id" validate:"required"`
	Quantity  int64 `json:"quantity" validate:"gt=0"`
}

type ReservedProducts struct {
	StockID  int        `json:"stock_id" validate:"required"`
	Products []Products `validate:"dive"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Stock struct {
	StockID int `json:"stock_id" validate:"required"`
}
