package evm_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/evm"
	chaintypes "github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/logger"
	"github.com/bandprotocol/falcon/relayer/wallet"
	"github.com/bandprotocol/falcon/relayer/wallet/geth"
)

const (
	testKey        = "testKey"
	testPrivateKey = "0x72d4772a70645a5a5ec3fdc27afda98d2860a6f7903bff5fd45c0a23d7982121"
	testAddress    = "0x990Ec0f6dFc9e8eE20dec3Ab855D03007A9dD946"
	testMnemonic   = "repeat sugar clarify visa chief soon walnut kangaroo rude parrot height piano spoil desk basket swim income catalog more plunge supreme above later worry"
)

var evmCfg = &evm.EVMChainProviderConfig{
	BaseChainProviderConfig: chains.BaseChainProviderConfig{
		Endpoints: []string{"http://localhost:8545"},
		ChainType: chaintypes.ChainTypeEVM,
		MaxRetry:  3,
		ChainID:   31337,

		QueryTimeout:   3 * time.Second,
		ExecuteTimeout: 3 * time.Second,
	},
	TunnelRouterAddress:        "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9",
	BlockConfirmation:          5,
	WaitingTxDuration:          time.Second * 3,
	CheckingTxInterval:         time.Second,
	LivelinessCheckingInterval: 15 * time.Minute,
	GasType:                    evm.GasTypeEIP1559,
	GasMultiplier:              1.1,
}

func TestKeysTestSuite(t *testing.T) {
	suite.Run(t, new(KeysTestSuite))
}

type KeysTestSuite struct {
	suite.Suite

	chainProvider *evm.EVMChainProvider
	log           logger.Logger
	homePath      string
	wallet        wallet.Wallet
}

func (s *KeysTestSuite) loadChainProvider() {
	s.log = logger.NewZapLogWrapper(zap.NewNop().Sugar())

	chainName := "testnet"
	client := evm.NewClient(chainName, evmCfg, s.log, nil)

	wallet, err := geth.NewGethWallet("", s.homePath, chainName)
	s.Require().NoError(err)

	chainProvider, err := evm.NewEVMChainProvider(chainName, client, evmCfg, s.log, wallet, nil)
	s.Require().NoError(err)

	s.chainProvider = chainProvider
	s.wallet = wallet
}

// SetupTest sets up the test suite by creating a temporary directory and declare mock objects.
func (s *KeysTestSuite) SetupTest() {
	s.homePath = s.T().TempDir()
	s.loadChainProvider()
}

func (s *KeysTestSuite) TestAddKeyByPrivateKey() {
	type Input struct {
		keyName string
		privKey string
	}
	testcases := []struct {
		name  string
		input Input
		err   error
		out   *chaintypes.Key
	}{
		{
			name: "success",
			input: Input{
				keyName: "testkey2",
				privKey: testPrivateKey,
			},
			out: chaintypes.NewKey("", testAddress, ""),
		},
		{
			name: "invalid private key",
			input: Input{
				keyName: "testkey2",
				privKey: "x72d4772a70645a5a5ec3fdc27afda98d2860a6f7903bff5fd45c0a23d7982121",
			},
			err: fmt.Errorf("invalid hex character"),
		},
	}

	for _, tc := range testcases {
		s.T().Run(tc.name, func(t *testing.T) {
			key, err := chains.AddKeyByPrivateKey(s.wallet, tc.input.keyName, tc.input.privKey)

			if tc.err != nil {
				s.Require().ErrorContains(err, tc.err.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.out, key)
			}
		})
	}
}

func (s *KeysTestSuite) TestAddKeyByMnemonic() {
	type Input struct {
		keyName  string
		mnemonic string
		coinType uint32
		account  uint
		index    uint
	}
	testcases := []struct {
		name  string
		input Input
		err   error
		out   *chaintypes.Key
	}{
		{
			name: "success",
			input: Input{
				keyName:  "testkey",
				mnemonic: testMnemonic,
				coinType: 60,
				account:  0,
				index:    0,
			},
			out: chaintypes.NewKey("", testAddress, ""),
		},
		{
			name: "success with different index",
			input: Input{
				keyName:  "testkey2",
				mnemonic: testMnemonic,
				coinType: 60,
				account:  0,
				index:    1,
			},
			out: chaintypes.NewKey("", "0x01AF9badF97c97C9444E0b7fa94b69b8CB3C28e7", ""),
		},
		{
			name: "success with no mnemonic",
			input: Input{
				keyName:  "testkey3",
				mnemonic: "",
				coinType: 60,
				account:  0,
				index:    0,
			},
		},
		{
			name: "invalid mnemonic",
			input: Input{
				keyName:  "testkey4",
				mnemonic: "mnemonic",
				coinType: 60,
				account:  0,
				index:    0,
			},
			err: fmt.Errorf("mnemonic is invalid"),
		},
	}

	for _, tc := range testcases {
		s.T().Run(tc.name, func(t *testing.T) {
			key, err := s.chainProvider.AddKeyByMnemonic(
				tc.input.keyName,
				tc.input.mnemonic,
				tc.input.coinType,
				tc.input.account,
				tc.input.index,
			)

			if tc.err != nil {
				s.Require().ErrorContains(err, tc.err.Error())
			} else {
				s.Require().NoError(err)

				if tc.out != nil {
					s.Require().Equal(tc.out, key)
				}

				if tc.input.mnemonic == "" {
					s.Require().NotEmpty(
						key.Mnemonic,
						"expected generated mnemonic to be returned when none is provided",
					)
				}
			}
		})
	}
}

