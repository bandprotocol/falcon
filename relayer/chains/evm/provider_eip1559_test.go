package evm_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/internal/relayertest/mocks"
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	"github.com/bandprotocol/falcon/relayer/chains/evm"
	"github.com/bandprotocol/falcon/relayer/logger"
	"github.com/bandprotocol/falcon/relayer/wallet"
	"github.com/bandprotocol/falcon/relayer/wallet/geth"
)

type EIP1559ProviderTestSuite struct {
	suite.Suite
	ctrl *gomock.Controller

	chainProvider *evm.EVMChainProvider
	client        *mocks.MockEVMClient
	homePath      string
	chainName     string

	relayingPacket    bandtypes.Packet
	relayingCalldata  []byte
	gasInfo           evm.GasInfo
	mockSigner        wallet.Signer
	mockSignerAddress common.Address
}

func TestEIP1559ProviderTestSuite(t *testing.T) {
	suite.Run(t, new(EIP1559ProviderTestSuite))
}

func (s *EIP1559ProviderTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.client = mocks.NewMockEVMClient(s.ctrl)

	evmConfig := *baseEVMCfg
	evmConfig.GasType = evm.GasTypeEIP1559
	s.chainName = "testnet"
	s.homePath = s.T().TempDir()

	gethWallet, err := geth.NewGethWallet("", s.homePath, s.chainName)
	s.Require().NoError(err)

	log := logger.NewZapLogWrapper(zap.NewNop().Sugar())
	chainProvider, err := evm.NewEVMChainProvider(s.chainName, s.client, &evmConfig, log, gethWallet, nil)
	s.Require().NoError(err)
	s.chainProvider = chainProvider

	priv, err := crypto.HexToECDSA(evm.StripPrivateKeyPrefix(testPrivateKey))
	s.Require().NoError(err)

	s.mockSigner = geth.NewLocalSigner("test", priv)

	gethAddr, err := evm.HexToAddress(s.mockSigner.GetAddress())
	s.Require().NoError(err)
	s.mockSignerAddress = gethAddr

	s.chainProvider.FreeSigners = make(chan wallet.Signer, 1)
	s.chainProvider.FreeSigners <- s.mockSigner

	s.relayingPacket = mockPacket()
	evmSig := s.relayingPacket.CurrentGroupSigning.EVMSignature
	s.relayingCalldata, err = s.chainProvider.CreateCalldata(
		s.relayingPacket.CurrentGroupSigning.Message,
		evmSig.RAddress,
		evmSig.Signature,
	)
	s.Require().NoError(err)

	s.gasInfo = evm.NewGasEIP1559Info(big.NewInt(10_000_000_000), big.NewInt(8_000_000_000))
}

func (s *EIP1559ProviderTestSuite) MockDefaultResponses() {
	gasInfoCalldata, err := hex.DecodeString("658612e9")
	s.Require().NoError(err)
	gasInfoResponse, err := hex.DecodeString(uint256ToHex(big.NewInt(12_000_000_000)))
	s.Require().NoError(err)

	mockCtx := gomock.Any()
	s.client.EXPECT().CheckAndConnect(mockCtx).Return(nil).AnyTimes()
	s.client.EXPECT().EstimateGasTipCap(mockCtx).Return(s.gasInfo.GasPriorityFee, nil).AnyTimes()
	s.client.EXPECT().EstimateBaseFee(mockCtx).Return(s.gasInfo.GasBaseFee, nil).AnyTimes()
	s.client.EXPECT().NonceAt(mockCtx, s.mockSignerAddress).Return(uint64(100), nil).AnyTimes()
	s.client.EXPECT().
		Query(mockCtx, s.chainProvider.TunnelRouterAddress, gasInfoCalldata).
		Return(gasInfoResponse, nil).
		AnyTimes()
}

