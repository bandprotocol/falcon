package evm_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/evm"
	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/logger"
	"github.com/bandprotocol/falcon/relayer/wallet/geth"
)

const (
	keyName1    = "testkey1"
	address1    = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	privateKey1 = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	keyName2 = "testkey2"
	address2 = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
	url      = "http://127.0.0.1:8545"
)

var evmCfg = &evm.EVMChainProviderConfig{
	BaseChainProviderConfig: chains.BaseChainProviderConfig{
		Endpoints:           []string{"http://localhost:8545"},
		ChainType:           chainstypes.ChainTypeEVM,
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
	GasType:                    evm.GasTypeEIP1559,
	GasMultiplier:              1.1,
}

type SenderTestSuite struct {
	suite.Suite

	homePath  string
	chainName string
}

func TestSenderTestSuite(t *testing.T) {
	suite.Run(t, new(SenderTestSuite))
}

// SetupTest sets up the test suite by creating a temporary directory and declare mock objects.
func (s *SenderTestSuite) SetupTest() {
	var err error
	tmpDir := s.T().TempDir()
	s.homePath = tmpDir

	s.chainName = "testnet"

	log := logger.NewZapLogWrapper(zap.NewNop().Sugar())
	s.Require().NoError(err)

	client := evm.NewClient(s.chainName, evmCfg, log, nil)

	wallet, err := geth.NewGethWallet("", s.homePath, s.chainName)
	s.Require().NoError(err)

	chainProvider, err := evm.NewEVMChainProvider(s.chainName, client, evmCfg, log, wallet, nil)
	s.Require().NoError(err)

	// Add two mock keys to the chain provider
	_, err = chainProvider.AddKeyByPrivateKey(keyName1, privateKey1)
	s.Require().NoError(err)

	testKey := "testKey"
	_, err = chainProvider.AddRemoteSignerKey(keyName2, address2, url, &testKey)
	s.Require().NoError(err)
}

func (s *SenderTestSuite) TestLoadFreeSenders() {
	log := logger.NewZapLogWrapper(zap.NewNop().Sugar())

	client := evm.NewClient(s.chainName, evmCfg, log, nil)

	wallet, err := geth.NewGethWallet("", s.homePath, s.chainName)
	s.Require().NoError(err)

	chainProvider, err := evm.NewEVMChainProvider(s.chainName, client, evmCfg, log, wallet, nil)
	s.Require().NoError(err)

	err = chainProvider.LoadSigners()
	s.Require().NoError(err)

	count := len(chainProvider.Wallet.GetSigners())
	s.Require().
		Equal(count, len(chainProvider.FreeSigners))

	expectedSenders := map[string]string{
		keyName1: address1,
		keyName2: address2,
	}

	// Check all signers in the channel
	for i := 0; i < count; i++ {
		sender := <-chainProvider.FreeSigners
		s.Require().NotNil(sender)

		name := sender.GetName()
		actualAddress := sender.GetAddress()

		expectedAddress, exists := expectedSenders[name]
		s.Require().True(exists, "Unexpected signer name: %s", name)
		s.Require().Equal(expectedAddress, actualAddress)

		// Remove the validated sender from the map
		delete(expectedSenders, name)
	}
}
