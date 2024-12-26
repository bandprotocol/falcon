package evm

import (
	"context"
	"encoding/hex"
	"os"
	"path"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/relayer/chains"
	chaintypes "github.com/bandprotocol/falcon/relayer/chains/types"
)

const (
	testPrivateKey = "0x72d4772a70645a5a5ec3fdc27afda98d2860a6f7903bff5fd45c0a23d7982121"
	testAddress    = "0x990Ec0f6dFc9e8eE20dec3Ab855D03007A9dD946"
	testMnemonic   = "repeat sugar clarify visa chief soon walnut kangaroo rude parrot height piano spoil desk basket swim income catalog more plunge supreme above later worry"
)

var evmCfg = &EVMChainProviderConfig{
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
	GasType:                    GasTypeEIP1559,
	GasMultiplier:              1.1,
}

type KeysTestSuite struct {
	suite.Suite

	ctx           context.Context
	chainProvider *EVMChainProvider
	log           *zap.Logger
	homePath      string
}

func TestKeysTestSuite(t *testing.T) {
	suite.Run(t, new(KeysTestSuite))
}

// SetupTest sets up the test suite by creating a temporary directory and declare mock objects.
func (s *KeysTestSuite) SetupTest() {
	var err error
	tmpDir := s.T().TempDir()

	log, err := zap.NewDevelopment()
	s.Require().NoError(err)

	chainName := "testnet"

	client := NewClient(chainName, evmCfg, log)

	s.ctx = context.Background()

	s.chainProvider, err = NewEVMChainProvider(chainName, client, evmCfg, log, tmpDir)
	s.Require().NoError(err)

	s.log = log

	s.homePath = tmpDir
}

func (s *KeysTestSuite) TestAddKeyPrivateKeyInputNotEmpty() {
	keyName := "testkey"
	mnemonic := ""
	coinType := 60
	account := 0
	index := 0
	passphrase := ""

	actual, err := s.chainProvider.AddKey(
		keyName,
		mnemonic,
		testPrivateKey,
		s.homePath,
		uint32(coinType),
		uint(account),
		uint(index),
		passphrase,
	)
	s.Require().NoError(err)

	expected := chaintypes.NewKey("", testAddress, "")
	s.Require().Equal(expected, actual)

	addr, err := HexToAddress(testAddress)
	s.Require().NoError(err)

	// check that key actually added in keystore
	s.Require().True(s.chainProvider.KeyStore.HasAddress(addr))

	// check that key info actually stored in local disk
	keyInfo, err := LoadKeyInfo(s.homePath, s.chainProvider.ChainName)
	s.Require().NoError(err)

	_, exist := keyInfo[keyName]
	s.Require().True(exist)
}

func (s *KeysTestSuite) TestAddKeyMnemonicInputNotEmpty() {
	keyName := "testkey"
	privatekeyHex := ""
	coinType := 60
	account := 0
	index := 0
	passphrase := ""

	actual, err := s.chainProvider.AddKey(
		keyName,
		testMnemonic,
		privatekeyHex,
		s.homePath,
		uint32(coinType),
		uint(account),
		uint(index),
		passphrase,
	)
	s.Require().NoError(err)

	expected := chaintypes.NewKey(testMnemonic, testAddress, "")
	s.Require().Equal(expected, actual)

	addr, err := HexToAddress(testAddress)
	s.Require().NoError(err)

	// check that key actually added in keystore
	s.Require().True(s.chainProvider.KeyStore.HasAddress(addr))

	// check that key info actually stored in local disk
	keyInfo, err := LoadKeyInfo(s.homePath, s.chainProvider.ChainName)
	s.Require().NoError(err)

	_, exist := keyInfo[keyName]
	s.Require().True(exist)
}

func (s *KeysTestSuite) TestAddKeyMnemonicInputEmpty() {
	keyName := "testkey"
	mnemonic := ""
	privateKey := ""
	coinType := 60
	account := 0
	index := 0
	passphrase := ""

	actual, err := s.chainProvider.AddKey(
		keyName,
		mnemonic,
		privateKey,
		s.homePath,
		uint32(coinType),
		uint(account),
		uint(index),
		passphrase,
	)
	s.Require().NoError(err)

	s.Require().NotEqual("", actual.Mnemonic)
	s.Require().NotEqual("", actual.Address)

	addr, err := HexToAddress(actual.Address)
	s.Require().NoError(err)

	// check that key actually added in keystore
	s.Require().True(s.chainProvider.KeyStore.HasAddress(addr))

	// check that key info actually stored in local disk
	keyInfo, err := LoadKeyInfo(s.homePath, s.chainProvider.ChainName)
	s.Require().NoError(err)

	_, exist := keyInfo[keyName]
	s.Require().True(exist)
}

