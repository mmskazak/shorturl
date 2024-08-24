// Code generated by MockGen. DO NOT EDIT.
// Source: save_batch_urls.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	storage "mmskazak/shorturl/internal/storage"

	gomock "github.com/golang/mock/gomock"
)

// MockISaveBatch is a mock of ISaveBatch interface.
type MockISaveBatch struct {
	ctrl     *gomock.Controller
	recorder *MockISaveBatchMockRecorder
}

// MockISaveBatchMockRecorder is the mock recorder for MockISaveBatch.
type MockISaveBatchMockRecorder struct {
	mock *MockISaveBatch
}

// NewMockISaveBatch creates a new mock instance.
func NewMockISaveBatch(ctrl *gomock.Controller) *MockISaveBatch {
	mock := &MockISaveBatch{ctrl: ctrl}
	mock.recorder = &MockISaveBatchMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockISaveBatch) EXPECT() *MockISaveBatchMockRecorder {
	return m.recorder
}

// SaveBatch mocks base method.
func (m *MockISaveBatch) SaveBatch(ctx context.Context, items []storage.Incoming, baseHost, userID string, generator storage.IGenIDForURL) ([]storage.Output, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveBatch", ctx, items, baseHost, userID, generator)
	ret0, _ := ret[0].([]storage.Output)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveBatch indicates an expected call of SaveBatch.
func (mr *MockISaveBatchMockRecorder) SaveBatch(ctx, items, baseHost, userID, generator interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveBatch", reflect.TypeOf((*MockISaveBatch)(nil).SaveBatch), ctx, items, baseHost, userID, generator)
}
