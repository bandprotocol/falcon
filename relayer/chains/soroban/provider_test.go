package soroban_test

import (
	"context"
	"testing"
	"time"

	cmbytes "github.com/cometbft/cometbft/libs/bytes"
	"github.com/shopspring/decimal"
	hProtocol "github.com/stellar/go-stellar-sdk/protocols/horizon"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/internal/bandchain/tss"
	"github.com/bandprotocol/falcon/internal/relayertest/mocks"
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/soroban"
	chaintypes "github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/logger"
	walletsoroban "github.com/bandprotocol/falcon/relayer/wallet/soroban"
)

var baseSorobanCfg = &soroban.SorobanChainProviderConfig{
	BaseChainProviderConfig: chains.BaseChainProviderConfig{
		Endpoints:                  []string{"http://localhost:8545"},
		ChainType:                  chaintypes.ChainTypeSoroban,
		MaxRetry:                   3,
		QueryTimeout:               3 * time.Second,
		ExecuteTimeout:             3 * time.Second,
		LivelinessCheckingInterval: 15 * time.Minute,
	},
	HorizonEndpoints:   []string{"https://horizon-testnet.stellar.org"},
	NetworkPassphrase:  "Test SDF Network ; September 2015",
	WaitingTxDuration:  time.Second * 3,
	CheckingTxInterval: time.Second,
}

func mockPacket() bandtypes.Packet {
	relatedMsg := cmbytes.HexBytes("0xdeadbeef")
	rAddr := cmbytes.HexBytes("0xfad9c8855b740a0b7ed4c221dbad0f33a83a49ca")
	signature := cmbytes.HexBytes("0xabcd")

	evmSignature := bandtypes.NewEVMSignature(rAddr, signature)
	signingInfo := bandtypes.NewSigning(
		1,
		relatedMsg,
		evmSignature,
		tss.SIGNING_STATUS_SUCCESS,
	)

	return bandtypes.Packet{
		TunnelID:      1,
		Sequence:      42,
		TargetAddress: "GBV4CVR37D6TCX3I7FVSVUX2R6EIOH2AODVQLYY3474FUTG3YNDJ6P75",
		SignalPrices: []bandtypes.SignalPrice{
			{SignalID: "signal1", Price: 100},
			{SignalID: "signal2", Price: 200},
		},
		CurrentGroupSigning:  signingInfo,
		IncomingGroupSigning: nil,
	}
}

type ProviderTestSuite struct {
	suite.Suite

	ctrl          *gomock.Controller
	chainProvider *soroban.SorobanChainProvider
	client        *mocks.MockSorobanClient
	log           *zap.Logger
	homePath      string
	chainName     string
}

func TestProviderTestSuite(t *testing.T) {
	suite.Run(t, new(ProviderTestSuite))
}

func (s *ProviderTestSuite) SetupTest() {
	var err error
	tmpDir := s.T().TempDir()
	s.homePath = tmpDir

	s.ctrl = gomock.NewController(s.T())
	s.client = mocks.NewMockSorobanClient(s.ctrl)

	s.log = zap.NewNop()

	s.chainName = "soroban-testnet"

	wallet, err := walletsoroban.NewWallet("", s.homePath, s.chainName)
	s.Require().NoError(err)

	log := logger.NewZapLogWrapper(zap.NewNop().Sugar())
	s.chainProvider = soroban.NewSorobanChainProvider(s.chainName, s.client, baseSorobanCfg, log, wallet, nil)

	s.chainProvider.Client = s.client
}

func (s *ProviderTestSuite) TestQueryTunnelInfo() {
	tunnelID := uint64(1)
	tunnelAddr := "GBV4CVR37D6TCX3I7FVSVUX2R6EIOH2AODVQLYY3474FUTG3YNDJ6P75"

	tunnel, err := s.chainProvider.QueryTunnelInfo(context.Background(), tunnelID, tunnelAddr)
	s.Require().NoError(err)
	s.Require().Equal(tunnelID, tunnel.ID)
	s.Require().Equal(tunnelAddr, tunnel.TargetAddress)
	s.Require().True(tunnel.IsActive)
	s.Require().Nil(tunnel.LatestSequence)
	s.Require().Nil(tunnel.Balance)
}

func (s *ProviderTestSuite) TestCheckConfirmedTx() {
	txHash := "abc123"

	testcases := []struct {
		name       string
		preProcess func()
		err        error
		out        soroban.TxResult
	}{
		{
			name: "success",
			preProcess: func() {
				s.client.EXPECT().GetTransactionStatus(txHash).Return(hProtocol.Transaction{
					Hash:       txHash,
					Ledger:     100,
					FeeCharged: 100,
					Successful: true,
				}, nil)
			},
			out: soroban.NewTxResult(
				chaintypes.TX_STATUS_SUCCESS,
				txHash,
				100,
				decimal.NewNullDecimal(decimal.NewFromInt(100)),
				"",
			),
		},
		{
			name: "transaction failed",
			preProcess: func() {
				s.client.EXPECT().GetTransactionStatus(txHash).Return(hProtocol.Transaction{
					Hash:       txHash,
					Ledger:     100,
					FeeCharged: 100,
					Successful: false,
				}, nil)
			},
			out: soroban.NewTxResult(
				chaintypes.TX_STATUS_FAILED,
				txHash,
				100,
				decimal.NewNullDecimal(decimal.NewFromInt(100)),
				"transaction failed on-chain",
			),
		},
		{
			name: "client error",
			preProcess: func() {
				s.client.EXPECT().
					GetTransactionStatus(txHash).
					Return(hProtocol.Transaction{}, context.DeadlineExceeded)
			},
			err: context.DeadlineExceeded,
			out: soroban.NewTxResult(
				chaintypes.TX_STATUS_PENDING,
				txHash,
				0,
				decimal.NullDecimal{},
				context.DeadlineExceeded.Error(),
			),
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			if tc.preProcess != nil {
				tc.preProcess()
			}

			res, err := s.chainProvider.CheckConfirmedTx(txHash)
			if tc.err != nil {
				s.Require().ErrorIs(err, tc.err)
			} else {
				s.Require().NoError(err)
			}
			s.Require().Equal(tc.out, res)
		})
	}
}
