package service

import (
	mocks "lamoda/internal/service/mocks"
	mrand "math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ServiceTestSuite struct {
	suite.Suite
	productStorer *mocks.MockproductStorer
	stockStorer   *mocks.MockstockStorer
	rnd           *mrand.Rand
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) SetupSuite() {
	suite.rnd = mrand.New(mrand.NewSource(time.Now().UnixNano()))
}

func (suite *ServiceTestSuite) SetupTest() {

	ctrl := gomock.NewController(suite.T())
	suite.productStorer = mocks.NewMockproductStorer(ctrl)
	suite.stockStorer = mocks.NewMockstockStorer(ctrl)
}

func (suite *ServiceTestSuite) TearDownTest() {
}
