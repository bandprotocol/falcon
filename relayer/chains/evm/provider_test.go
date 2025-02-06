package evm_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	cmbytes "github.com/cometbft/cometbft/libs/bytes"
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
	"github.com/bandprotocol/falcon/relayer/chains/evm"
	chaintypes "github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

var baseEVMCfg = &evm.EVMChainProviderConfig{
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

func mockPacket() bandtypes.Packet {
	relatedMsg := cmbytes.HexBytes("0xdeadbeef")
	rAddr := gethcommon.HexToAddress("0xfad9c8855b740a0b7ed4c221dbad0f33a83a49ca")
	signature := cmbytes.HexBytes("0xabcd")

	evmSignature := bandtypes.NewEVMSignature(rAddr.Bytes(), signature)
	signingInfo := bandtypes.NewSigning(
		1,
		relatedMsg,
		evmSignature,
		"SIGNING_STATUS_SUCCESS",
	)

	return bandtypes.Packet{
		TunnelID: 1,
		Sequence: 42,
		SignalPrices: []bandtypes.SignalPrice{
			{SignalID: "signal1", Price: 100},
			{SignalID: "signal2", Price: 200},
		},
		CurrentGroupSigning:  signingInfo,
		IncomingGroupSigning: nil,
	}
}

func mockSender() (evm.Sender, error) {
	addr, err := evm.HexToAddress(testAddress)
	if err != nil {
		return evm.Sender{}, err
	}

	priv, err := crypto.HexToECDSA(evm.StripPrivateKeyPrefix(testPrivateKey))
	if err != nil {
		return evm.Sender{}, err
	}

	return evm.Sender{
		Address:    addr,
		PrivateKey: priv,
	}, nil
}

func uint256ToHex(value *big.Int) string {
	return fmt.Sprintf("%064x", value)
}

type ProviderTestSuite struct {
	suite.Suite

	ctrl          *gomock.Controller
	chainProvider *evm.EVMChainProvider
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
	s.homePath = tmpDir

	s.ctrl = gomock.NewController(s.T())
	s.client = mocks.NewMockEVMClient(s.ctrl)

	// mock objects.
	s.log = zap.NewNop()

	chainName := "testnet"
	s.chainName = chainName

	wallet, err := wallet.NewGethWallet("", s.homePath, s.chainName)
	s.Require().NoError(err)

	s.chainProvider, err = evm.NewEVMChainProvider(s.chainName, s.client, baseEVMCfg, s.log, s.homePath, wallet)
	s.Require().NoError(err)

	s.chainProvider.Client = s.client
}

func (s *ProviderTestSuite) TestQueryTunnelInfo() {
	queryTunnelCalldata, err := hex.DecodeString(
		"077071ef0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000e688b84b23f322a994a53dbf8e15fa82cdb71127",
	)
	s.Require().NoError(err)

	// abi-encoded from {"isActive": True,"latestSequence": 1,"balance": 1000000000000000000}
	queryTunnelResponse, err := hex.DecodeString(
		"000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000de0b6b3a7640000",
	)
	s.Require().NoError(err)

	type Input struct {
		tunnelID   uint64
		tunnelAddr string
	}
	testcases := []struct {
		name       string
		input      Input
		preProcess func()
		err        error
		out        *chaintypes.Tunnel
	}{
		{
			name:  "success",
			input: Input{1, "0xe688b84b23f322a994A53dbF8E15FA82CDB71127"},
			preProcess: func() {
				s.client.EXPECT().CheckAndConnect(gomock.Any()).Return(nil)
				s.client.EXPECT().
					Query(gomock.Any(), s.chainProvider.TunnelRouterAddress, queryTunnelCalldata).
					Return(queryTunnelResponse, nil)
			},
			out: &chaintypes.Tunnel{
				ID:             1,
				TargetAddress:  "0xe688b84b23f322a994A53dbF8E15FA82CDB71127",
				IsActive:       true,
				LatestSequence: 1,
				Balance:        big.NewInt(1000000000000000000),
			},
		},
		{
			name:  "failed to connect client",
			input: Input{1, "0xe688b84b23f322a994A53dbF8E15FA82CDB71127"},
			preProcess: func() {
				s.client.EXPECT().CheckAndConnect(gomock.Any()).Return(fmt.Errorf("Connect client error"))
			},
			err: fmt.Errorf("Connect client error"),
		},
		{
			name:  "invalid target address",
			input: Input{1, "0xincorrect"},
			preProcess: func() {
				s.client.EXPECT().CheckAndConnect(gomock.Any()).Return(nil)
			},
			err: fmt.Errorf("invalid address"),
		},
		{
			name:  "cannot unpack data",
			input: Input{1, "0xe688b84b23f322a994A53dbF8E15FA82CDB71127"},
			preProcess: func() {
				s.client.EXPECT().CheckAndConnect(gomock.Any()).Return(nil)
				s.client.EXPECT().
					Query(gomock.Any(), s.chainProvider.TunnelRouterAddress, queryTunnelCalldata).
					Return([]uint8{0, 124}, nil)
			},
			err: fmt.Errorf("failed to unpack data"),
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			if tc.preProcess != nil {
				tc.preProcess()
			}
			defer s.ctrl.Finish()

			tunnel, err := s.chainProvider.QueryTunnelInfo(
				context.Background(),
				tc.input.tunnelID,
				tc.input.tunnelAddr,
			)

			if tc.err != nil {
				s.Require().ErrorContains(err, tc.err.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.out, tunnel)
			}
		})
	}
}

