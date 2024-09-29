// Code generated by MockGen. DO NOT EDIT.
// Source: contracts.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	contracts "mmskazak/shorturl/internal/contracts"
	dtos "mmskazak/shorturl/internal/dtos"
	models "mmskazak/shorturl/internal/models"

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

// MockPinger is a mock of Pinger interface.
type MockPinger struct {
	ctrl     *gomock.Controller
	recorder *MockPingerMockRecorder
}

// MockPingerMockRecorder is the mock recorder for MockPinger.
type MockPingerMockRecorder struct {
	mock *MockPinger
}

// NewMockPinger creates a new mock instance.
func NewMockPinger(ctrl *gomock.Controller) *MockPinger {
	mock := &MockPinger{ctrl: ctrl}
	mock.recorder = &MockPingerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPinger) EXPECT() *MockPingerMockRecorder {
	return m.recorder
}

// Ping mocks base method.
func (m *MockPinger) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockPingerMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockPinger)(nil).Ping), ctx)
}

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
func (m *MockISaveBatch) SaveBatch(ctx context.Context, items []models.Incoming, baseHost, userID string, generator contracts.IGenIDForURL) ([]models.Output, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveBatch", ctx, items, baseHost, userID, generator)
	ret0, _ := ret[0].([]models.Output)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveBatch indicates an expected call of SaveBatch.
func (mr *MockISaveBatchMockRecorder) SaveBatch(ctx, items, baseHost, userID, generator interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveBatch", reflect.TypeOf((*MockISaveBatch)(nil).SaveBatch), ctx, items, baseHost, userID, generator)
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

// MockIDeleteUserURLs is a mock of IDeleteUserURLs interface.
type MockIDeleteUserURLs struct {
	ctrl     *gomock.Controller
	recorder *MockIDeleteUserURLsMockRecorder
}

// MockIDeleteUserURLsMockRecorder is the mock recorder for MockIDeleteUserURLs.
type MockIDeleteUserURLsMockRecorder struct {
	mock *MockIDeleteUserURLs
}

// NewMockIDeleteUserURLs creates a new mock instance.
func NewMockIDeleteUserURLs(ctrl *gomock.Controller) *MockIDeleteUserURLs {
	mock := &MockIDeleteUserURLs{ctrl: ctrl}
	mock.recorder = &MockIDeleteUserURLsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIDeleteUserURLs) EXPECT() *MockIDeleteUserURLsMockRecorder {
	return m.recorder
}

// DeleteURLs mocks base method.
func (m *MockIDeleteUserURLs) DeleteURLs(ctx context.Context, urlIDs []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteURLs", ctx, urlIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteURLs indicates an expected call of DeleteURLs.
func (mr *MockIDeleteUserURLsMockRecorder) DeleteURLs(ctx, urlIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteURLs", reflect.TypeOf((*MockIDeleteUserURLs)(nil).DeleteURLs), ctx, urlIDs)
}

// MockIGetUserURLs is a mock of IGetUserURLs interface.
type MockIGetUserURLs struct {
	ctrl     *gomock.Controller
	recorder *MockIGetUserURLsMockRecorder
}

// MockIGetUserURLsMockRecorder is the mock recorder for MockIGetUserURLs.
type MockIGetUserURLsMockRecorder struct {
	mock *MockIGetUserURLs
}

// NewMockIGetUserURLs creates a new mock instance.
func NewMockIGetUserURLs(ctrl *gomock.Controller) *MockIGetUserURLs {
	mock := &MockIGetUserURLs{ctrl: ctrl}
	mock.recorder = &MockIGetUserURLsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIGetUserURLs) EXPECT() *MockIGetUserURLsMockRecorder {
	return m.recorder
}

// GetUserURLs mocks base method.
func (m *MockIGetUserURLs) GetUserURLs(ctx context.Context, userID, baseHost string) ([]models.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserURLs", ctx, userID, baseHost)
	ret0, _ := ret[0].([]models.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserURLs indicates an expected call of GetUserURLs.
func (mr *MockIGetUserURLsMockRecorder) GetUserURLs(ctx, userID, baseHost interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserURLs", reflect.TypeOf((*MockIGetUserURLs)(nil).GetUserURLs), ctx, userID, baseHost)
}

// MockIGetShortURL is a mock of IGetShortURL interface.
type MockIGetShortURL struct {
	ctrl     *gomock.Controller
	recorder *MockIGetShortURLMockRecorder
}

// MockIGetShortURLMockRecorder is the mock recorder for MockIGetShortURL.
type MockIGetShortURLMockRecorder struct {
	mock *MockIGetShortURL
}

// NewMockIGetShortURL creates a new mock instance.
func NewMockIGetShortURL(ctrl *gomock.Controller) *MockIGetShortURL {
	mock := &MockIGetShortURL{ctrl: ctrl}
	mock.recorder = &MockIGetShortURLMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIGetShortURL) EXPECT() *MockIGetShortURLMockRecorder {
	return m.recorder
}

// GetShortURL mocks base method.
func (m *MockIGetShortURL) GetShortURL(ctx context.Context, idShortPath string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShortURL", ctx, idShortPath)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShortURL indicates an expected call of GetShortURL.
func (mr *MockIGetShortURLMockRecorder) GetShortURL(ctx, idShortPath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShortURL", reflect.TypeOf((*MockIGetShortURL)(nil).GetShortURL), ctx, idShortPath)
}

// MockIShortURLService is a mock of IShortURLService interface.
type MockIShortURLService struct {
	ctrl     *gomock.Controller
	recorder *MockIShortURLServiceMockRecorder
}

// MockIShortURLServiceMockRecorder is the mock recorder for MockIShortURLService.
type MockIShortURLServiceMockRecorder struct {
	mock *MockIShortURLService
}

// NewMockIShortURLService creates a new mock instance.
func NewMockIShortURLService(ctrl *gomock.Controller) *MockIShortURLService {
	mock := &MockIShortURLService{ctrl: ctrl}
	mock.recorder = &MockIShortURLServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIShortURLService) EXPECT() *MockIShortURLServiceMockRecorder {
	return m.recorder
}

// GenerateShortURL mocks base method.
func (m *MockIShortURLService) GenerateShortURL(ctx context.Context, dto dtos.DTOShortURL, generator contracts.IGenIDForURL, data contracts.ISetShortURL) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateShortURL", ctx, dto, generator, data)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateShortURL indicates an expected call of GenerateShortURL.
func (mr *MockIShortURLServiceMockRecorder) GenerateShortURL(ctx, dto, generator, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateShortURL", reflect.TypeOf((*MockIShortURLService)(nil).GenerateShortURL), ctx, dto, generator, data)
}

// MockIInternalStats is a mock of IInternalStats interface.
type MockIInternalStats struct {
	ctrl     *gomock.Controller
	recorder *MockIInternalStatsMockRecorder
}

// MockIInternalStatsMockRecorder is the mock recorder for MockIInternalStats.
type MockIInternalStatsMockRecorder struct {
	mock *MockIInternalStats
}

// NewMockIInternalStats creates a new mock instance.
func NewMockIInternalStats(ctrl *gomock.Controller) *MockIInternalStats {
	mock := &MockIInternalStats{ctrl: ctrl}
	mock.recorder = &MockIInternalStatsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIInternalStats) EXPECT() *MockIInternalStatsMockRecorder {
	return m.recorder
}

// InternalStats mocks base method.
func (m *MockIInternalStats) InternalStats(ctx context.Context) (models.Stats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InternalStats", ctx)
	ret0, _ := ret[0].(models.Stats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InternalStats indicates an expected call of InternalStats.
func (mr *MockIInternalStatsMockRecorder) InternalStats(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InternalStats", reflect.TypeOf((*MockIInternalStats)(nil).InternalStats), ctx)
}
