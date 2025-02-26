package evm_test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/relayer/chains/evm"
	chaintypes "github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

const (
	testPrivateKey = "0x72d4772a70645a5a5ec3fdc27afda98d2860a6f7903bff5fd45c0a23d7982121"
	testAddress    = "0x990Ec0f6dFc9e8eE20dec3Ab855D03007A9dD946"
	testMnemonic   = "repeat sugar clarify visa chief soon walnut kangaroo rude parrot height piano spoil desk basket swim income catalog more plunge supreme above later worry"
)

type KeysTestSuite struct {
	suite.Suite

	chainProvider *evm.EVMChainProvider
	log           *zap.Logger
	homePath      string
	wallet        wallet.Wallet
}

func TestKeysTestSuite(t *testing.T) {
	suite.Run(t, new(KeysTestSuite))
}

// SetupTest sets up the test suite by creating a temporary directory and declare mock objects.
func (s *KeysTestSuite) SetupTest() {
	s.homePath = s.T().TempDir()
	s.log = zap.NewNop()

	chainName := "testnet"
	client := evm.NewClient(chainName, evmCfg, s.log)

	wallet, err := wallet.NewGethWallet("", s.homePath, chainName)
	s.Require().NoError(err)

	chainProvider, err := evm.NewEVMChainProvider(chainName, client, evmCfg, s.log, wallet)
	s.Require().NoError(err)

	s.chainProvider = chainProvider
	s.wallet = wallet
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
				keyName: "testkey",
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
		{
			name: "duplicate private key",
			input: Input{
				keyName: "testkey3",
				privKey: testPrivateKey,
			},
			err: fmt.Errorf("account already exists"),
		},
	}

	for _, tc := range testcases {
		s.T().Run(tc.name, func(t *testing.T) {
			key, err := s.chainProvider.AddKeyByPrivateKey(tc.input.keyName, tc.input.privKey)

			if tc.err != nil {
				s.Require().ErrorContains(err, tc.err.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.out, key)

				// check that key info actually stored in local disk
				_, ok := s.wallet.GetAddress(tc.input.keyName)
				s.Require().True(ok)
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
			out: chaintypes.NewKey(testMnemonic, testAddress, ""),
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
			out: chaintypes.NewKey(testMnemonic, "0x01AF9badF97c97C9444E0b7fa94b69b8CB3C28e7", ""),
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
			name: "duplicate key name",
			input: Input{
				keyName:  "testkey",
				mnemonic: "",
				coinType: 60,
				account:  0,
				index:    0,
			},
			err: fmt.Errorf("key name exists"),
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

				// check that key info actually stored in local disk
				_, ok := s.wallet.GetAddress(tc.input.keyName)
				s.Require().True(ok)
			}
		})
	}
}

func (s *KeysTestSuite) TestDeleteKey() {
	keyName := "deletablekey"
	privatekeyHex := testPrivateKey

	// Add a key to delete
	_, err := s.chainProvider.AddKeyByPrivateKey(keyName, privatekeyHex)
	s.Require().NoError(err)

	// Delete the key
	err = s.chainProvider.DeleteKey(keyName)
	s.Require().NoError(err)

	// Ensure the key is no longer in the KeyInfo or KeyStore
	s.Require().False(s.chainProvider.IsKeyNameExist(keyName))

	// Delete the key again should return error
	err = s.chainProvider.DeleteKey(keyName)
	s.Require().ErrorContains(err, "key name does not exist")
}

func (s *KeysTestSuite) TestExportPrivateKey() {
	keyName := "exportkey"
	privatekeyHex := testPrivateKey

	// Add a key to export
	_, err := s.chainProvider.AddKeyByPrivateKey(keyName, privatekeyHex)
	s.Require().NoError(err)

	// Export the private key
	exportedKey, err := s.chainProvider.ExportPrivateKey(keyName)
	s.Require().NoError(err)

	s.Require().Equal(evm.StripPrivateKeyPrefix(privatekeyHex), evm.StripPrivateKeyPrefix(exportedKey))
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

	key2, err := s.chainProvider.AddKeyByMnemonic(
		keyName2,
		mnemonic,
		uint32(coinType),
		uint(account),
		uint(index),
	)
	s.Require().NoError(err)

	// List all keys
	actual := s.chainProvider.ListKeys()
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
	keyName := "showkey"
	privatekeyHex := testPrivateKey

	// Add a key to show
	_, err := s.chainProvider.AddKeyByPrivateKey(keyName, privatekeyHex)
	s.Require().NoError(err)

	// Show the key
	address, err := s.chainProvider.ShowKey(keyName)
	s.Require().Equal(address, address)
	s.Require().NoError(err)
}

func (s *KeysTestSuite) TestIsKeyNameExist() {
	priv, err := crypto.HexToECDSA(evm.StripPrivateKeyPrefix(testPrivateKey))
	s.Require().NoError(err)

	_, err = s.chainProvider.Wallet.SavePrivateKey("testkey1", priv)
	s.Require().NoError(err)

	expected := s.chainProvider.IsKeyNameExist("testkey1")

	s.Require().Equal(expected, true)

	expected = s.chainProvider.IsKeyNameExist("testkey2")
	s.Require().Equal(expected, false)
}
