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

type LegacyProviderTestSuite struct {
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

func TestLegacyProviderTestSuite(t *testing.T) {
	suite.Run(t, new(LegacyProviderTestSuite))
}

func (s *LegacyProviderTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.client = mocks.NewMockEVMClient(s.ctrl)

	evmConfig := *baseEVMCfg
	evmConfig.GasType = evm.GasTypeLegacy
	s.chainName = "testnet"
	s.homePath = s.T().TempDir()

	gethWallet, err := geth.NewGethWallet("", s.homePath, s.chainName)
	s.Require().NoError(err)

	log := logger.NewZapLogWrapper(zap.NewNop())
	chainProvider, err := evm.NewEVMChainProvider(s.chainName, s.client, &evmConfig, log, gethWallet)
	s.Require().NoError(err)
	s.chainProvider = chainProvider

	priv, err := crypto.HexToECDSA(evm.StripPrivateKeyPrefix(testPrivateKey))
	s.Require().NoError(err)

	s.mockSigner = geth.NewLocalSigner("testkey", priv)

	gethAddr, err := evm.HexToAddress(s.mockSigner.GetAddress())
	s.Require().NoError(err)
	s.mockSignerAddress = gethAddr

	s.chainProvider.FreeSigners = make(chan wallet.Signer, 1)
	s.chainProvider.FreeSigners <- s.mockSigner

	s.relayingPacket = mockPacket()
	s.relayingCalldata, err = s.chainProvider.CreateCalldata(&s.relayingPacket)
	s.Require().NoError(err)

	s.gasInfo = evm.NewGasLegacyInfo(big.NewInt(10_000_000_000))
}

func (s *LegacyProviderTestSuite) MockDefaultResponses() {
	gasInfoCalldata, err := hex.DecodeString("658612e9")
	s.Require().NoError(err)
	gasInfoResponse, err := hex.DecodeString(uint256ToHex(big.NewInt(12_000_000_000)))
	s.Require().NoError(err)

	mockCtx := gomock.Any()
	s.client.EXPECT().CheckAndConnect(mockCtx).Return(nil).AnyTimes()
	s.client.EXPECT().EstimateGasPrice(mockCtx).Return(s.gasInfo.GasPrice, nil).AnyTimes()
	s.client.EXPECT().NonceAt(mockCtx, s.mockSignerAddress).Return(uint64(100), nil).AnyTimes()
	s.client.EXPECT().
		Query(mockCtx, s.chainProvider.TunnelRouterAddress, gasInfoCalldata).
		Return(gasInfoResponse, nil).
		AnyTimes()
}

func (s *LegacyProviderTestSuite) TestRelayPacketSuccess() {
	// mock client responses
	s.client.EXPECT().EstimateGas(gomock.Any(), ethereum.CallMsg{
		From:     s.mockSignerAddress,
		To:       &s.chainProvider.TunnelRouterAddress,
		Data:     s.relayingCalldata,
		GasPrice: s.gasInfo.GasPrice,
	}).Return(uint64(200_000), nil).AnyTimes()
	txHash := "0xabc123"
	s.client.EXPECT().BroadcastTx(gomock.Any(), gomock.Any()).Return(txHash, nil)
	s.client.EXPECT().GetTxReceipt(gomock.Any(), txHash).Return(&gethtypes.Receipt{
		Status:      gethtypes.ReceiptStatusSuccessful,
		GasUsed:     21000,
		BlockNumber: big.NewInt(100),
	}, nil)

	s.client.EXPECT().GetBlockHeight(gomock.Any()).Return(uint64(105), nil)
	s.MockDefaultResponses()

	err := s.chainProvider.RelayPacket(context.Background(), &s.relayingPacket)
	s.Require().NoError(err)
}

func (s *LegacyProviderTestSuite) TestRelayPacketSuccessWithoutQueryMaxGasFee() {
	s.chainProvider.Config.MaxGasPrice = 2_000_000_000

	// mock client responses
	s.client.EXPECT().EstimateGas(gomock.Any(), ethereum.CallMsg{
		From:     s.mockSignerAddress,
		To:       &s.chainProvider.TunnelRouterAddress,
		Data:     s.relayingCalldata,
		GasPrice: big.NewInt(2_000_000_000),
	}).Return(uint64(200_000), nil)

	txHash := "0xabc123"
	s.client.EXPECT().BroadcastTx(gomock.Any(), gomock.Any()).Return(txHash, nil)
	s.client.EXPECT().GetTxReceipt(gomock.Any(), txHash).Return(&gethtypes.Receipt{
		Status:      gethtypes.ReceiptStatusSuccessful,
		GasUsed:     21000,
		BlockNumber: big.NewInt(100),
	}, nil)

	s.client.EXPECT().GetBlockHeight(gomock.Any()).Return(uint64(105), nil)
	s.MockDefaultResponses()

	err := s.chainProvider.RelayPacket(context.Background(), &s.relayingPacket)
	s.Require().NoError(err)
}

func (s *LegacyProviderTestSuite) TestRelayPacketFailedGasEstimation() {
	// mock client responses
	s.client.EXPECT().EstimateGasPrice(gomock.Any()).Return(nil, fmt.Errorf("failed to estimate gas price"))
	s.MockDefaultResponses()

	err := s.chainProvider.RelayPacket(context.Background(), &s.relayingPacket)
	s.Require().ErrorContains(err, "failed to estimate gas price")
}

func (s *LegacyProviderTestSuite) TestBumpAndBoundGas() {
	s.MockDefaultResponses()

	testCases := []struct {
		name             string
		maxGasPrice      uint64
		initialGasPrice  int64
		multiplier       float64
		expectedGasPrice int64
	}{
		{
			name:             "Gas price within limit",
			maxGasPrice:      15_000_000_000,
			initialGasPrice:  10_000_000_000,
			multiplier:       1.2,
			expectedGasPrice: 12_000_000_000,
		},
		{
			name:             "Gas price exceeding limit",
			maxGasPrice:      15_000_000_000,
			initialGasPrice:  14_000_000_000,
			multiplier:       1.2,
			expectedGasPrice: 15_000_000_000,
		},
		{
			name:             "No gas price cap, use relayer fee",
			maxGasPrice:      0,
			initialGasPrice:  11_000_000_000,
			multiplier:       1.2,
			expectedGasPrice: 12_000_000_000,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.chainProvider.Config.MaxGasPrice = tc.maxGasPrice

			actual, err := s.chainProvider.BumpAndBoundGas(
				context.Background(),
				evm.NewGasLegacyInfo(big.NewInt(tc.initialGasPrice)),
				tc.multiplier,
			)
			s.Require().NoError(err)

			expected := evm.NewGasLegacyInfo(big.NewInt(tc.expectedGasPrice))
			s.Require().Equal(expected, actual, "Failed test case: %s", tc.name)
		})
	}
}

