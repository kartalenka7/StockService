package service

import (
	"context"
	"lamoda/internal/model"
	product "lamoda/internal/service/product_service"
	stock "lamoda/internal/service/stock_service"
)

type serviceStorer interface {
	ReserveProducts(ctx context.Context, ReservedProducts model.ReservedProducts) error
	DeleteReservation(ctx context.Context, ReservedProducts model.ReservedProducts) error
	GetAvailableQty(ctx context.Context, stockID int) ([]model.Products, error)
	CheckStockAvailability(ctx context.Context, stockID int) error
}

type ServiceMethods interface {
	product.ProductService
	stock.StockService
}
type service struct {
	product.ProductService
	stock.StockService
}

func NewService(storer serviceStorer) ServiceMethods {
	return &service{
		product.NewProductService(storer),
		stock.NewStockService(storer),
	}
}