func (s *KeysTestSuite) TestAddRemoteSignerKey() {
	testKey := "testKey"
	type Input struct {
		keyName string
		addr    string
		url     string
		key     *string
	}
	testcases := []struct {
		name  string
		input Input
		err   error
		out   *chaintypes.Key
	}{
		{
			name: "success",
			input: Input{
				keyName: "remotekey",
				addr:    testAddress,
				url:     "http://127.0.0.1:8545",
				key:     &testKey,
			},
			out: chaintypes.NewKey("", testAddress, ""),
		},
		{
			name: "nil key",
			input: Input{
				keyName: "nilkey",
				addr:    testAddress,
				url:     "http://127.0.0.1:8545",
				key:     nil,
			},
			out: chaintypes.NewKey("", testAddress, ""),
		},
	}

	for _, tc := range testcases {
		s.T().Run(tc.name, func(t *testing.T) {
			key, err := chains.AddRemoteSignerKey(
				s.wallet,
				tc.input.keyName,
				tc.input.addr,
				tc.input.url,
				tc.input.key,
			)

			s.Require().NoError(err)
			s.Require().Equal(tc.out, key)
		})
	}
}

func (s *KeysTestSuite) TestDeleteKey() {
	// Add key to delete
	_, err := chains.AddKeyByPrivateKey(s.wallet, testKey, testPrivateKey)
	s.Require().NoError(err)

	s.loadChainProvider()

	// Delete the key
	err = chains.DeleteKey(s.wallet, testKey)
	s.Require().NoError(err)
}

func (s *KeysTestSuite) TestExportPrivateKey() {
	tests := []struct {
		name      string
		keyName   string
		setup     func()
		wantErr   bool
		errSubstr string
	}{
		{
			name:    "success",
			keyName: testKey,
			setup: func() {
				_, err := chains.AddKeyByPrivateKey(s.wallet, testKey, testPrivateKey)
				s.Require().NoError(err)
			},
		},
		{
			name:      "key name does not exist",
			keyName:   "doesNotExist",
			wantErr:   true,
			errSubstr: "key name does not exist",
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			s.loadChainProvider()
			exported, err := chains.ExportPrivateKey(s.wallet, tc.keyName)
			if tc.wantErr {
				s.Require().ErrorContains(err, tc.errSubstr)
				return
			}
			s.Require().NoError(err)
			s.Require().Equal(
				evm.StripPrivateKeyPrefix(testPrivateKey),
				evm.StripPrivateKeyPrefix(exported),
			)
		})
	}
}

func (s *KeysTestSuite) TestListKeys() {
	// Add multiple keys
	keyName1 := "key1"
	keyName2 := "key2"
	mnemonic := ""
	coinType := 60
	account := 0
	index := 0

	key1, err := s.chainProvider.AddKeyByMnemonic(
		keyName1,
		mnemonic,
		uint32(coinType),
		uint(account),
		uint(index),
	)
	s.Require().NoError(err)

	s.loadChainProvider()

	key2, err := s.chainProvider.AddKeyByMnemonic(
		keyName2,
		mnemonic,
		uint32(coinType),
		uint(account),
		uint(index),
	)
	s.Require().NoError(err)

	s.loadChainProvider()

	// List all keys
	actual := chains.ListKeys(s.wallet)
	s.Require().Equal(2, len(actual))

	expected1 := chaintypes.NewKey("", key1.Address, keyName1)
	expected2 := chaintypes.NewKey("", key2.Address, keyName2)

	// Check if expected1 and expected2 are in actual
	foundExpected1 := false
	foundExpected2 := false

	for _, key := range actual {
		if key.Address == expected1.Address {
			foundExpected1 = true
		}
		if key.Address == expected2.Address {
			foundExpected2 = true
		}
	}

	s.Require().True(foundExpected1)
	s.Require().True(foundExpected2)
}

func (s *KeysTestSuite) TestShowKey() {
	tests := []struct {
		name      string
		keyName   string
		setup     func()
		wantErr   bool
		errSubstr string
	}{
		{
			name:    "success",
			keyName: testKey,
			setup: func() {
				_, err := chains.AddKeyByPrivateKey(s.wallet, testKey, testPrivateKey)
				s.Require().NoError(err)
			},
		},
		{
			name:      "key name not exist",
			keyName:   "doesNotExist",
			wantErr:   true,
			errSubstr: "key name does not exist",
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			s.loadChainProvider()
			address, err := chains.ShowKey(s.wallet, tc.keyName)
			if tc.wantErr {
				s.Require().ErrorContains(err, tc.errSubstr)
				return
			}
			s.Require().NoError(err)
			s.Require().Equal(
				testAddress,
				address,
			)
		})
	}
}