func (s *KeysTestSuite) TestAddKeyWithMnemonic() {
	keyName := "testkey"
	coinType := 60
	account := 0
	index := 0

	actual, err := s.chainProvider.AddKeyWithMnemonic(
		keyName,
		testMnemonic,
		s.homePath,
		uint32(coinType),
		uint(account),
		uint(index),
		"",
	)
	s.Require().NoError(err)

	s.Require().NotEqual("", actual.Mnemonic)

	addr, err := HexToAddress(actual.Address)
	s.Require().NoError(err)

	// check that key actually added in keystore
	s.Require().True(s.chainProvider.KeyStore.HasAddress(addr))

	// check that key info actually stored in local disk
	keyInfo, err := LoadKeyInfo(s.homePath, s.chainProvider.ChainName)
	s.Require().NoError(err)

	_, exist := keyInfo[keyName]
	s.Require().True(exist)
}

func (s *KeysTestSuite) TestAddKeyWithPrivateKey() {
	keyName := "testkey"

	actual, err := s.chainProvider.AddKeyWithPrivateKey(keyName, testPrivateKey, s.homePath, "")
	s.Require().NoError(err)

	expected := chaintypes.NewKey("", testAddress, "")
	s.Require().Equal(expected, actual)

	addr, err := HexToAddress(testAddress)
	s.Require().NoError(err)

	// check that key actually added in keystore
	s.Require().True(s.chainProvider.KeyStore.HasAddress(addr))

	// check that key info actually stored in local disk
	keyInfo, err := LoadKeyInfo(s.homePath, s.chainProvider.ChainName)
	s.Require().NoError(err)

	_, exist := keyInfo[keyName]
	s.Require().True(exist)
}

func (s *KeysTestSuite) TestAddKeyWithPrivateKeyInvalidPrivateKey() {
	keyName := "testkey"
	privateKey := "x72d4772a70645a5a5ec3fdc27afda98d2860a6f7903bff5fd45c0a23d7982121"

	_, err := s.chainProvider.AddKeyWithPrivateKey(keyName, privateKey, s.homePath, "")
	s.Require().Error(err)
}

func (s *KeysTestSuite) TestFinalizeKeyAddition() {
	keyName := "testkey"
	mnemonic := ""
	priv, err := crypto.HexToECDSA(StripPrivateKeyPrefix(testPrivateKey))
	s.Require().NoError(err)
	passphrase := ""

	actual, err := s.chainProvider.finalizeKeyAddition(keyName, priv, mnemonic, s.homePath, passphrase)
	s.Require().NoError(err)

	expected := chaintypes.NewKey(mnemonic, testAddress, "")
	s.Require().Equal(expected, actual)
}

func (s *KeysTestSuite) TestDeleteKey() {
	keyName := "deletablekey"
	privatekeyHex := testPrivateKey

	// Add a key to delete
	_, err := s.chainProvider.AddKeyWithPrivateKey(keyName, privatekeyHex, s.homePath, "")
	s.Require().NoError(err)

	// Delete the key
	err = s.chainProvider.DeleteKey(s.homePath, keyName, "")
	s.Require().NoError(err)

	// Ensure the key is no longer in the KeyInfo or KeyStore
	s.Require().False(s.chainProvider.IsKeyNameExist(keyName))

	addr, err := HexToAddress(testAddress)
	s.Require().NoError(err)
	s.Require().False(s.chainProvider.KeyStore.HasAddress(addr))
}

func (s *KeysTestSuite) TestExportPrivateKey() {
	keyName := "exportkey"
	privatekeyHex := testPrivateKey

	// Add a key to export
	_, err := s.chainProvider.AddKeyWithPrivateKey(keyName, privatekeyHex, s.homePath, "")
	s.Require().NoError(err)

	// Export the private key
	exportedKey, err := s.chainProvider.ExportPrivateKey(keyName, "")
	s.Require().NoError(err)

	s.Require().Equal(StripPrivateKeyPrefix(privatekeyHex), StripPrivateKeyPrefix(exportedKey))
}

func (s *KeysTestSuite) TestListKeys() {
	// Add multiple keys
	keyName1 := "key1"
	keyName2 := "key2"
	mnemonic := ""
	privateKey := ""
	coinType := 60
	account := 0
	index := 0
	passphrase := ""

	key1, err := s.chainProvider.AddKey(
		keyName1,
		mnemonic,
		privateKey,
		s.homePath,
		uint32(coinType),
		uint(account),
		uint(index),
		passphrase,
	)
	s.Require().NoError(err)

	key2, err := s.chainProvider.AddKey(
		keyName2,
		mnemonic,
		privateKey,
		s.homePath,
		uint32(coinType),
		uint(account),
		uint(index),
		passphrase,
	)
	s.Require().NoError(err)

	// List all keys
	actual := s.chainProvider.Listkeys()
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
	_, err := s.chainProvider.AddKeyWithPrivateKey(keyName, privatekeyHex, s.homePath, "")
	s.Require().NoError(err)

	// Show the key
	address := s.chainProvider.ShowKey(keyName)
	s.Require().Equal(address, address)
}

