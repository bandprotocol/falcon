// Code generated by MockGen. DO NOT EDIT.
// Source: relayer/chains/provider.go
//
// Generated by this command:
//
//	mockgen -source=relayer/chains/provider.go -package mocks -destination internal/relayertest/mocks/chain_provider.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	big "math/big"
	reflect "reflect"

	types "github.com/bandprotocol/falcon/relayer/band/types"
	types0 "github.com/bandprotocol/falcon/relayer/chains/types"
	gomock "go.uber.org/mock/gomock"
)

// MockChainProvider is a mock of ChainProvider interface.
type MockChainProvider struct {
	ctrl     *gomock.Controller
	recorder *MockChainProviderMockRecorder
	isgomock struct{}
}

// MockChainProviderMockRecorder is the mock recorder for MockChainProvider.
type MockChainProviderMockRecorder struct {
	mock *MockChainProvider
}

// NewMockChainProvider creates a new mock instance.
func NewMockChainProvider(ctrl *gomock.Controller) *MockChainProvider {
	mock := &MockChainProvider{ctrl: ctrl}
	mock.recorder = &MockChainProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChainProvider) EXPECT() *MockChainProviderMockRecorder {
	return m.recorder
}

// AddKey mocks base method.
func (m *MockChainProvider) AddKey(keyName, mnemonic, privateKeyHex, homePath string, coinType uint32, account, index uint) (*types0.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddKey", keyName, mnemonic, privateKeyHex, homePath, coinType, account, index)
	ret0, _ := ret[0].(*types0.Key)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddKey indicates an expected call of AddKey.
func (mr *MockChainProviderMockRecorder) AddKey(keyName, mnemonic, privateKeyHex, homePath, coinType, account, index any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddKey", reflect.TypeOf((*MockChainProvider)(nil).AddKey), keyName, mnemonic, privateKeyHex, homePath, coinType, account, index)
}

// DeleteKey mocks base method.
func (m *MockChainProvider) DeleteKey(homePath, keyName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteKey", homePath, keyName)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteKey indicates an expected call of DeleteKey.
func (mr *MockChainProviderMockRecorder) DeleteKey(homePath, keyName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteKey", reflect.TypeOf((*MockChainProvider)(nil).DeleteKey), homePath, keyName)
}

// ExportPrivateKey mocks base method.
func (m *MockChainProvider) ExportPrivateKey(keyName string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExportPrivateKey", keyName)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExportPrivateKey indicates an expected call of ExportPrivateKey.
func (mr *MockChainProviderMockRecorder) ExportPrivateKey(keyName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExportPrivateKey", reflect.TypeOf((*MockChainProvider)(nil).ExportPrivateKey), keyName)
}

// Init mocks base method.
func (m *MockChainProvider) Init(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockChainProviderMockRecorder) Init(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockChainProvider)(nil).Init), ctx)
}

// IsKeyNameExist mocks base method.
func (m *MockChainProvider) IsKeyNameExist(keyName string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsKeyNameExist", keyName)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsKeyNameExist indicates an expected call of IsKeyNameExist.
func (mr *MockChainProviderMockRecorder) IsKeyNameExist(keyName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsKeyNameExist", reflect.TypeOf((*MockChainProvider)(nil).IsKeyNameExist), keyName)
}

// Listkeys mocks base method.
func (m *MockChainProvider) Listkeys() []*types0.Key {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Listkeys")
	ret0, _ := ret[0].([]*types0.Key)
	return ret0
}

// Listkeys indicates an expected call of Listkeys.
func (mr *MockChainProviderMockRecorder) Listkeys() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Listkeys", reflect.TypeOf((*MockChainProvider)(nil).Listkeys))
}

// LoadFreeSenders mocks base method.
func (m *MockChainProvider) LoadFreeSenders(homePath string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadFreeSenders", homePath)
	ret0, _ := ret[0].(error)
	return ret0
}

// LoadFreeSenders indicates an expected call of LoadFreeSenders.
func (mr *MockChainProviderMockRecorder) LoadFreeSenders(homePath any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadFreeSenders", reflect.TypeOf((*MockChainProvider)(nil).LoadFreeSenders), homePath)
}

