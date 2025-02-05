// Code generated by MockGen. DO NOT EDIT.
// Source: relayer/store/store.go
//
// Generated by this command:
//
//	mockgen -source=relayer/store/store.go -package mocks -destination internal/relayertest/mocks/store.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	chains "github.com/bandprotocol/falcon/relayer/chains"
	config "github.com/bandprotocol/falcon/relayer/config"
	wallet "github.com/bandprotocol/falcon/relayer/wallet"
	gomock "go.uber.org/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
	isgomock struct{}
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

// GetConfig mocks base method.
func (m *MockStore) GetConfig() (*config.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfig")
	ret0, _ := ret[0].(*config.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConfig indicates an expected call of GetConfig.
func (mr *MockStoreMockRecorder) GetConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfig", reflect.TypeOf((*MockStore)(nil).GetConfig))
}

// GetHashedPassphrase mocks base method.
func (m *MockStore) GetHashedPassphrase() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHashedPassphrase")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHashedPassphrase indicates an expected call of GetHashedPassphrase.
func (mr *MockStoreMockRecorder) GetHashedPassphrase() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHashedPassphrase", reflect.TypeOf((*MockStore)(nil).GetHashedPassphrase))
}

// HasConfig mocks base method.
func (m *MockStore) HasConfig() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasConfig")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasConfig indicates an expected call of HasConfig.
func (mr *MockStoreMockRecorder) HasConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasConfig", reflect.TypeOf((*MockStore)(nil).HasConfig))
}

// NewWallet mocks base method.
func (m *MockStore) NewWallet(chainType chains.ChainType, chainName, passphrase string) (wallet.Wallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewWallet", chainType, chainName, passphrase)
	ret0, _ := ret[0].(wallet.Wallet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewWallet indicates an expected call of NewWallet.
func (mr *MockStoreMockRecorder) NewWallet(chainType, chainName, passphrase any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewWallet", reflect.TypeOf((*MockStore)(nil).NewWallet), chainType, chainName, passphrase)
}

// SaveConfig mocks base method.
func (m *MockStore) SaveConfig(cfg *config.Config) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveConfig", cfg)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveConfig indicates an expected call of SaveConfig.
func (mr *MockStoreMockRecorder) SaveConfig(cfg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveConfig", reflect.TypeOf((*MockStore)(nil).SaveConfig), cfg)
}

// SaveHashedPassphrase mocks base method.
func (m *MockStore) SaveHashedPassphrase(hashedPassphrase []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveHashedPassphrase", hashedPassphrase)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveHashedPassphrase indicates an expected call of SaveHashedPassphrase.
func (mr *MockStoreMockRecorder) SaveHashedPassphrase(hashedPassphrase any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveHashedPassphrase", reflect.TypeOf((*MockStore)(nil).SaveHashedPassphrase), hashedPassphrase)
}