func (s *EIP1559ProviderTestSuite) TestRelayPacketSuccess() {
	// mock client responses
	s.client.EXPECT().EstimateGas(gomock.Any(), ethereum.CallMsg{
		From:      s.mockSignerAddress,
		To:        &s.chainProvider.TunnelRouterAddress,
		Data:      s.relayingCalldata,
		GasFeeCap: s.gasInfo.GasFeeCap,
		GasTipCap: s.gasInfo.GasPriorityFee,
	}).Return(uint64(200_000), nil)

	txHash := "0xabc123"
	s.client.EXPECT().BroadcastTx(gomock.Any(), gomock.Any()).Return(txHash, nil)
	s.client.EXPECT().GetTxReceipt(gomock.Any(), txHash).Return(&evm.TxReceipt{
		Status:            gethtypes.ReceiptStatusSuccessful,
		GasUsed:           21000,
		EffectiveGasPrice: big.NewInt(20000),
		BlockNumber:       big.NewInt(100),
	}, nil)

	s.client.EXPECT().GetBlockHeight(gomock.Any()).Return(uint64(105), nil)
	s.MockDefaultResponses()

	err := s.chainProvider.RelayPacket(context.Background(), &s.relayingPacket)
	s.Require().NoError(err)
}

func (s *EIP1559ProviderTestSuite) TestRelayPacketSuccessWithoutQueryMaxGasFee() {
	s.chainProvider.Config.MaxBaseFee = 2_000_000_000
	s.chainProvider.Config.MaxPriorityFee = 3_000_000_000

	// mock client responses
	s.client.EXPECT().EstimateGas(gomock.Any(), ethereum.CallMsg{
		From:      s.mockSignerAddress,
		To:        &s.chainProvider.TunnelRouterAddress,
		Data:      s.relayingCalldata,
		GasFeeCap: big.NewInt(5_000_000_000),
		GasTipCap: big.NewInt(3_000_000_000),
	}).Return(uint64(200_000), nil)

	txHash := "0xabc123"
	s.client.EXPECT().BroadcastTx(gomock.Any(), gomock.Any()).Return(txHash, nil)
	s.client.EXPECT().GetTxReceipt(gomock.Any(), txHash).Return(&evm.TxReceipt{
		Status:            gethtypes.ReceiptStatusSuccessful,
		GasUsed:           21000,
		EffectiveGasPrice: big.NewInt(20000),
		BlockNumber:       big.NewInt(100),
	}, nil)

	s.client.EXPECT().GetBlockHeight(gomock.Any()).Return(uint64(105), nil)
	s.MockDefaultResponses()

	err := s.chainProvider.RelayPacket(context.Background(), &s.relayingPacket)
	s.Require().NoError(err)
}

func (s *EIP1559ProviderTestSuite) TestRelayPacketFailedConnect() {
	// mock client responses
	s.client.EXPECT().CheckAndConnect(gomock.Any()).Return(fmt.Errorf("failed to connect client"))
	s.MockDefaultResponses()

	err := s.chainProvider.RelayPacket(context.Background(), &s.relayingPacket)
	s.Require().ErrorContains(err, "failed to connect client")
}

func (s *EIP1559ProviderTestSuite) TestRelayPacketFailedGasEstimation() {
	// mock client responses
	s.client.EXPECT().EstimateGasTipCap(gomock.Any()).Return(nil, fmt.Errorf("failed to estimate gas tip cap"))
	s.MockDefaultResponses()

	err := s.chainProvider.RelayPacket(context.Background(), &s.relayingPacket)
	s.Require().ErrorContains(err, "failed to estimate gas tip cap")
}

