// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jautyw/isa-investment-funds/internal/transport (interfaces: Service)

// Package transport is a generated GoMock package.
package transport

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	service "github.com/jautyw/isa-investment-funds/internal/service"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// GetFunds mocks base method.
func (m *MockService) GetFunds(arg0 context.Context, arg1 string) (*service.Funds, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFunds", arg0, arg1)
	ret0, _ := ret[0].(*service.Funds)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFunds indicates an expected call of GetFunds.
func (mr *MockServiceMockRecorder) GetFunds(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFunds", reflect.TypeOf((*MockService)(nil).GetFunds), arg0, arg1)
}

// GetInvestmentOverview mocks base method.
func (m *MockService) GetInvestmentOverview(arg0 context.Context, arg1 int) (*service.Overview, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInvestmentOverview", arg0, arg1)
	ret0, _ := ret[0].(*service.Overview)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInvestmentOverview indicates an expected call of GetInvestmentOverview.
func (mr *MockServiceMockRecorder) GetInvestmentOverview(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInvestmentOverview", reflect.TypeOf((*MockService)(nil).GetInvestmentOverview), arg0, arg1)
}
