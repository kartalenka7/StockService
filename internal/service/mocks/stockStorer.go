// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/stock_service/stock_service.go
//
// Generated by this command:
//
//	mockgen -source=internal/service/stock_service/stock_service.go -destination=internal/service/mocks/stockStorer.go
//

// Package mock_stock_service is a generated GoMock package.
package mocks

import (
	context "context"
	model "lamoda/internal/model"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockstockStorer is a mock of stockStorer interface.
type MockstockStorer struct {
	ctrl     *gomock.Controller
	recorder *MockstockStorerMockRecorder
}

// MockstockStorerMockRecorder is the mock recorder for MockstockStorer.
type MockstockStorerMockRecorder struct {
	mock *MockstockStorer
}

// NewMockstockStorer creates a new mock instance.
func NewMockstockStorer(ctrl *gomock.Controller) *MockstockStorer {
	mock := &MockstockStorer{ctrl: ctrl}
	mock.recorder = &MockstockStorerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockstockStorer) EXPECT() *MockstockStorerMockRecorder {
	return m.recorder
}

// CheckStockAvailability mocks base method.
func (m *MockstockStorer) CheckStockAvailability(ctx context.Context, stockID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckStockAvailability", ctx, stockID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckStockAvailability indicates an expected call of CheckStockAvailability.
func (mr *MockstockStorerMockRecorder) CheckStockAvailability(ctx, stockID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckStockAvailability", reflect.TypeOf((*MockstockStorer)(nil).CheckStockAvailability), ctx, stockID)
}

// GetAvailableQty mocks base method.
func (m *MockstockStorer) GetAvailableQty(ctx context.Context, stockID int) ([]model.Products, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAvailableQty", ctx, stockID)
	ret0, _ := ret[0].([]model.Products)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAvailableQty indicates an expected call of GetAvailableQty.
func (mr *MockstockStorerMockRecorder) GetAvailableQty(ctx, stockID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvailableQty", reflect.TypeOf((*MockstockStorer)(nil).GetAvailableQty), ctx, stockID)
}