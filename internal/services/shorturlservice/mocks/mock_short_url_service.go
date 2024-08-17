// Code generated by MockGen. DO NOT EDIT.
// Source: short_url_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIGenIDForURL is a mock of IGenIDForURL interface.
type MockIGenIDForURL struct {
	ctrl     *gomock.Controller
	recorder *MockIGenIDForURLMockRecorder
}

// MockIGenIDForURLMockRecorder is the mock recorder for MockIGenIDForURL.
type MockIGenIDForURLMockRecorder struct {
	mock *MockIGenIDForURL
}

// NewMockIGenIDForURL creates a new mock instance.
func NewMockIGenIDForURL(ctrl *gomock.Controller) *MockIGenIDForURL {
	mock := &MockIGenIDForURL{ctrl: ctrl}
	mock.recorder = &MockIGenIDForURLMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIGenIDForURL) EXPECT() *MockIGenIDForURLMockRecorder {
	return m.recorder
}

// Generate mocks base method.
func (m *MockIGenIDForURL) Generate() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Generate indicates an expected call of Generate.
func (mr *MockIGenIDForURLMockRecorder) Generate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockIGenIDForURL)(nil).Generate))
}

// MockISetShortURL is a mock of ISetShortURL interface.
type MockISetShortURL struct {
	ctrl     *gomock.Controller
	recorder *MockISetShortURLMockRecorder
}

// MockISetShortURLMockRecorder is the mock recorder for MockISetShortURL.
type MockISetShortURLMockRecorder struct {
	mock *MockISetShortURL
}

// NewMockISetShortURL creates a new mock instance.
func NewMockISetShortURL(ctrl *gomock.Controller) *MockISetShortURL {
	mock := &MockISetShortURL{ctrl: ctrl}
	mock.recorder = &MockISetShortURLMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockISetShortURL) EXPECT() *MockISetShortURLMockRecorder {
	return m.recorder
}

// SetShortURL mocks base method.
func (m *MockISetShortURL) SetShortURL(ctx context.Context, idShortPath, targetURL, userID string, deleted bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetShortURL", ctx, idShortPath, targetURL, userID, deleted)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetShortURL indicates an expected call of SetShortURL.
func (mr *MockISetShortURLMockRecorder) SetShortURL(ctx, idShortPath, targetURL, userID, deleted interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetShortURL", reflect.TypeOf((*MockISetShortURL)(nil).SetShortURL), ctx, idShortPath, targetURL, userID, deleted)
}
