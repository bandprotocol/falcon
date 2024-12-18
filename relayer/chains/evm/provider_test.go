package evm

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	cmbytes "github.com/cometbft/cometbft/libs/bytes"
	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/internal/relayertest/mocks"
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	"github.com/bandprotocol/falcon/relayer/chains"
	chaintypes "github.com/bandprotocol/falcon/relayer/chains/types"
)

var baseEVMCfg = &EVMChainProviderConfig{
	BaseChainProviderConfig: chains.BaseChainProviderConfig{
		Endpoints:           []string{"http://localhost:8545"},
		ChainType:           chains.ChainTypeEVM,
		MaxRetry:            3,
		ChainID:             31337,
		TunnelRouterAddress: "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9",
		QueryTimeout:        3 * time.Second,
		ExecuteTimeout:      3 * time.Second,
	},
	BlockConfirmation:          5,
	WaitingTxDuration:          time.Second * 3,
	CheckingTxInterval:         time.Second,
	LivelinessCheckingInterval: 15 * time.Minute,
	GasMultiplier:              1,
}

type ProviderTestSuite struct {
	suite.Suite

	ctx           context.Context
	chainProvider *EVMChainProvider
	client        *mocks.MockEVMClient
	log           *zap.Logger
	homePath      string
	chainName     string
}

func TestProviderTestSuite(t *testing.T) {
	suite.Run(t, new(ProviderTestSuite))
}

// SetupTest sets up the test suite by creating a temporary directory and declare mock objects.
func (s *ProviderTestSuite) SetupTest() {
	var err error
	tmpDir := s.T().TempDir()

	ctrl := gomock.NewController(s.T())
	s.client = mocks.NewMockEVMClient(ctrl)

	log, err := zap.NewDevelopment()
	s.Require().NoError(err)

	// mock objects.
	s.log = log

	chainName := "testnet"
	s.chainName = chainName

	s.chainProvider, err = NewEVMChainProvider(s.chainName, s.client, baseEVMCfg, s.log, s.homePath)
	s.Require().NoError(err)

	s.chainProvider.Client = s.client

	s.ctx = context.Background()
	s.homePath = tmpDir
}

func (s *ProviderTestSuite) TestQueryTunnelInfo() {
	tunnelID := 1
	tunnelDestinationAddr := "0xe00F1f85abDB2aF6760759547d450da68CE66Bb1"
	addr, err := HexToAddress(tunnelDestinationAddr)
	s.Require().NoError(err)

	calldata, err := s.chainProvider.TunnelRouterABI.Pack("tunnelInfo", uint64(tunnelID), addr)
	s.Require().NoError(err)

	m, exist := s.chainProvider.TunnelRouterABI.Methods["tunnelInfo"]
	s.Require().True(exist)

	mockResponse := struct {
		IsActive       bool
		LatestSequence uint64
		Balance        *big.Int
	}{
		IsActive:       true,
		LatestSequence: 42,
		Balance:        big.NewInt(1000),
	}

	b, err := m.Outputs.Pack(mockResponse)
	s.Require().NoError(err)

	s.client.EXPECT().Query(s.ctx, s.chainProvider.TunnelRouterAddress, calldata).Return(b, nil)
	s.client.EXPECT().CheckAndConnect(s.ctx).Return(nil)

	actual, err := s.chainProvider.QueryTunnelInfo(s.ctx, uint64(tunnelID), tunnelDestinationAddr)
	s.Require().NoError(err)

	expected := &chaintypes.Tunnel{
		ID:             uint64(tunnelID),
		TargetAddress:  tunnelDestinationAddr,
		IsActive:       true,
		LatestSequence: 42,
		Balance:        big.NewInt(1000),
	}

	s.Require().Equal(expected, actual)
}

func (s *ProviderTestSuite) TestQueryTunnelInfoFailedConnection() {
	tunnelID := 1
	tunnelDestinationAddr := "invalid-address"

	s.client.EXPECT().CheckAndConnect(s.ctx).Return(fmt.Errorf("Connect client error"))

	_, err := s.chainProvider.QueryTunnelInfo(s.ctx, uint64(tunnelID), tunnelDestinationAddr)
	s.Require().ErrorContains(err, "[EVMProvider] failed to connect client:")
}

