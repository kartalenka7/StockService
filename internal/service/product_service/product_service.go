package product_service

import (
	"context"
	"lamoda/internal/model"
)

type productStorer interface {
	ReserveProduct(ctx context.Context, stockId int, product model.Products) error
	DeleteReservation(ctx context.Context, stockId int, product model.Products) error
	CheckStockAvailability(ctx context.Context, stockID int) error
}

type ProductService struct {
	storage productStorer
}

func NewProductService(storer productStorer) *ProductService {
	return &ProductService{storage: storer}
}

func (p *ProductService) MakeReservation(ctx context.Context,
	products model.ReservedProducts) error {
	if err := p.storage.CheckStockAvailability(ctx, products.StockID); err != nil {
		return err
	}
	for _, product := range products.Products {
		err := p.storage.ReserveProduct(ctx, products.StockID, product)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *ProductService) DeleteReservation(ctx context.Context,
	products model.ReservedProducts) error {
	if err := p.storage.CheckStockAvailability(ctx, products.StockID); err != nil {
		return err
	}
	for _, product := range products.Products {
		err := p.storage.DeleteReservation(ctx, products.StockID, product)
		if err != nil {
			return err
		}
	}
	return nil
}
