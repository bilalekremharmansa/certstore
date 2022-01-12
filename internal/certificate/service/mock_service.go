// Code generated by MockGen. DO NOT EDIT.
// Source: internal/certificate/service/service.go

// Package service is a generated GoMock package.
package service

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCertificateService is a mock of CertificateService interface.
type MockCertificateService struct {
	ctrl     *gomock.Controller
	recorder *MockCertificateServiceMockRecorder
}

// MockCertificateServiceMockRecorder is the mock recorder for MockCertificateService.
type MockCertificateServiceMockRecorder struct {
	mock *MockCertificateService
}

// NewMockCertificateService creates a new mock instance.
func NewMockCertificateService(ctrl *gomock.Controller) *MockCertificateService {
	mock := &MockCertificateService{ctrl: ctrl}
	mock.recorder = &MockCertificateServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCertificateService) EXPECT() *MockCertificateServiceMockRecorder {
	return m.recorder
}

// CreateCertificate mocks base method.
func (m *MockCertificateService) CreateCertificate(arg0 *NewCertificateRequest) (*NewCertificateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCertificate", arg0)
	ret0, _ := ret[0].(*NewCertificateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCertificate indicates an expected call of CreateCertificate.
func (mr *MockCertificateServiceMockRecorder) CreateCertificate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCertificate", reflect.TypeOf((*MockCertificateService)(nil).CreateCertificate), arg0)
}