func (s *ProviderTestSuite) TestQueryTunnelInfoInvalidAddress() {
	tunnelID := 1
	tunnelDestinationAddr := "invalid-address"

	s.client.EXPECT().CheckAndConnect(s.ctx).Return(nil)

	_, err := s.chainProvider.QueryTunnelInfo(s.ctx, uint64(tunnelID), tunnelDestinationAddr)
	s.Require().ErrorContains(err, "[EVMProvider] invalid address:")
}

func (s *ProviderTestSuite) TestRelayPacket() {
	var err error
	gasEVMCfg := *baseEVMCfg
	gasEVMCfg.GasType = GasTypeEIP1559
	gasEVMCfg.GasLimit = 30000
	gasEVMCfg.MaxPriorityFee = 60
	s.chainProvider, err = NewEVMChainProvider(s.chainName, s.client, &gasEVMCfg, s.log, s.homePath)
	s.Require().NoError(err)

	// Create a mock EVMSignature
	evmSignature := bandtypes.NewEVMSignature(
		gethcommon.HexToAddress("0xfad9c8855b740a0b7ed4c221dbad0f33a83a49ca").Bytes(),
		cmbytes.HexBytes("0xabcd"),
	)

	// Create mock signing information
	signingInfo := bandtypes.NewSigning(
		1,
		cmbytes.HexBytes("0xdeadbeef"),
		evmSignature,
		"SIGNING_STATUS_SUCCESS",
	)

	packet := &bandtypes.Packet{
		TunnelID: 1,
		Sequence: 42,
		SignalPrices: []bandtypes.SignalPrice{
			{SignalID: "signal1", Price: 100},
			{SignalID: "signal2", Price: 200},
		},
		CurrentGroupSigning:  signingInfo,
		IncomingGroupSigning: nil,
	}

	addr, err := HexToAddress(testAddress)
	s.Require().NoError(err)
	priv, err := crypto.HexToECDSA(StripPrivateKeyPrefix(testPrivateKey))
	s.Require().NoError(err)

	// Mock sender
	mockSender := &Sender{
		Address:    addr,
		PrivateKey: priv,
	}

	s.chainProvider.FreeSenders = make(chan *Sender, 1)
	s.chainProvider.FreeSenders <- mockSender

	receipt := &gethtypes.Receipt{
		Status:      gethtypes.ReceiptStatusSuccessful, // Mock a successful transaction
		GasUsed:     21000,
		BlockNumber: big.NewInt(100),
	}
	latestBlock := uint64(105) // Mock the latest block height

	// Mock client responses
	s.client.EXPECT().EstimateGasTipCap(s.ctx).Return(big.NewInt(50), nil)
	s.client.EXPECT().EstimateBaseFee(s.ctx).Return(big.NewInt(100), nil)
	s.client.EXPECT().CheckAndConnect(s.ctx).Return(nil)
	s.client.EXPECT().PendingNonceAt(s.ctx, addr).Return(uint64(100), nil)
	s.client.EXPECT().BroadcastTx(s.ctx, gomock.Any()).Return("0xabc123", nil)
	s.client.EXPECT().GetTxReceipt(s.ctx, "0xabc123").Return(receipt, nil)
	s.client.EXPECT().GetBlockHeight(s.ctx).Return(latestBlock, nil)

	// Call the method
	err = s.chainProvider.RelayPacket(s.ctx, packet)
	s.Require().NoError(err)
}