func (s *LegacyProviderTestSuite) TestEstimateGas() {
	s.client.EXPECT().EstimateGasPrice(gomock.Any()).Return(big.NewInt(5_000_000_000), nil)
	s.MockDefaultResponses()

	actual, err := s.chainProvider.EstimateGasFee(context.Background())
	s.Require().NoError(err)

	expected := evm.GasInfo{
		Type:           evm.GasTypeLegacy,
		GasPrice:       big.NewInt(5_000_000_000),
		GasPriorityFee: nil,
		GasBaseFee:     nil,
		GasFeeCap:      nil,
	}

	s.Require().Equal(expected, actual)
}

func (s *LegacyProviderTestSuite) TestNewRelayTx() {
	data := []byte("mock calldata")
	callMsg := ethereum.CallMsg{
		From:     s.mockSignerAddress,
		To:       &s.chainProvider.TunnelRouterAddress,
		Data:     data,
		GasPrice: s.gasInfo.GasPrice,
	}

	s.client.EXPECT().EstimateGas(gomock.Any(), callMsg).Return(uint64(100), nil)
	s.client.EXPECT().NonceAt(gomock.Any(), s.mockSignerAddress).Return(uint64(1), nil)

	actual, err := s.chainProvider.NewRelayTx(context.Background(), data, s.mockSigner, s.gasInfo)
	s.Require().NoError(err)

	expected := gethtypes.NewTx(&gethtypes.LegacyTx{
		Nonce:    uint64(1),
		To:       &s.chainProvider.TunnelRouterAddress,
		Value:    decimal.NewFromInt(0).BigInt(),
		Data:     data,
		Gas:      uint64(100),
		GasPrice: s.gasInfo.GasPrice,
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
