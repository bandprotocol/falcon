package evm_test

import (
	"context"
	"encoding/hex"
	"os"
	"path"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/evm"
)

const (
	privateKey1 = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	address1    = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	privateKey2 = "0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"
	address2    = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
)

var evmCfg = &evm.EVMChainProviderConfig{
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
	GasType:                    evm.GasTypeEIP1559,
	GasMultiplier:              1.1,
}

type SenderTestSuite struct {
	suite.Suite

	ctx           context.Context
	chainProvider *evm.EVMChainProvider
	log           *zap.Logger
	homePath      string
}

func TestSenderTestSuite(t *testing.T) {
	suite.Run(t, new(SenderTestSuite))
}

// SetupTest sets up the test suite by creating a temporary directory and declare mock objects.
func (s *SenderTestSuite) SetupTest() {
	var err error
	tmpDir := s.T().TempDir()

	log, err := zap.NewDevelopment()
	s.Require().NoError(err)

	// mock objects.
	s.log = zap.NewNop()

	chainName := "testnet"

	client := evm.NewClient(chainName, evmCfg, log)

	s.chainProvider, err = evm.NewEVMChainProvider(chainName, client, evmCfg, log, tmpDir)
	s.Require().NoError(err)

	s.ctx = context.Background()
	s.homePath = tmpDir
}

func TestLoadKeyInfo(t *testing.T) {
	tmpDir := t.TempDir()
	chainName := "testnet"

	// write mock keyInfo at keyInfo's path
	keyInfo := make(evm.KeyInfo)
	keyInfo["key1"] = ""
	keyInfo["key2"] = ""
	b, err := toml.Marshal(&keyInfo)
	require.NoError(t, err)

	keyInfoDir := path.Join(tmpDir, "keys", chainName, "info")
	keyInfoPath := path.Join(keyInfoDir, "info.toml")
	// Create the info folder if doesn't exist
	err = os.MkdirAll(keyInfoDir, os.ModePerm)
	require.NoError(t, err)
	// Create the file and write the default config to the given location.
	f, err := os.Create(keyInfoPath)
	require.NoError(t, err)
	defer f.Close()

	_, err = f.Write(b)
	require.NoError(t, err)

	// load keyInfo
	actual, err := evm.LoadKeyInfo(tmpDir, chainName)
	require.NoError(t, err)

	require.Equal(t, keyInfo, actual)
}

func (s *SenderTestSuite) TestLoadFreeSenders() {
	keyName1 := "key1"
	keyName2 := "key2"

	// Add two mock keys to the chain provider
	_, err := s.chainProvider.AddKeyWithPrivateKey(keyName1, privateKey1, s.homePath, "")
	s.Require().NoError(err)

	_, err = s.chainProvider.AddKeyWithPrivateKey(keyName2, privateKey2, s.homePath, "")
	s.Require().NoError(err)

	// Load free senders
	err = s.chainProvider.LoadFreeSenders(s.homePath, "")
	s.Require().NoError(err)

	// Validate the FreeSenders channel is populated correctly
	count := len(s.chainProvider.KeyInfo)
	s.Require().
		Equal(count, len(s.chainProvider.FreeSenders))

	// Create a map to check properties of retrieved senders
	expectedSenders := map[string]string{
		address1: privateKey1,
		address2: privateKey2,
	}

	// Check all senders in the channel
	for i := 0; i < count; i++ {
		sender := <-s.chainProvider.FreeSenders
		s.Require().NotNil(sender)

		actualAddress := sender.Address.Hex()
		actualPrivateKey := evm.StripPrivateKeyPrefix(
			hex.EncodeToString(crypto.FromECDSA(sender.PrivateKey)),
		)

		expectedPrivateKey, exists := expectedSenders[actualAddress]
		s.Require().True(exists, "Unexpected sender address: %s", actualAddress)

		// Validate the private key matches
		s.Require().
			Equal(evm.StripPrivateKeyPrefix(expectedPrivateKey), evm.StripPrivateKeyPrefix(actualPrivateKey))

		// Remove the validated sender from the map
		delete(expectedSenders, actualAddress)
	}
}