func (s *ProviderTestSuite) TestRelayPacketFailedConnection() {
	var err error
	gasEVMCfg := *baseEVMCfg
	gasEVMCfg.GasType = GasTypeEIP1559
	gasEVMCfg.GasLimit = 30000
	gasEVMCfg.MaxPriorityFee = 60
	s.chainProvider, err = NewEVMChainProvider(s.chainName, s.client, &gasEVMCfg, s.log, s.homePath)
	s.Require().NoError(err)

	// Create a mock EVMSignature
	evmSignature := bandtypes.NewEVMSignature(
		gethcommon.HexToAddress("0xfad9c8855b740a0b7ed4c221dbad0f33a83a49ca").Bytes(),
		cmbytes.HexBytes("0xabcd"),
	)

	// Create mock signing information
	signingInfo := bandtypes.NewSigning(
		1,
		cmbytes.HexBytes("0xdeadbeef"),
		evmSignature,
		"SIGNING_STATUS_SUCCESS",
	)

	packet := &bandtypes.Packet{
		TunnelID: 1,
		Sequence: 42,
		SignalPrices: []bandtypes.SignalPrice{
			{SignalID: "signal1", Price: 100},
			{SignalID: "signal2", Price: 200},
		},
		CurrentGroupSigning:  signingInfo,
		IncomingGroupSigning: nil,
	}

	s.client.EXPECT().CheckAndConnect(s.ctx).Return(fmt.Errorf("failed to connect client"))

	err = s.chainProvider.RelayPacket(s.ctx, packet)
	s.Require().ErrorContains(err, "failed to connect client")
}

func (s *ProviderTestSuite) TestRelayPacketFailedGasEstimation() {
	var err error
	gasEVMCfg := *baseEVMCfg
	gasEVMCfg.GasType = GasTypeEIP1559
	gasEVMCfg.GasLimit = 30000
	gasEVMCfg.MaxPriorityFee = 60

	s.chainProvider, err = NewEVMChainProvider(s.chainName, s.client, &gasEVMCfg, s.log, s.homePath)
	s.Require().NoError(err)

	// Create a mock EVMSignature
	evmSignature := bandtypes.NewEVMSignature(
		gethcommon.HexToAddress("0xfad9c8855b740a0b7ed4c221dbad0f33a83a49ca").Bytes(),
		cmbytes.HexBytes("0xabcd"),
	)

	// Create mock signing information
	signingInfo := bandtypes.NewSigning(
		1,
		cmbytes.HexBytes("0xdeadbeef"),
		evmSignature,
		"SIGNING_STATUS_SUCCESS",
	)

	packet := &bandtypes.Packet{
		TunnelID: 1,
		Sequence: 42,
		SignalPrices: []bandtypes.SignalPrice{
			{SignalID: "signal1", Price: 100},
			{SignalID: "signal2", Price: 200},
		},
		CurrentGroupSigning:  signingInfo,
		IncomingGroupSigning: nil,
	}

	s.client.EXPECT().CheckAndConnect(s.ctx).Return(nil)
	s.client.EXPECT().EstimateGasTipCap(s.ctx).Return(nil, fmt.Errorf("failed to estimate tip cap"))

	err = s.chainProvider.RelayPacket(s.ctx, packet)
	s.Require().ErrorContains(err, "failed to estimate gas")
}

