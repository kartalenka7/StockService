package service

import (
	"context"
	"lamoda/internal/model"
	"lamoda/internal/service/product_service"
)

func (suite *ServiceTestSuite) TestMakeReservation() {

	ctx := context.Background()

	tests := []struct {
		name              string
		products          model.ReservedProducts
		StockAvailableErr error
		ReservationErr    error
	}{
		// TODO: Add test cases.
		{
			name: "Success",
			products: model.ReservedProducts{
				StockID: suite.rnd.Intn(20),
				Products: []model.Products{
					{
						ProductID: suite.rnd.Intn(20),
						Quantity:  int64(suite.rnd.Intn(20)),
					},
				},
			},
			StockAvailableErr: nil,
			ReservationErr:    nil,
		},
		{
			name: "Stock is not available",
			products: model.ReservedProducts{
				StockID: suite.rnd.Intn(20),
				Products: []model.Products{
					{
						ProductID: suite.rnd.Intn(20),
						Quantity:  int64(suite.rnd.Intn(20)),
					},
				},
			},
			StockAvailableErr: model.ErrStockIsNotAvailable,
			ReservationErr:    model.ErrStockIsNotAvailable,
		},
		{
			name: "Not enough Qty",
			products: model.ReservedProducts{
				StockID: suite.rnd.Intn(20),
				Products: []model.Products{
					{
						ProductID: suite.rnd.Intn(20),
						Quantity:  int64(suite.rnd.Intn(20)),
					},
				},
			},
			StockAvailableErr: nil,
			ReservationErr:    model.ErrNotEnoughAvailableQty,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.productStorer.EXPECT().CheckStockAvailability(ctx,
				tc.products.StockID).Return(tc.StockAvailableErr)
			if tc.StockAvailableErr == nil {
				suite.productStorer.EXPECT().ReserveProduct(ctx, tc.products.StockID,
					tc.products.Products[0]).Return(tc.ReservationErr)
			}

			service := product_service.NewProductService(suite.productStorer)
			err := service.MakeReservation(ctx, tc.products)
			suite.Equal(tc.ReservationErr, err)
		})
	}
}

func (suite *ServiceTestSuite) TestDeleteReservation() {
	ctx := context.Background()

	tests := []struct {
		name              string
		products          model.ReservedProducts
		StockAvailableErr error
		DeleteErr         error
	}{
		// TODO: Add test cases.
		{
			name: "Success",
			products: model.ReservedProducts{
				StockID: suite.rnd.Intn(20),
				Products: []model.Products{
					{
						ProductID: suite.rnd.Intn(20),
						Quantity:  int64(suite.rnd.Intn(20)),
					},
				},
			},
			StockAvailableErr: nil,
			DeleteErr:         nil,
		},
		{
			name: "Stock is not available",
			products: model.ReservedProducts{
				StockID: suite.rnd.Intn(20),
				Products: []model.Products{
					{
						ProductID: suite.rnd.Intn(20),
						Quantity:  int64(suite.rnd.Intn(20)),
					},
				},
			},
			StockAvailableErr: model.ErrStockIsNotAvailable,
			DeleteErr:         model.ErrStockIsNotAvailable,
		},
		{
			name: "Product not found",
			products: model.ReservedProducts{
				StockID: suite.rnd.Intn(20),
				Products: []model.Products{
					{
						ProductID: suite.rnd.Intn(20),
						Quantity:  int64(suite.rnd.Intn(20)),
					},
				},
			},
			StockAvailableErr: nil,
			DeleteErr:         model.ErrProductNotFound,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.productStorer.EXPECT().CheckStockAvailability(ctx,
				tc.products.StockID).Return(tc.StockAvailableErr)
			if tc.StockAvailableErr == nil {
				suite.productStorer.EXPECT().DeleteReservation(ctx, tc.products.StockID,
					tc.products.Products[0]).Return(tc.DeleteErr)
			}
			service := product_service.NewProductService(suite.productStorer)
			err := service.DeleteReservation(ctx, tc.products)
			suite.Equal(tc.DeleteErr, err)
		})
	}
}
