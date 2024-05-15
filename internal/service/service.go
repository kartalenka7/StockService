package service

import (
	"context"
	"lamoda/internal/model"
	product "lamoda/internal/service/product_service"
	stock "lamoda/internal/service/stock_service"
)

type serviceStorer interface {
	ReserveProduct(ctx context.Context, stockId int, product model.Products) error
	DeleteReservation(ctx context.Context, stockId int, product model.Products) error
	GetAvailableQty(ctx context.Context, stockID int) ([]model.Products, error)
	CheckStockAvailability(ctx context.Context, stockID int) error
}

type service struct {
	*product.ProductService
	*stock.StockService
}

func NewService(storer serviceStorer) *service {
	return &service{
		product.NewProductService(storer),
		stock.NewStockService(storer),
	}
}