func (s *ProviderTestSuite) TestRelayPacketFailedHandleRelay() {
	var err error
	gasEVMCfg := *baseEVMCfg
	gasEVMCfg.GasType = GasTypeEIP1559
	gasEVMCfg.GasLimit = 30000
	gasEVMCfg.MaxPriorityFee = 60
	s.chainProvider, err = NewEVMChainProvider(s.chainName, s.client, &gasEVMCfg, s.log, s.homePath)
	s.Require().NoError(err)

	// Create a mock EVMSignature
	evmSignature := bandtypes.NewEVMSignature(
		gethcommon.HexToAddress("0xfad9c8855b740a0b7ed4c221dbad0f33a83a49ca").Bytes(),
		cmbytes.HexBytes("0xabcd"),
	)

	// Create mock signing information
	signingInfo := bandtypes.NewSigning(
		1,
		cmbytes.HexBytes("0xdeadbeef"),
		evmSignature,
		"SIGNING_STATUS_SUCCESS",
	)

	packet := &bandtypes.Packet{
		TunnelID: 1,
		Sequence: 42,
		SignalPrices: []bandtypes.SignalPrice{
			{SignalID: "signal1", Price: 100},
			{SignalID: "signal2", Price: 200},
		},
		CurrentGroupSigning:  signingInfo,
		IncomingGroupSigning: nil,
	}

	addr, err := HexToAddress(testAddress)
	s.Require().NoError(err)
	priv, err := crypto.HexToECDSA(StripPrivateKeyPrefix(testPrivateKey))
	s.Require().NoError(err)

	// Mock sender
	mockSender := &Sender{
		Address:    addr,
		PrivateKey: priv,
	}

	s.chainProvider.FreeSenders = make(chan *Sender, 1)
	s.chainProvider.FreeSenders <- mockSender

	// Mock client responses
	s.client.EXPECT().EstimateGasTipCap(s.ctx).Return(big.NewInt(50), nil)
	s.client.EXPECT().EstimateBaseFee(s.ctx).Return(big.NewInt(100), nil)
	s.client.EXPECT().CheckAndConnect(s.ctx).Return(nil)
	s.client.EXPECT().PendingNonceAt(s.ctx, addr).Return(uint64(100), nil).Times(s.chainProvider.Config.MaxRetry)
	s.client.EXPECT().
		BroadcastTx(s.ctx, gomock.Any()).
		Return("", fmt.Errorf("failed to broadcast an evm transaction")).
		Times(s.chainProvider.Config.MaxRetry)

	// Call the method
	err = s.chainProvider.RelayPacket(s.ctx, packet)
	s.Require().ErrorContains(err, "failed to relay packet after")
}

func (s *ProviderTestSuite) TestRelayPacketCheckTxFailed() {
	var err error
	gasEVMCfg := *baseEVMCfg
	gasEVMCfg.GasType = GasTypeEIP1559
	gasEVMCfg.GasLimit = 30000
	gasEVMCfg.MaxPriorityFee = 60
	s.chainProvider, err = NewEVMChainProvider(s.chainName, s.client, &gasEVMCfg, s.log, s.homePath)
	s.Require().NoError(err)

	// Create a mock EVMSignature
	evmSignature := bandtypes.NewEVMSignature(
		gethcommon.HexToAddress("0xfad9c8855b740a0b7ed4c221dbad0f33a83a49ca").Bytes(),
		cmbytes.HexBytes("0xabcd"),
	)

	// Create mock signing information
	signingInfo := bandtypes.NewSigning(
		1,
		cmbytes.HexBytes("0xdeadbeef"),
		evmSignature,
		"SIGNING_STATUS_SUCCESS",
	)

	packet := &bandtypes.Packet{
		TunnelID: 1,
		Sequence: 42,
		SignalPrices: []bandtypes.SignalPrice{
			{SignalID: "signal1", Price: 100},
			{SignalID: "signal2", Price: 200},
		},
		CurrentGroupSigning:  signingInfo,
		IncomingGroupSigning: nil,
	}

	addr, err := HexToAddress(testAddress)
	s.Require().NoError(err)
	priv, err := crypto.HexToECDSA(StripPrivateKeyPrefix(testPrivateKey))
	s.Require().NoError(err)

	// Mock sender
	mockSender := &Sender{
		Address:    addr,
		PrivateKey: priv,
	}

	s.chainProvider.FreeSenders = make(chan *Sender, 1)
	s.chainProvider.FreeSenders <- mockSender

	receipt := &gethtypes.Receipt{
		Status:      gethtypes.ReceiptStatusFailed,
		GasUsed:     21000,
		BlockNumber: big.NewInt(100),
	}

	// Mock client responses
	s.client.EXPECT().EstimateGasTipCap(s.ctx).Return(big.NewInt(50), nil)
	s.client.EXPECT().EstimateBaseFee(s.ctx).Return(big.NewInt(100), nil)
	s.client.EXPECT().CheckAndConnect(s.ctx).Return(nil)
	s.client.EXPECT().PendingNonceAt(s.ctx, addr).Return(uint64(100), nil).Times(s.chainProvider.Config.MaxRetry)
	s.client.EXPECT().BroadcastTx(s.ctx, gomock.Any()).Return("0xabc123", nil).Times(s.chainProvider.Config.MaxRetry)
	s.client.EXPECT().GetTxReceipt(s.ctx, "0xabc123").Return(receipt, nil).Times(s.chainProvider.Config.MaxRetry)

	// Call the method
	err = s.chainProvider.RelayPacket(s.ctx, packet)
	s.Require().ErrorContains(err, "failed to relay packet after")
}

