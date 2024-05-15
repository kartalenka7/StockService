package service

import (
	"context"
	"lamoda/internal/model"
	"lamoda/internal/service/stock_service"

	"github.com/jackc/pgx/v5"
)

func (suite *ServiceTestSuite) TestStockService_SelectAvailableQty() {
	ctx := context.Background()

	tests := []struct {
		name              string
		stockID           int
		want              []model.Products
		StockAvailableErr error
		GetErr            error
	}{
		// TODO: Add test cases.
		{
			name:    "Success",
			stockID: suite.rnd.Intn(20),
			want: []model.Products{
				{
					ProductID: suite.rnd.Intn(20),
					Quantity:  int64(suite.rnd.Intn(20)),
				},
			},
			StockAvailableErr: nil,
			GetErr:            nil,
		},
		{
			name:              "Stock is not available",
			want:              nil,
			StockAvailableErr: model.ErrStockIsNotAvailable,
			GetErr:            model.ErrStockIsNotAvailable,
		},
		{
			name:              "No rows found",
			want:              nil,
			StockAvailableErr: nil,
			GetErr:            pgx.ErrNoRows,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.stockStorer.EXPECT().CheckStockAvailability(ctx, tc.stockID).Return(tc.StockAvailableErr)
			if tc.StockAvailableErr == nil {
				suite.stockStorer.EXPECT().GetAvailableQty(ctx,
					tc.stockID).Return(tc.want, tc.GetErr)
			}
			service := stock_service.NewStockService(suite.stockStorer)
			got, err := service.SelectAvailableQty(ctx, tc.stockID)
			suite.Equal(tc.GetErr, err)
			if err == nil {
				suite.Equal(got, tc.want)
			}
		})
	}

}
