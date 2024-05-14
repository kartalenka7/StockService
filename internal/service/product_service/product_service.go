package product_service

import (
	"context"
	"lamoda/internal/model"
)

type productStorer interface {
	ReserveProducts(ctx context.Context, products model.ReservedProducts) error
	DeleteReservation(ctx context.Context, ReservedProducts model.ReservedProducts) error
	CheckStockAvailability(ctx context.Context, stockID int) error
}

type ProductService interface {
	MakeReservation(ctx context.Context, products model.ReservedProducts) error
	DeleteReservation(ctx context.Context, products model.ReservedProducts) error
}

type productService struct {
	storage productStorer
}

func NewProductService(storer productStorer) ProductService {
	return &productService{storage: storer}
}

func (p *productService) MakeReservation(ctx context.Context, products model.ReservedProducts) error {
	if err := p.storage.CheckStockAvailability(ctx, products.StockID); err != nil {
		return err
	}
	return p.storage.ReserveProducts(ctx, products)
}

func (p *productService) DeleteReservation(ctx context.Context, products model.ReservedProducts) error {
	if err := p.storage.CheckStockAvailability(ctx, products.StockID); err != nil {
		return err
	}
	return p.storage.DeleteReservation(ctx, products)
}