func (s *ProviderTestSuite) TestEstimateGasLegacy() {
	var err error
	gasEVMCfg := *baseEVMCfg
	gasEVMCfg.GasType = GasTypeLegacy
	gasEVMCfg.MaxGasPrice = 120

	s.chainProvider, err = NewEVMChainProvider(s.chainName, s.client, &gasEVMCfg, s.log, s.homePath)
	s.Require().NoError(err)

	s.client.EXPECT().EstimateGasPrice(s.ctx).Return(big.NewInt(100), nil)

	actual, err := s.chainProvider.EstimateGas(s.ctx)
	s.Require().NoError(err)

	expected := GasInfo{
		Type:           GasTypeLegacy,
		GasPrice:       big.NewInt(100),
		GasPriorityFee: nil,
		GasBaseFee:     nil,
	}

	s.Require().Equal(expected, actual)
}

func (s *ProviderTestSuite) TestEstimateGasEIP1559() {
	var err error
	gasEVMCfg := *baseEVMCfg
	gasEVMCfg.GasType = GasTypeEIP1559
	gasEVMCfg.MaxPriorityFee = 60

	s.chainProvider, err = NewEVMChainProvider(s.chainName, s.client, &gasEVMCfg, s.log, s.homePath)
	s.Require().NoError(err)

	s.client.EXPECT().EstimateGasTipCap(s.ctx).Return(big.NewInt(50), nil)
	s.client.EXPECT().EstimateBaseFee(s.ctx).Return(big.NewInt(100), nil)

	actual, err := s.chainProvider.EstimateGas(s.ctx)
	s.Require().NoError(err)

	expected := GasInfo{
		Type:           GasTypeEIP1559,
		GasPrice:       nil,
		GasPriorityFee: big.NewInt(50),
		GasBaseFee:     big.NewInt(100),
	}

	s.Require().Equal(expected, actual)
}

func (s *ProviderTestSuite) TestEstimateGasUnsupportedGas() {
	_, err := s.chainProvider.EstimateGas(s.ctx)
	s.Require().ErrorContains(err, "unsupported gas type:")
}

func (s *ProviderTestSuite) TestBumpAndBoundGasLegacy() {
	var err error

	// Test cases
	testCases := []struct {
		name             string
		maxGasPrice      uint64
		initialGasPrice  int64
		multiplier       float64
		mockRelayerFee   *big.Int
		expectedGasPrice int64
	}{
		{
			name:             "Gas price within limit",
			maxGasPrice:      150,
			initialGasPrice:  100,
			multiplier:       1.2,
			expectedGasPrice: 120, // due to big.Float imprecision
		},
		{
			name:             "Gas price exceeding limit",
			maxGasPrice:      150,
			initialGasPrice:  140,
			multiplier:       1.2,
			expectedGasPrice: 150, // Capped at maxGasPrice
		},
		{
			name:             "No gas price cap, use relayer fee",
			maxGasPrice:      0,
			initialGasPrice:  140,
			multiplier:       1.2,
			mockRelayerFee:   big.NewInt(160),
			expectedGasPrice: 160, // Use mocked relayer fee
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			gasEVMCfg := *baseEVMCfg
			gasEVMCfg.GasType = GasTypeLegacy
			gasEVMCfg.MaxGasPrice = tc.maxGasPrice

			s.chainProvider, err = NewEVMChainProvider(s.chainName, s.client, &gasEVMCfg, s.log, s.homePath)
			s.Require().NoError(err)

			// Mock relayer fee if applicable
			if tc.mockRelayerFee != nil {
				calldata, err := s.chainProvider.TunnelRouterABI.Pack("gasFee")
				s.Require().NoError(err)

				m, exist := s.chainProvider.TunnelRouterABI.Methods["gasFee"]
				s.Require().True(exist)

				b, err := m.Outputs.Pack(tc.mockRelayerFee)
				s.Require().NoError(err)

				s.client.EXPECT().
					Query(s.ctx, s.chainProvider.TunnelRouterAddress, calldata).
					Return(b, nil)
			}

			actual, err := s.chainProvider.BumpAndBoundGas(
				s.ctx,
				NewGasLegacyInfo(big.NewInt(tc.initialGasPrice)),
				tc.multiplier,
			)
			s.Require().NoError(err)

			expected := GasInfo{
				Type:           GasTypeLegacy,
				GasPrice:       big.NewInt(tc.expectedGasPrice),
				GasPriorityFee: nil,
				GasBaseFee:     nil,
			}

			s.Require().Equal(expected, actual, "Failed test case: %s", tc.name)
		})
	}
}