func (s *EIP1559ProviderTestSuite) TestRelayPacketFailedBroadcastTx() {
	// mock client responses
	s.client.EXPECT().EstimateGas(gomock.Any(), ethereum.CallMsg{
		From:      s.mockSignerAddress,
		To:        &s.chainProvider.TunnelRouterAddress,
		Data:      s.relayingCalldata,
		GasFeeCap: s.gasInfo.GasFeeCap,
		GasTipCap: s.gasInfo.GasPriorityFee,
	}).Return(uint64(200_000), nil).Times(s.chainProvider.Config.MaxRetry)

	s.client.EXPECT().
		BroadcastTx(gomock.Any(), gomock.Any()).
		Return("", fmt.Errorf("failed to broadcast an evm transaction")).
		Times(s.chainProvider.Config.MaxRetry)
	s.MockDefaultResponses()

	err := s.chainProvider.RelayPacket(context.Background(), &s.relayingPacket)
	s.Require().ErrorContains(err, "failed to relay packet after")
}

func (s *EIP1559ProviderTestSuite) TestRelayPacketFailedTxReceiptStatus() {
	// mock client responses
	s.client.EXPECT().EstimateGas(gomock.Any(), ethereum.CallMsg{
		From:      s.mockSignerAddress,
		To:        &s.chainProvider.TunnelRouterAddress,
		Data:      s.relayingCalldata,
		GasFeeCap: s.gasInfo.GasFeeCap,
		GasTipCap: s.gasInfo.GasPriorityFee,
	}).Return(uint64(200_000), nil).Times(s.chainProvider.Config.MaxRetry)

	txHash := "0xabc123"
	s.client.EXPECT().
		BroadcastTx(gomock.Any(), gomock.Any()).
		Return(txHash, nil).
		Times(s.chainProvider.Config.MaxRetry)

	s.client.EXPECT().
		GetTxReceipt(gomock.Any(), txHash).
		Return(&evm.TxReceipt{
			Status:            gethtypes.ReceiptStatusFailed,
			GasUsed:           21000,
			EffectiveGasPrice: big.NewInt(20000),
			BlockNumber:       big.NewInt(100),
		}, nil).
		Times(s.chainProvider.Config.MaxRetry)
	s.MockDefaultResponses()

	err := s.chainProvider.RelayPacket(context.Background(), &s.relayingPacket)
	s.Require().ErrorContains(err, "failed to relay packet after")
}

func (s *EIP1559ProviderTestSuite) TestBumpAndBoundGas() {
	s.MockDefaultResponses()

	// Test cases
	testCases := []struct {
		name                string
		maxPriorityFee      uint64
		maxBaseFee          uint64
		initialPriorityFee  int64
		initialBaseFee      int64
		multiplier          float64
		expectedPriorityFee int64
		expectedBaseFee     int64
	}{
		{
			name:                "Priority and base fee within limits",
			maxPriorityFee:      10_000_000_000,
			maxBaseFee:          20_000_000_000,
			initialPriorityFee:  5_000_000_000,
			initialBaseFee:      15_000_000_000,
			multiplier:          1.2,
			expectedPriorityFee: 6_000_000_000,  // due to big.Float imprecision
			expectedBaseFee:     15_000_000_000, // Unchanged
		},
		{
			name:                "Priority fee exceeds cap",
			maxPriorityFee:      8_000_000_000,
			maxBaseFee:          20_000_000_000,
			initialPriorityFee:  7_000_000_000,
			initialBaseFee:      15_000_000_000,
			multiplier:          1.2,
			expectedPriorityFee: 8_000_000_000, // Capped at maxPriorityFee
			expectedBaseFee:     15_000_000_000,
		},
		{
			name:                "Base fee exceeds cap",
			maxPriorityFee:      10_000_000_000,
			maxBaseFee:          18_000_000_000,
			initialPriorityFee:  5_000_000_000,
			initialBaseFee:      19_000_000_000,
			multiplier:          1.2,
			expectedPriorityFee: 6_000_000_000, // due to big.Float imprecision
			expectedBaseFee:     18_000_000_000,
		},
		{
			name:                "No priority fee cap, use relayer fee",
			maxPriorityFee:      0,
			maxBaseFee:          0,
			initialPriorityFee:  11_000_000_000,
			initialBaseFee:      18_000_000_000,
			multiplier:          1.2,
			expectedPriorityFee: 12_000_000_000, // due to big.Float imprecision
			expectedBaseFee:     18_000_000_000,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.chainProvider.Config.MaxPriorityFee = tc.maxPriorityFee
			s.chainProvider.Config.MaxBaseFee = tc.maxBaseFee

			actual, err := s.chainProvider.BumpAndBoundGas(
				context.Background(),
				evm.NewGasEIP1559Info(big.NewInt(tc.initialPriorityFee), big.NewInt(tc.initialBaseFee)),
				tc.multiplier,
			)
			s.Require().NoError(err)

			expected := evm.NewGasEIP1559Info(big.NewInt(tc.expectedPriorityFee), big.NewInt(tc.expectedBaseFee))
			s.Require().Equal(expected, actual, "Failed test case: %s", tc.name)
		})
	}
}

