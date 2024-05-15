package product

import (
	"context"
	"lamoda/internal/model"
	"lamoda/internal/service/mocks"
	"testing"
	"time"

	mrand "math/rand"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ProductTestSuite struct {
	suite.Suite
	mockCtrl       *gomock.Controller
	storer         *mocks.MockproductStorer
	rnd            *mrand.Rand
	productService *ProductService
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ProductTestSuite))
}

func (suite *ProductTestSuite) SetupSuite() {
	suite.rnd = mrand.New(mrand.NewSource(time.Now().UnixNano()))
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.storer = mocks.NewMockproductStorer(suite.mockCtrl)
	suite.productService = NewProductService(suite.storer)
}

func (suite *ProductTestSuite) TearDownSuite() {
	suite.mockCtrl.Finish()
}

func (suite *ProductTestSuite) TestMakeReservation_Success() {
	ctx := context.Background()

	products := model.ReservedProducts{
		StockID: suite.rnd.Intn(20),
		Products: []model.Products{
			{
				ProductID: suite.rnd.Intn(20),
				Quantity:  int64(suite.rnd.Intn(20)),
			},
		},
	}

	gomock.InOrder(
		suite.storer.EXPECT().CheckStockAvailability(
			gomock.AssignableToTypeOf(ctx),
			products.StockID,
		).Return(nil),
		suite.storer.EXPECT().ReserveProduct(
			gomock.AssignableToTypeOf(ctx),
			products.StockID, products.Products[0]).Return(nil),
	)

	err := suite.productService.MakeReservation(ctx, products)
	require.NoError(suite.T(), err)
}

func (suite *ProductTestSuite) TestMakeReservation_StockNotAvailable() {
	ctx := context.Background()

	products := model.ReservedProducts{
		StockID: suite.rnd.Intn(20),
		Products: []model.Products{
			{
				ProductID: suite.rnd.Intn(20),
				Quantity:  int64(suite.rnd.Intn(20)),
			},
		},
	}

	gomock.InOrder(
		suite.storer.EXPECT().CheckStockAvailability(
			gomock.AssignableToTypeOf(ctx),
			products.StockID,
		).Return(model.ErrStockIsNotAvailable),
	)

	err := suite.productService.MakeReservation(ctx, products)
	require.EqualError(suite.T(), err, model.ErrStockIsNotAvailable.Error())
}

func (suite *ProductTestSuite) TestMakeReservation_NotEnoughQty() {
	ctx := context.Background()

	products := model.ReservedProducts{
		StockID: suite.rnd.Intn(20),
		Products: []model.Products{
			{
				ProductID: suite.rnd.Intn(20),
				Quantity:  int64(suite.rnd.Intn(20)),
			},
		},
	}

	gomock.InOrder(
		suite.storer.EXPECT().CheckStockAvailability(
			gomock.AssignableToTypeOf(ctx),
			products.StockID,
		).Return(nil),
		suite.storer.EXPECT().ReserveProduct(
			gomock.AssignableToTypeOf(ctx),
			products.StockID, products.Products[0]).Return(model.ErrNotEnoughAvailableQty),
	)

	err := suite.productService.MakeReservation(ctx, products)
	require.EqualError(suite.T(), err, model.ErrNotEnoughAvailableQty.Error())
}

func (suite *ProductTestSuite) TestDeleteReservation_Success() {
	ctx := context.Background()

	products := model.ReservedProducts{
		StockID: suite.rnd.Intn(20),
		Products: []model.Products{
			{
				ProductID: suite.rnd.Intn(20),
				Quantity:  int64(suite.rnd.Intn(20)),
			},
		},
	}

	gomock.InOrder(
		suite.storer.EXPECT().CheckStockAvailability(
			gomock.AssignableToTypeOf(ctx),
			products.StockID,
		).Return(nil),
		suite.storer.EXPECT().DeleteReservation(
			gomock.AssignableToTypeOf(ctx),
			products.StockID, products.Products[0]).Return(nil),
	)

	err := suite.productService.DeleteReservation(ctx, products)
	require.NoError(suite.T(), err)
}

func (suite *ProductTestSuite) TestDeleteReservation_StockNotAvailable() {
	ctx := context.Background()

	products := model.ReservedProducts{
		StockID: suite.rnd.Intn(20),
		Products: []model.Products{
			{
				ProductID: suite.rnd.Intn(20),
				Quantity:  int64(suite.rnd.Intn(20)),
			},
		},
	}

	gomock.InOrder(
		suite.storer.EXPECT().CheckStockAvailability(
			gomock.AssignableToTypeOf(ctx),
			products.StockID,
		).Return(model.ErrStockIsNotAvailable),
	)

	err := suite.productService.DeleteReservation(ctx, products)
	require.EqualError(suite.T(), err, model.ErrStockIsNotAvailable.Error())
}

func (suite *ProductTestSuite) TestDeleteReservation_ProductNotFound() {
	ctx := context.Background()

	products := model.ReservedProducts{
		StockID: suite.rnd.Intn(20),
		Products: []model.Products{
			{
				ProductID: suite.rnd.Intn(20),
				Quantity:  int64(suite.rnd.Intn(20)),
			},
		},
	}

	gomock.InOrder(
		suite.storer.EXPECT().CheckStockAvailability(
			gomock.AssignableToTypeOf(ctx),
			products.StockID,
		).Return(nil),
		suite.storer.EXPECT().DeleteReservation(
			gomock.AssignableToTypeOf(ctx),
			products.StockID, products.Products[0]).Return(model.ErrProductNotFound),
	)

	err := suite.productService.DeleteReservation(ctx, products)
	require.EqualError(suite.T(), err, model.ErrProductNotFound.Error())
}