func (s *ProviderTestSuite) TestBumpAndBoundGasEIP1559() {
	var err error

	// Test cases
	testCases := []struct {
		name                string
		maxPriorityFee      uint64
		maxBaseFee          uint64
		initialPriorityFee  int64
		initialBaseFee      int64
		multiplier          float64
		mockRelayerFee      *big.Int
		expectedPriorityFee int64
		expectedBaseFee     int64
	}{
		{
			name:                "Priority and base fee within limits",
			maxPriorityFee:      100,
			maxBaseFee:          200,
			initialPriorityFee:  50,
			initialBaseFee:      150,
			multiplier:          1.2,
			expectedPriorityFee: 60,  // due to big.Float imprecision
			expectedBaseFee:     150, // Unchanged
		},
		{
			name:                "Priority fee exceeds cap",
			maxPriorityFee:      80,
			maxBaseFee:          200,
			initialPriorityFee:  70,
			initialBaseFee:      150,
			multiplier:          1.2,
			expectedPriorityFee: 80, // Capped at maxPriorityFee
			expectedBaseFee:     150,
		},
		{
			name:                "Base fee exceeds cap",
			maxPriorityFee:      100,
			maxBaseFee:          180,
			initialPriorityFee:  50,
			initialBaseFee:      190,
			multiplier:          1.2,
			expectedPriorityFee: 60, // due to big.Float imprecision
			expectedBaseFee:     180,
		},
		{
			name:                "No priority fee cap, use relayer fee",
			maxPriorityFee:      0,
			maxBaseFee:          200,
			initialPriorityFee:  70,
			initialBaseFee:      150,
			multiplier:          1.2,
			mockRelayerFee:      big.NewInt(90),
			expectedPriorityFee: 84, // due to big.Float imprecision
			expectedBaseFee:     150,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			gasEVMCfg := *baseEVMCfg
			gasEVMCfg.GasType = GasTypeEIP1559
			gasEVMCfg.MaxPriorityFee = tc.maxPriorityFee
			gasEVMCfg.MaxBaseFee = tc.maxBaseFee

			s.chainProvider, err = NewEVMChainProvider(s.chainName, s.client, &gasEVMCfg, s.log, s.homePath)
			s.Require().NoError(err)

			// Mock relayer fee if applicable
			if tc.mockRelayerFee != nil {
				calldata, err := s.chainProvider.TunnelRouterABI.Pack("gasFee")
				s.Require().NoError(err)

				m, exist := s.chainProvider.TunnelRouterABI.Methods["gasFee"]
				s.Require().True(exist)

				b, err := m.Outputs.Pack(tc.mockRelayerFee)
				s.Require().NoError(err)

				s.client.EXPECT().
					Query(s.ctx, s.chainProvider.TunnelRouterAddress, calldata).
					Return(b, nil)
			}

			actual, err := s.chainProvider.BumpAndBoundGas(
				s.ctx,
				NewGasEIP1559Info(big.NewInt(tc.initialPriorityFee), big.NewInt(tc.initialBaseFee)),
				tc.multiplier,
			)
			s.Require().NoError(err)

			expected := GasInfo{
				Type:           GasTypeEIP1559,
				GasPrice:       nil,
				GasPriorityFee: big.NewInt(tc.expectedPriorityFee),
				GasBaseFee:     big.NewInt(tc.expectedBaseFee),
			}

			s.Require().Equal(expected, actual, "Failed test case: %s", tc.name)
		})
	}
}

