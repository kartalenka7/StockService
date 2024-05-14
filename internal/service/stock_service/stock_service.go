package stock_service

import (
	"context"
	"lamoda/internal/model"
)

type stockStorer interface {
	GetAvailableQty(ctx context.Context, stockID int) ([]model.Products, error)
	CheckStockAvailability(ctx context.Context, stockID int) error
}

type StockService interface {
	SelectAvailableQty(ctx context.Context, stockID int) ([]model.Products, error)
}

type stockService struct {
	storage stockStorer
}

func NewStockService(storer stockStorer) StockService {
	return &stockService{
		storage: storer,
	}
}

func (s *stockService) SelectAvailableQty(ctx context.Context, stockID int) ([]model.Products, error) {
	if err := s.storage.CheckStockAvailability(ctx, stockID); err != nil {
		return nil, err
	}
	return s.storage.GetAvailableQty(ctx, stockID)
}
