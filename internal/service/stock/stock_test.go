package stock

import (
	"context"
	"lamoda/internal/model"
	"testing"
	"time"

	"lamoda/internal/service/mocks"

	mrand "math/rand"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type StockTestSuite struct {
	suite.Suite
	mockCtrl     *gomock.Controller
	storer       *mocks.MockstockStorer
	rnd          *mrand.Rand
	stockService *StockService
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(StockTestSuite))
}

func (suite *StockTestSuite) SetupSuite() {
	suite.rnd = mrand.New(mrand.NewSource(time.Now().UnixNano()))
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.storer = mocks.NewMockstockStorer(suite.mockCtrl)
	suite.stockService = NewStockService(suite.storer)
}

func (suite *StockTestSuite) TearDownSuite() {
	suite.mockCtrl.Finish()
}

func (suite *StockTestSuite) TestSelectAvailableQty_Success() {
	ctx := context.Background()
	stockID := suite.rnd.Intn(20)

	expectedProducts := []model.Products{
		{
			ProductID: suite.rnd.Intn(20),
			Quantity:  int64(suite.rnd.Intn(20)),
		},
	}

	gomock.InOrder(
		suite.storer.EXPECT().CheckStockAvailability(
			gomock.AssignableToTypeOf(ctx),
			stockID,
		).Return(nil),
		suite.storer.EXPECT().GetAvailableQty(
			gomock.AssignableToTypeOf(ctx),
			stockID,
		).Return(expectedProducts, nil),
	)

	actualProducts, err := suite.stockService.SelectAvailableQty(ctx, stockID)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), expectedProducts, actualProducts)
}

func (suite *StockTestSuite) TestSelectAvailableQty_StockIsNotAvailable() {
	ctx := context.Background()
	stockID := suite.rnd.Intn(20)

	gomock.InOrder(
		suite.storer.EXPECT().CheckStockAvailability(
			gomock.AssignableToTypeOf(ctx),
			stockID,
		).Return(model.ErrStockIsNotAvailable),
	)

	_, err := suite.stockService.SelectAvailableQty(ctx, stockID)
	require.EqualError(suite.T(), err, model.ErrStockIsNotAvailable.Error())
}

func (suite *StockTestSuite) TestSelectAvailableQty_NowRowsFound() {
	ctx := context.Background()
	stockID := suite.rnd.Intn(20)

	gomock.InOrder(
		suite.storer.EXPECT().CheckStockAvailability(
			gomock.AssignableToTypeOf(ctx),
			stockID,
		).Return(nil),
		suite.storer.EXPECT().GetAvailableQty(
			gomock.AssignableToTypeOf(ctx),
			stockID,
		).Return(nil, pgx.ErrNoRows),
	)

	actualProducts, err := suite.stockService.SelectAvailableQty(ctx, stockID)
	require.EqualError(suite.T(), err, pgx.ErrNoRows.Error())
	require.Nil(suite.T(), actualProducts)
}