func (s *ProviderTestSuite) TestCheckConfirmedTx() {
	txHash := "0xabc123"
	receipt := &gethtypes.Receipt{
		Status:      gethtypes.ReceiptStatusSuccessful,
		GasUsed:     21000,
		BlockNumber: big.NewInt(100),
	}
	latestBlock := uint64(105)

	// Mock client behavior
	s.client.EXPECT().GetTxReceipt(s.ctx, txHash).Return(receipt, nil)
	s.client.EXPECT().GetBlockHeight(s.ctx).Return(latestBlock, nil)

	// Call checkConfirmedTx
	result, err := s.chainProvider.checkConfirmedTx(s.ctx, txHash)
	s.Require().NoError(err)

	// Define expected result
	expected := NewConfirmTxResult(
		txHash,
		TX_STATUS_SUCCESS,
		decimal.NewNullDecimal(decimal.New(21000, 0)), // Gas used
		s.chainProvider.GasType,
	)

	// Assert the results
	s.Require().Equal(expected.Status, result.Status)
	s.Require().Equal(expected.GasUsed, result.GasUsed)
}

func (s *ProviderTestSuite) TestCheckConfirmedTxReceiptFailed() {
	txHash := "0xabc123"
	receipt := &gethtypes.Receipt{
		Status:      gethtypes.ReceiptStatusFailed,
		GasUsed:     21000,
		BlockNumber: big.NewInt(100),
	}

	// Mock client behavior
	s.client.EXPECT().GetTxReceipt(s.ctx, txHash).Return(receipt, nil)

	// Call checkConfirmedTx
	result, err := s.chainProvider.checkConfirmedTx(s.ctx, txHash)
	s.Require().NoError(err)

	// Define expected result
	expected := NewConfirmTxResult(
		txHash,
		TX_STATUS_FAILED,
		decimal.NullDecimal{},
		s.chainProvider.GasType,
	)

	// Assert the results
	s.Require().Equal(expected.Status, result.Status)
	s.Require().Equal(expected.GasUsed, result.GasUsed)
}

func (s *ProviderTestSuite) TestCheckConfirmedTxReceiptUnmined() {
	txHash := "0xabc123"
	receipt := &gethtypes.Receipt{
		Status:      gethtypes.ReceiptStatusSuccessful,
		GasUsed:     21000,
		BlockNumber: big.NewInt(100),
	}
	latestBlock := uint64(103)

	// Mock client behavior
	s.client.EXPECT().GetTxReceipt(s.ctx, txHash).Return(receipt, nil)
	s.client.EXPECT().GetBlockHeight(s.ctx).Return(latestBlock, nil)

	// Call checkConfirmedTx
	result, err := s.chainProvider.checkConfirmedTx(s.ctx, txHash)
	s.Require().NoError(err)

	// Define expected result
	expected := NewConfirmTxResult(
		txHash,
		TX_STATUS_UNMINED,
		decimal.NullDecimal{},
		s.chainProvider.GasType,
	)

	// Assert the results
	s.Require().Equal(expected.Status, result.Status)
	s.Require().Equal(expected.GasUsed, result.GasUsed)
}

