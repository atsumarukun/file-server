// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/api/usecase/auth.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	dto "file-server/internal/app/api/usecase/dto"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAuthUsecase is a mock of AuthUsecase interface.
type MockAuthUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockAuthUsecaseMockRecorder
}

// MockAuthUsecaseMockRecorder is the mock recorder for MockAuthUsecase.
type MockAuthUsecaseMockRecorder struct {
	mock *MockAuthUsecase
}

// NewMockAuthUsecase creates a new mock instance.
func NewMockAuthUsecase(ctrl *gomock.Controller) *MockAuthUsecase {
	mock := &MockAuthUsecase{ctrl: ctrl}
	mock.recorder = &MockAuthUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthUsecase) EXPECT() *MockAuthUsecaseMockRecorder {
	return m.recorder
}

// Signin mocks base method.
func (m *MockAuthUsecase) Signin(arg0 string) (*dto.AuthDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Signin", arg0)
	ret0, _ := ret[0].(*dto.AuthDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Signin indicates an expected call of Signin.
func (mr *MockAuthUsecaseMockRecorder) Signin(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Signin", reflect.TypeOf((*MockAuthUsecase)(nil).Signin), arg0)
}
