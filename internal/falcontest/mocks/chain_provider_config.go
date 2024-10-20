// Code generated by MockGen. DO NOT EDIT.
// Source: falcon/chains/config.go
//
// Generated by this command:
//
//	mockgen -source=falcon/chains/config.go -package mocks -destination internal/falcontest/mocks/chain_provider_config.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	chains "github.com/bandprotocol/falcon/falcon/chains"
	gomock "go.uber.org/mock/gomock"
	zap "go.uber.org/zap"
)

// MockChainProviderConfig is a mock of ChainProviderConfig interface.
type MockChainProviderConfig struct {
	ctrl     *gomock.Controller
	recorder *MockChainProviderConfigMockRecorder
	isgomock struct{}
}

// MockChainProviderConfigMockRecorder is the mock recorder for MockChainProviderConfig.
type MockChainProviderConfigMockRecorder struct {
	mock *MockChainProviderConfig
}

// NewMockChainProviderConfig creates a new mock instance.
func NewMockChainProviderConfig(ctrl *gomock.Controller) *MockChainProviderConfig {
	mock := &MockChainProviderConfig{ctrl: ctrl}
	mock.recorder = &MockChainProviderConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChainProviderConfig) EXPECT() *MockChainProviderConfigMockRecorder {
	return m.recorder
}

// NewChainProvider mocks base method.
func (m *MockChainProviderConfig) NewChainProvider(chainName string, log *zap.Logger, homePath string, debug bool) (chains.ChainProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewChainProvider", chainName, log, homePath, debug)
	ret0, _ := ret[0].(chains.ChainProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewChainProvider indicates an expected call of NewChainProvider.
func (mr *MockChainProviderConfigMockRecorder) NewChainProvider(chainName, log, homePath, debug any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewChainProvider", reflect.TypeOf((*MockChainProviderConfig)(nil).NewChainProvider), chainName, log, homePath, debug)
}

// Validate mocks base method.
func (m *MockChainProviderConfig) Validate() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate")
	ret0, _ := ret[0].(error)
	return ret0
}

// Validate indicates an expected call of Validate.
func (mr *MockChainProviderConfigMockRecorder) Validate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockChainProviderConfig)(nil).Validate))
}