func (s *ProviderTestSuite) TestNewRelayTxLegacy() {
	var err error
	gasEVMCfg := *baseEVMCfg
	gasEVMCfg.GasType = GasTypeLegacy

	fmt.Println("gasEVMCfg.GasLimit", gasEVMCfg.GasLimit)

	s.chainProvider, err = NewEVMChainProvider(s.chainName, s.client, &gasEVMCfg, s.log, s.homePath)
	s.Require().NoError(err)

	data := []byte("mock calldata")
	sender, err := HexToAddress(testAddress)
	s.Require().NoError(err)

	gasInfo := NewGasLegacyInfo(big.NewInt(1000000000))

	callMsg := ethereum.CallMsg{
		From: sender,
		To:   &s.chainProvider.TunnelRouterAddress,
		Data: data,
	}

	s.client.EXPECT().EstimateGas(s.ctx, callMsg).Return(uint64(100), nil)
	s.client.EXPECT().PendingNonceAt(s.ctx, sender).Return(uint64(1), nil)

	actual, err := s.chainProvider.newRelayTx(s.ctx, data, sender, gasInfo)
	s.Require().NoError(err)

	expected := gethtypes.NewTx(&gethtypes.LegacyTx{
		Nonce:    uint64(1),
		To:       &s.chainProvider.TunnelRouterAddress,
		Value:    decimal.NewFromInt(0).BigInt(),
		Data:     data,
		Gas:      uint64(100),
		GasPrice: gasInfo.GasPrice,
	})

	s.Require().Equal(expected.Nonce(), actual.Nonce(), "Nonce mismatch")
	s.Require().Equal(expected.To(), actual.To(), "To address mismatch")
	s.Require().Equal(expected.Data(), actual.Data(), "Data mismatch")
	s.Require().Equal(expected.Gas(), actual.Gas(), "Gas limit mismatch")
	s.Require().Equal(expected.GasPrice(), actual.GasPrice(), "GasPrice mismatch")
	s.Require().Equal(expected.GasTipCap(), actual.GasTipCap(), "GasTipCap mismatch")
	s.Require().Equal(expected.GasFeeCap(), actual.GasFeeCap(), "GasFeeCap mismatch")
	s.Require().Equal(expected.ChainId(), actual.ChainId(), "ChainID mismatch")
}

func (s *ProviderTestSuite) TestNewRelayTxEIP1559() {
	var err error
	gasEVMCfg := *baseEVMCfg
	gasEVMCfg.GasLimit = 100000
	gasEVMCfg.GasType = GasTypeEIP1559

	s.chainProvider, err = NewEVMChainProvider(s.chainName, s.client, &gasEVMCfg, s.log, s.homePath)
	s.Require().NoError(err)

	data := []byte("mock calldata")
	sender, err := HexToAddress(testAddress)
	s.Require().NoError(err)

	gasInfo := NewGasEIP1559Info(big.NewInt(2000000000), big.NewInt(1000000000))

	s.client.EXPECT().PendingNonceAt(s.ctx, sender).Return(uint64(1), nil)

	actual, err := s.chainProvider.newRelayTx(s.ctx, data, sender, gasInfo)
	s.Require().NoError(err)

	expected := gethtypes.NewTx(&gethtypes.DynamicFeeTx{
		ChainID:   big.NewInt(int64(s.chainProvider.Config.ChainID)),
		Nonce:     uint64(1),
		To:        &s.chainProvider.TunnelRouterAddress,
		Value:     decimal.NewFromInt(0).BigInt(),
		Data:      data,
		Gas:       100000,
		GasFeeCap: big.NewInt(3000000000), // Sum of GasBaseFee and GasPriorityFee
		GasTipCap: gasInfo.GasPriorityFee,
	})

	s.Require().Equal(expected.Nonce(), actual.Nonce(), "Nonce mismatch")
	s.Require().Equal(expected.To(), actual.To(), "To address mismatch")
	s.Require().Equal(expected.Data(), actual.Data(), "Data mismatch")
	s.Require().Equal(expected.Gas(), actual.Gas(), "Gas limit mismatch")
	s.Require().Equal(expected.GasPrice(), actual.GasPrice(), "GasPrice mismatch")
	s.Require().Equal(expected.GasTipCap(), actual.GasTipCap(), "GasTipCap mismatch")
	s.Require().Equal(expected.GasFeeCap(), actual.GasFeeCap(), "GasFeeCap mismatch")
	s.Require().Equal(expected.ChainId(), actual.ChainId(), "ChainID mismatch")
}