func (s *ProviderTestSuite) TestEstimateGasUnsupportedGas() {
	_, err := s.chainProvider.EstimateGasFee(context.Background())
	s.Require().ErrorContains(err, "unsupported gas type:")
}

func (s *ProviderTestSuite) TestCheckConfirmedTx() {
	txHash := "0xabc123"
	txBlock := int64(100)

	testcases := []struct {
		name       string
		preProcess func()
		err        error
		out        *evm.ConfirmTxResult
	}{
		{
			name: "success",
			preProcess: func() {
				currentBlock := txBlock + int64(s.chainProvider.Config.BlockConfirmation) + 10

				s.client.EXPECT().GetTxReceipt(gomock.Any(), txHash).Return(&gethtypes.Receipt{
					Status:      gethtypes.ReceiptStatusSuccessful,
					GasUsed:     21000,
					BlockNumber: big.NewInt(txBlock),
				}, nil)
				s.client.EXPECT().GetBlockHeight(gomock.Any()).Return(uint64(currentBlock), nil)
			},
			out: evm.NewConfirmTxResult(
				txHash,
				evm.TX_STATUS_SUCCESS,
				decimal.NewNullDecimal(decimal.New(21000, 0)),
			),
		},
		{
			name: "get tx receipt with failed status",
			preProcess: func() {
				s.client.EXPECT().GetTxReceipt(gomock.Any(), txHash).Return(&gethtypes.Receipt{
					Status:      gethtypes.ReceiptStatusFailed,
					GasUsed:     21000,
					BlockNumber: big.NewInt(txBlock),
				}, nil)
			},
			out: evm.NewConfirmTxResult(
				txHash,
				evm.TX_STATUS_FAILED,
				decimal.NullDecimal{},
			),
		},
		{
			name: "get tx receipt but not confirmed block",
			preProcess: func() {
				currentBlock := txBlock + int64(s.chainProvider.Config.BlockConfirmation) - 1

				s.client.EXPECT().GetTxReceipt(gomock.Any(), txHash).Return(&gethtypes.Receipt{
					Status:      gethtypes.ReceiptStatusSuccessful,
					GasUsed:     21000,
					BlockNumber: big.NewInt(txBlock),
				}, nil)
				s.client.EXPECT().GetBlockHeight(gomock.Any()).Return(uint64(currentBlock), nil)
			},
			out: evm.NewConfirmTxResult(
				txHash,
				evm.TX_STATUS_UNMINED,
				decimal.NullDecimal{},
			),
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			if tc.preProcess != nil {
				tc.preProcess()
			}

			expect, err := s.chainProvider.CheckConfirmedTx(context.Background(), txHash)
			if tc.err != nil {
				s.Require().ErrorContains(err, tc.err.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.out, expect)
			}
		})
	}
}