// QueryBalance mocks base method.
func (m *MockChainProvider) QueryBalance(ctx context.Context, keyName string) (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryBalance", ctx, keyName)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryBalance indicates an expected call of QueryBalance.
func (mr *MockChainProviderMockRecorder) QueryBalance(ctx, keyName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryBalance", reflect.TypeOf((*MockChainProvider)(nil).QueryBalance), ctx, keyName)
}

// QueryTunnelInfo mocks base method.
func (m *MockChainProvider) QueryTunnelInfo(ctx context.Context, tunnelID uint64, tunnelDestinationAddr string) (*types0.Tunnel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryTunnelInfo", ctx, tunnelID, tunnelDestinationAddr)
	ret0, _ := ret[0].(*types0.Tunnel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryTunnelInfo indicates an expected call of QueryTunnelInfo.
func (mr *MockChainProviderMockRecorder) QueryTunnelInfo(ctx, tunnelID, tunnelDestinationAddr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryTunnelInfo", reflect.TypeOf((*MockChainProvider)(nil).QueryTunnelInfo), ctx, tunnelID, tunnelDestinationAddr)
}

// RelayPacket mocks base method.
func (m *MockChainProvider) RelayPacket(ctx context.Context, packet *types.Packet) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RelayPacket", ctx, packet)
	ret0, _ := ret[0].(error)
	return ret0
}

// RelayPacket indicates an expected call of RelayPacket.
func (mr *MockChainProviderMockRecorder) RelayPacket(ctx, packet any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RelayPacket", reflect.TypeOf((*MockChainProvider)(nil).RelayPacket), ctx, packet)
}

// ShowKey mocks base method.
func (m *MockChainProvider) ShowKey(keyName string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShowKey", keyName)
	ret0, _ := ret[0].(string)
	return ret0
}

// ShowKey indicates an expected call of ShowKey.
func (mr *MockChainProviderMockRecorder) ShowKey(keyName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowKey", reflect.TypeOf((*MockChainProvider)(nil).ShowKey), keyName)
}

// MockKeyProvider is a mock of KeyProvider interface.
type MockKeyProvider struct {
	ctrl     *gomock.Controller
	recorder *MockKeyProviderMockRecorder
	isgomock struct{}
}

// MockKeyProviderMockRecorder is the mock recorder for MockKeyProvider.
type MockKeyProviderMockRecorder struct {
	mock *MockKeyProvider
}

// NewMockKeyProvider creates a new mock instance.
func NewMockKeyProvider(ctrl *gomock.Controller) *MockKeyProvider {
	mock := &MockKeyProvider{ctrl: ctrl}
	mock.recorder = &MockKeyProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKeyProvider) EXPECT() *MockKeyProviderMockRecorder {
	return m.recorder
}

// AddKey mocks base method.
func (m *MockKeyProvider) AddKey(keyName, mnemonic, privateKeyHex, homePath string, coinType uint32, account, index uint) (*types0.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddKey", keyName, mnemonic, privateKeyHex, homePath, coinType, account, index)
	ret0, _ := ret[0].(*types0.Key)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddKey indicates an expected call of AddKey.
func (mr *MockKeyProviderMockRecorder) AddKey(keyName, mnemonic, privateKeyHex, homePath, coinType, account, index any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddKey", reflect.TypeOf((*MockKeyProvider)(nil).AddKey), keyName, mnemonic, privateKeyHex, homePath, coinType, account, index)
}

// DeleteKey mocks base method.
func (m *MockKeyProvider) DeleteKey(homePath, keyName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteKey", homePath, keyName)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteKey indicates an expected call of DeleteKey.
func (mr *MockKeyProviderMockRecorder) DeleteKey(homePath, keyName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteKey", reflect.TypeOf((*MockKeyProvider)(nil).DeleteKey), homePath, keyName)
}

// ExportPrivateKey mocks base method.
func (m *MockKeyProvider) ExportPrivateKey(keyName string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExportPrivateKey", keyName)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExportPrivateKey indicates an expected call of ExportPrivateKey.
func (mr *MockKeyProviderMockRecorder) ExportPrivateKey(keyName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExportPrivateKey", reflect.TypeOf((*MockKeyProvider)(nil).ExportPrivateKey), keyName)
}

// IsKeyNameExist mocks base method.
func (m *MockKeyProvider) IsKeyNameExist(keyName string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsKeyNameExist", keyName)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsKeyNameExist indicates an expected call of IsKeyNameExist.
func (mr *MockKeyProviderMockRecorder) IsKeyNameExist(keyName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsKeyNameExist", reflect.TypeOf((*MockKeyProvider)(nil).IsKeyNameExist), keyName)
}

// Listkeys mocks base method.
func (m *MockKeyProvider) Listkeys() []*types0.Key {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Listkeys")
	ret0, _ := ret[0].([]*types0.Key)
	return ret0
}

// Listkeys indicates an expected call of Listkeys.
func (mr *MockKeyProviderMockRecorder) Listkeys() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Listkeys", reflect.TypeOf((*MockKeyProvider)(nil).Listkeys))
}

// LoadFreeSenders mocks base method.
func (m *MockKeyProvider) LoadFreeSenders(homePath string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadFreeSenders", homePath)
	ret0, _ := ret[0].(error)
	return ret0
}

// LoadFreeSenders indicates an expected call of LoadFreeSenders.
func (mr *MockKeyProviderMockRecorder) LoadFreeSenders(homePath any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadFreeSenders", reflect.TypeOf((*MockKeyProvider)(nil).LoadFreeSenders), homePath)
}

// ShowKey mocks base method.
func (m *MockKeyProvider) ShowKey(keyName string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShowKey", keyName)
	ret0, _ := ret[0].(string)
	return ret0
}

// ShowKey indicates an expected call of ShowKey.
func (mr *MockKeyProviderMockRecorder) ShowKey(keyName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowKey", reflect.TypeOf((*MockKeyProvider)(nil).ShowKey), keyName)
}
