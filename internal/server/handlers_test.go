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

	"github.com/go-chi/chi/v5"
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
			if tt.wantStatus != http.StatusBadRequest {
				service.EXPECT().MakeReservation(req.Context(), gomock.Any()).Return(tt.serviceErr)
			}

			server.handlerMakeReservation(recorder, req)

			assert.Equal(t, recorder.Code, tt.wantStatus)
		})
	}
}

func TestHandlerDeleteReservation(t *testing.T) {
	rnd = mrand.New(mrand.NewSource(time.Now().UnixNano()))
	ctrl := gomock.NewController(t)
	service := mock_server.NewMockserviceInterface(ctrl)
	server := Server{
		service: service,
		log:     logger.InitLogger("fatal"),
	}

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
			req, err := http.NewRequest(http.MethodDelete, "/product", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			if tt.wantStatus != http.StatusBadRequest {
				service.EXPECT().DeleteReservation(req.Context(), gomock.Any()).Return(tt.serviceErr)
			}

			recorder := httptest.NewRecorder()
			server.handlerDeleteReservation(recorder, req)

			assert.Equal(t, recorder.Code, tt.wantStatus)
		})
	}
}

func TestHandlerGetAvailableQty(t *testing.T) {
	rnd = mrand.New(mrand.NewSource(time.Now().UnixNano()))
	ctrl := gomock.NewController(t)
	service := mock_server.NewMockserviceInterface(ctrl)
	server := Server{
		service: service,
		log:     logger.InitLogger("fatal"),
	}

	tests := []struct {
		name         string
		stockId      int
		serviceErr   error
		wantStatus   int
		wantProducts []model.Products
	}{
		// TODO: Add test cases.
		{
			name:       "Success",
			stockId:    rnd.Intn(20),
			serviceErr: nil,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Empty StockID",
			serviceErr: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Stock is not available",
			stockId:    rnd.Intn(20),
			serviceErr: model.ErrStockIsNotAvailable,
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := httptest.NewRequest("GET", "/{id}", nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", fmt.Sprint(tt.stockId))
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			if tt.wantStatus != http.StatusBadRequest {
				service.EXPECT().SelectAvailableQty(r.Context(), tt.stockId).Return(tt.wantProducts, tt.serviceErr)
			}
			recorder := httptest.NewRecorder()
			server.handlerGetAvailableQty(recorder, r)

			assert.Equal(t, tt.wantStatus, recorder.Code)

		})
	}
}
