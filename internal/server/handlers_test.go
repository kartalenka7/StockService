package server

import (
	"bytes"
	"context"
	"fmt"
	"lamoda/internal/logger"
	"lamoda/internal/model"
	mock_server "lamoda/internal/server/mocks"
	mrand "math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var rnd *mrand.Rand

func TestHandlerMakeReservation(t *testing.T) {
	rnd = mrand.New(mrand.NewSource(time.Now().UnixNano()))
	ctrl := gomock.NewController(t)
	service := mock_server.NewMockserviceInterface(ctrl)
	server := Server{
		service: service,
		log:     logger.InitLogger("fatal"),
	}
	ctx := context.Background()

	tests := []struct {
		name       string
		products   model.ReservedProducts
		serviceErr error
		wantStatus int
	}{
		// TODO: Add test cases.
		{
			name: "Success",
			products: model.ReservedProducts{
				StockID: rnd.Intn(20),
				Products: []model.Products{
					{
						ProductID: rnd.Intn(20),
						Quantity:  int64(rnd.Intn(20)),
					},
				},
			},
			serviceErr: nil,
			wantStatus: http.StatusOK,
		},

		{
			name: "Empty StockID",
			products: model.ReservedProducts{
				Products: []model.Products{
					{
						ProductID: rnd.Intn(20),
						Quantity:  int64(rnd.Intn(20)),
					},
				},
			},
			serviceErr: nil,
			wantStatus: http.StatusBadRequest,
		},

		{
			name: "Quantity < 0",
			products: model.ReservedProducts{
				StockID: rnd.Intn(20),
				Products: []model.Products{
					{
						ProductID: rnd.Intn(20),
						Quantity:  -int64(rnd.Intn(20)),
					},
				},
			},
			serviceErr: nil,
			wantStatus: http.StatusBadRequest,
		},

		{
			name: "Empty product part",
			products: model.ReservedProducts{
				StockID: rnd.Intn(20),
			},
			serviceErr: nil,
			wantStatus: http.StatusBadRequest,
		},

		{
			name: "Stock is not found",
			products: model.ReservedProducts{
				StockID: rnd.Intn(20),
				Products: []model.Products{
					{
						ProductID: rnd.Intn(20),
						Quantity:  int64(rnd.Intn(20)),
					},
				},
			},
			serviceErr: model.ErrStockNotFound,
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.wantStatus != http.StatusBadRequest {
				service.EXPECT().MakeReservation(ctx, gomock.Any()).Return(tt.serviceErr)
			}

			var productId int
			var qty int64

			if len(tt.products.Products) > 0 {
				productId = tt.products.Products[0].ProductID
				qty = tt.products.Products[0].Quantity
			}
			input := fmt.Sprintf(`{
				"products": [
				  {
					"product_id": %d,
					"quantity": %d
				  }
				],
				"stock_id": %d
			  }`, productId, qty, tt.products.StockID)
			req, err := http.NewRequest(http.MethodPost, "/product", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			recorder := httptest.NewRecorder()

			server.handlerMakeReservation(recorder, req)

			assert.Equal(t, recorder.Code, tt.wantStatus)
		})
	}
}