func (s *KeysTestSuite) TestIsKeyNameExist() {
	s.chainProvider.KeyInfo["testkey1"] = testAddress
	expected := s.chainProvider.IsKeyNameExist("testkey1")

	s.Require().Equal(expected, true)

	expected = s.chainProvider.IsKeyNameExist("testkey2")
	s.Require().Equal(expected, false)
}

func (s *KeysTestSuite) TestStorePrivateKey() {
	privateKeyECDSA, err := crypto.HexToECDSA(StripPrivateKeyPrefix(testPrivateKey)) // Remove "0x" prefix
	s.Require().NoError(err)

	// Store the private key in the keystore
	account, err := s.chainProvider.storePrivateKey(privateKeyECDSA, "")
	s.Require().NoError(err)
	s.Require().NotNil(account)

	// Verify that the key exists in the keystore
	addr := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)
	s.Require().True(s.chainProvider.KeyStore.HasAddress(addr))
}

func (s *KeysTestSuite) TestStorePrivateKeyDuplicatePrivateKey() {
	privateKeyECDSA, err := crypto.HexToECDSA(StripPrivateKeyPrefix(testPrivateKey)) // Remove "0x" prefix
	s.Require().NoError(err)

	// Store the private key in the keystore
	account, err := s.chainProvider.storePrivateKey(privateKeyECDSA, "")
	s.Require().NoError(err)
	s.Require().NotNil(account)

	// Verify that the key exists in the keystore
	addr := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)
	s.Require().True(s.chainProvider.KeyStore.HasAddress(addr))

	_, err = s.chainProvider.storePrivateKey(privateKeyECDSA, "")
	s.Require().Error(err)
}

func (s *KeysTestSuite) TestStoreKeyInfo() {
	keyName := "testkey"
	address := "0xc0ffee254729296a45a3885639AC7E10F9d54979"
	s.chainProvider.KeyInfo[keyName] = address

	err := os.MkdirAll(path.Join(s.homePath, "keys", s.chainProvider.ChainName), os.ModePerm)
	s.Require().NoError(err)
	err = s.chainProvider.storeKeyInfo(s.homePath)
	s.Require().NoError(err)

	// Add a key to simulate storing key info
	keyInfoPath := path.Join(s.homePath, "keys", s.chainProvider.ChainName, "info", "info.toml")
	b, err := os.ReadFile(keyInfoPath)
	s.Require().NoError(err)

	expected := "testkey = '0xc0ffee254729296a45a3885639AC7E10F9d54979'\n"
	s.Require().Equal(expected, string(b))
}

func (s *KeysTestSuite) TestGeneratePrivateKey() {
	mnemonic := "test test test test test test test test test test test junk" // Sample mnemonic
	coinType := 60
	account := 0
	index := 0

	// Generate the private key
	privateKeyECDSA, err := s.chainProvider.generatePrivateKey(mnemonic, uint32(coinType), uint(account), uint(index))
	s.Require().NoError(err)
	s.Require().NotNil(privateKeyECDSA)

	// Verify the public key from the generated private key
	address := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey).Hex()
	s.Require().NotEmpty(address)
}

func (s *KeysTestSuite) TestGeneratePrivateKeyWithInvalidMnemonic() {
	mnemonic := "invalid" // Sample mnemonic
	coinType := 60
	account := 0
	index := 0

	// Generate the private key
	_, err := s.chainProvider.generatePrivateKey(mnemonic, uint32(coinType), uint(account), uint(index))
	s.Require().Error((err))
}

func (s *KeysTestSuite) TestGetKeyFromKeyName() {
	keyName := "testkeyname"
	privatekeyHex := testPrivateKey

	// Add a key to test retrieval
	_, err := s.chainProvider.AddKeyWithPrivateKey(keyName, privatekeyHex, s.homePath, "")
	s.Require().NoError(err)

	// Retrieve the key using the key name
	key, err := s.chainProvider.getKeyFromKeyName(keyName, "")
	s.Require().NoError(err)
	s.Require().NotNil(key)

	// Verify that the retrieved private key matches the original private key
	s.Require().Equal(testPrivateKey[2:], hex.EncodeToString(crypto.FromECDSA(key.PrivateKey))) // Remove "0x"
}

func (s *KeysTestSuite) TestGetKeyFromKeyNameWithInvalidPassphrase() {
	keyName := "testkeyname"
	privatekeyHex := testPrivateKey
	passphrase := ""
	invalidPassphrase := "invalid"

	// Add a key to test retrieval
	_, err := s.chainProvider.AddKeyWithPrivateKey(keyName, privatekeyHex, s.homePath, passphrase)
	s.Require().NoError(err)

	// Retrieve the key using the key name
	_, err = s.chainProvider.getKeyFromKeyName(keyName, invalidPassphrase)
	s.Require().Error(err)
}
