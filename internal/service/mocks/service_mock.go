// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jautyw/isa-investment-funds/internal/service (interfaces: Store)

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	storage "github.com/jautyw/isa-investment-funds/internal/storage"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// GetAmountSpentCurrentTaxYear mocks base method.
func (m *MockStore) GetAmountSpentCurrentTaxYear(arg0 context.Context, arg1 int) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAmountSpentCurrentTaxYear", arg0, arg1)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAmountSpentCurrentTaxYear indicates an expected call of GetAmountSpentCurrentTaxYear.
func (mr *MockStoreMockRecorder) GetAmountSpentCurrentTaxYear(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAmountSpentCurrentTaxYear", reflect.TypeOf((*MockStore)(nil).GetAmountSpentCurrentTaxYear), arg0, arg1)
}

// GetFunds mocks base method.
func (m *MockStore) GetFunds(arg0 context.Context, arg1 string) (*storage.Funds, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFunds", arg0, arg1)
	ret0, _ := ret[0].(*storage.Funds)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFunds indicates an expected call of GetFunds.
func (mr *MockStoreMockRecorder) GetFunds(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFunds", reflect.TypeOf((*MockStore)(nil).GetFunds), arg0, arg1)
}

// GetInvestmentOverview mocks base method.
func (m *MockStore) GetInvestmentOverview(arg0 context.Context, arg1 int) ([]storage.InvestmentOverview, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInvestmentOverview", arg0, arg1)
	ret0, _ := ret[0].([]storage.InvestmentOverview)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInvestmentOverview indicates an expected call of GetInvestmentOverview.
func (mr *MockStoreMockRecorder) GetInvestmentOverview(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInvestmentOverview", reflect.TypeOf((*MockStore)(nil).GetInvestmentOverview), arg0, arg1)
}