func (s *EIP1559ProviderTestSuite) TestEstimateGas() {
	s.client.EXPECT().EstimateGasTipCap(gomock.Any()).Return(big.NewInt(5_000_000_000), nil)
	s.client.EXPECT().EstimateBaseFee(gomock.Any()).Return(big.NewInt(10_000_000_000), nil)
	s.MockDefaultResponses()

	actual, err := s.chainProvider.EstimateGasFee(context.Background())
	s.Require().NoError(err)

	expected := evm.GasInfo{
		Type:           evm.GasTypeEIP1559,
		GasPrice:       nil,
		GasPriorityFee: big.NewInt(5_000_000_000),
		GasBaseFee:     big.NewInt(10_000_000_000),
		GasFeeCap:      big.NewInt(15_000_000_000),
	}

	s.Require().Equal(expected, actual)
}

func (s *EIP1559ProviderTestSuite) TestNewRelayTx() {
	data := []byte("mock calldata")

	callMsg := ethereum.CallMsg{
		From:      s.mockSignerAddress,
		To:        &s.chainProvider.TunnelRouterAddress,
		Data:      data,
		GasFeeCap: s.gasInfo.GasFeeCap,
		GasTipCap: s.gasInfo.GasPriorityFee,
	}
	s.client.EXPECT().EstimateGas(gomock.Any(), callMsg).Return(uint64(100_000), nil)
	s.client.EXPECT().NonceAt(gomock.Any(), s.mockSignerAddress).Return(uint64(1), nil)

	actual, err := s.chainProvider.NewRelayTx(context.Background(), data, s.mockSigner, s.gasInfo)
	s.Require().NoError(err)

	expected := gethtypes.NewTx(&gethtypes.DynamicFeeTx{
		ChainID:   big.NewInt(int64(s.chainProvider.Config.ChainID)),
		Nonce:     uint64(1),
		To:        &s.chainProvider.TunnelRouterAddress,
		Value:     decimal.NewFromInt(0).BigInt(),
		Data:      data,
		Gas:       100_000,
		GasFeeCap: s.gasInfo.GasFeeCap,
		GasTipCap: s.gasInfo.GasPriorityFee,
	})

	// check only some parts of the received tx.
	s.Require().Equal(expected.Nonce(), actual.Nonce(), "Nonce mismatch")
	s.Require().Equal(expected.To(), actual.To(), "To address mismatch")
	s.Require().Equal(expected.Data(), actual.Data(), "Data mismatch")
	s.Require().Equal(expected.Gas(), actual.Gas(), "Gas limit mismatch")
	s.Require().Equal(expected.GasPrice(), actual.GasPrice(), "GasPrice mismatch")
	s.Require().Equal(expected.GasTipCap(), actual.GasTipCap(), "GasTipCap mismatch")
	s.Require().Equal(expected.GasFeeCap(), actual.GasFeeCap(), "GasFeeCap mismatch")
	s.Require().Equal(expected.ChainId(), actual.ChainId(), "ChainID mismatch")
}
