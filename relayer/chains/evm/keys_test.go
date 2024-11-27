package evm

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/relayer/chains"
	chaintypes "github.com/bandprotocol/falcon/relayer/chains/types"
)

const (
	privateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	address    = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
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

	// mock objects.
	s.log = log

	chainName := "testnet"

	client := NewClient(chainName, evmCfg, log)

	s.chainProvider, err = NewEVMChainProvider(chainName, client, evmCfg, log, tmpDir)
	s.Require().NoError(err)

	s.ctx = context.Background()
	s.homePath = tmpDir
}

func (s *KeysTestSuite) TestAddKeyWithPrivateKey() {
	keyName := "testkey"
	privatekeyHex := "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	actual, err := s.chainProvider.AddKeyWithPrivateKey(keyName, privatekeyHex, s.homePath, "")
	s.Require().NoError(err)

	expected := chaintypes.NewKey("", address, "")
	s.Require().Equal(expected, actual)

	addr, err := HexToAddress(address)
	s.Require().NoError(err)

	// check that key actually added in keystore
	s.Require().True(s.chainProvider.KeyStore.HasAddress(addr))

	// check that key info actually stored in local disk
	keyInfo, err := LoadKeyInfo(s.homePath, s.chainProvider.ChainName)
	s.Require().NoError(err)

	_, exist := keyInfo[keyName]
	s.Require().True(exist)
}

func (s *KeysTestSuite) TestAddKeyWithNoMnemonic() {
	keyName := "testkey"
	mnemonic := ""
	coinType := 60
	account := 0
	index := 0

	actual, err := s.chainProvider.AddKeyWithMnemonic(
		keyName,
		mnemonic,
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

func (s *KeysTestSuite) TestAddKeyWithGivenMnemonic() {
	keyName := "testkey"
	mnemonic := "evil cool swamp nurse emotion dumb lecture foam stamp cigar bamboo arctic leaf twin brand sight soda drill december dial raccoon race seek expose"
	coinType := 60
	account := 0
	index := 0

	actual, err := s.chainProvider.AddKeyWithMnemonic(
		keyName,
		mnemonic,
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

func (s *KeysTestSuite) TestAddKeyWithDuplicatePrivateKey() {
	keyName := "testkey"
	privatekeyHex := "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	actual, err := s.chainProvider.AddKeyWithPrivateKey(keyName, privatekeyHex, s.homePath, "")
	s.Require().NoError(err)

	expected := chaintypes.NewKey("", address, "")
	s.Require().Equal(expected, actual)

	addr, err := HexToAddress(address)
	s.Require().NoError(err)

	// check that key actually added in keystore
	s.Require().True(s.chainProvider.KeyStore.HasAddress(addr))

	// check that key info actually stored in local disk
	keyInfo, err := LoadKeyInfo(s.homePath, s.chainProvider.ChainName)
	s.Require().NoError(err)

	_, exist := keyInfo[keyName]
	s.Require().True(exist)

	_, err = s.chainProvider.AddKeyWithPrivateKey(keyName, privatekeyHex, s.homePath, "")
	s.Require().ErrorContains(err, "account already exists")
}

func (s *KeysTestSuite) TestAddKeyWithDuplicateMnemonic() {
	keyName := "testkey"
	mnemonic := "evil cool swamp nurse emotion dumb lecture foam stamp cigar bamboo arctic leaf twin brand sight soda drill december dial raccoon race seek expose"
	coinType := 60
	account := 0
	index := 0

	actual, err := s.chainProvider.AddKeyWithMnemonic(
		keyName,
		mnemonic,
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

	_, err = s.chainProvider.AddKeyWithMnemonic(
		keyName,
		mnemonic,
		s.homePath,
		uint32(coinType),
		uint(account),
		uint(index),
		"",
	)
	s.Require().ErrorContains(err, "account already exists")
}

func (s *KeysTestSuite) TestAddKeyWithInvalidPrivateKeyFormat() {
	keyName := "testkey"
	privatekeyHex := "privatekey"

	_, err := s.chainProvider.AddKeyWithPrivateKey(keyName, privatekeyHex, s.homePath, "")
	s.Require().Error(err)
}

func (s *KeysTestSuite) TestIsKeyNameExist() {
	keyName := "existingkey"
	privatekeyHex := privateKey
	s.Require().False(s.chainProvider.IsKeyNameExist(keyName))

	// Add a key to test existence
	_, err := s.chainProvider.AddKeyWithPrivateKey(keyName, privatekeyHex, s.homePath, "")
	s.Require().NoError(err)

	// Check that the key exists
	s.Require().True(s.chainProvider.IsKeyNameExist(keyName))

	// Check that a non-existent key returns false
	s.Require().False(s.chainProvider.IsKeyNameExist("nonexistentkey"))
}

func (s *KeysTestSuite) TestDeleteKey() {
	keyName := "deletablekey"
	privatekeyHex := privateKey

	// Add a key to delete
	_, err := s.chainProvider.AddKeyWithPrivateKey(keyName, privatekeyHex, s.homePath, "")
	s.Require().NoError(err)

	// Delete the key
	err = s.chainProvider.DeleteKey(s.homePath, keyName, "")
	s.Require().NoError(err)

	// Ensure the key is no longer in the KeyInfo or KeyStore
	s.Require().False(s.chainProvider.IsKeyNameExist(keyName))

	addr, err := HexToAddress(address)
	s.Require().NoError(err)
	s.Require().False(s.chainProvider.KeyStore.HasAddress(addr))
}

func (s *KeysTestSuite) TestExportPrivateKey() {
	keyName := "exportkey"
	privatekeyHex := privateKey

	// Add a key to export
	_, err := s.chainProvider.AddKeyWithPrivateKey(keyName, privatekeyHex, s.homePath, "")
	s.Require().NoError(err)

	// Export the private key
	exportedKey, err := s.chainProvider.ExportPrivateKey(keyName, "")
	s.Require().NoError(err)

	s.Require().Equal(ConvertPrivateKeyStrToHex(privatekeyHex), ConvertPrivateKeyStrToHex(exportedKey))
}

func (s *KeysTestSuite) TestListKeys() {
	// Add multiple keys
	keyName1 := "key1"
	keyName2 := "key2"
	coinType := 60
	account := 0
	index := 0

	key1, err := s.chainProvider.AddKeyWithMnemonic(
		keyName1,
		"",
		s.homePath,
		uint32(coinType),
		uint(account),
		uint(index),
		"",
	)
	s.Require().NoError(err)

	key2, err := s.chainProvider.AddKeyWithMnemonic(
		keyName2,
		"",
		s.homePath,
		uint32(coinType),
		uint(account),
		uint(index),
		"",
	)
	s.Require().NoError(err)

	// List all keys
	actual := s.chainProvider.Listkeys()

	expected := []*chaintypes.Key{
		chaintypes.NewKey("", key1.Address, keyName1),
		chaintypes.NewKey("", key2.Address, keyName2),
	}

	s.Require().Equal(expected, actual)
}

func (s *KeysTestSuite) TestShowKey() {
	keyName := "showkey"
	privatekeyHex := privateKey

	// Add a key to show
	_, err := s.chainProvider.AddKeyWithPrivateKey(keyName, privatekeyHex, s.homePath, "")
	s.Require().NoError(err)

	// Show the key
	address := s.chainProvider.ShowKey(keyName)
	s.Require().Equal(address, address)
}

func (s *KeysTestSuite) TestStorePrivateKey() {
	privateKeyECDSA, err := crypto.HexToECDSA(ConvertPrivateKeyStrToHex(privateKey)) // Remove "0x" prefix
	s.Require().NoError(err)

	// Store the private key in the keystore
	account, err := s.chainProvider.storePrivateKey(privateKeyECDSA, "")
	s.Require().NoError(err)
	s.Require().NotNil(account)

	// Verify that the key exists in the keystore
	addr := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)
	s.Require().True(s.chainProvider.KeyStore.HasAddress(addr))
}

func (s *KeysTestSuite) TestStoreKeyInfo() {
	keyName := "testkeyinfo"
	privatekeyHex := privateKey

	// Add a key to simulate storing key info
	_, err := s.chainProvider.AddKeyWithPrivateKey(keyName, privatekeyHex, s.homePath, "")
	s.Require().NoError(err)

	// Verify that key info is correctly stored in the file
	keyInfo, err := LoadKeyInfo(s.homePath, s.chainProvider.ChainName)
	s.Require().NoError(err)

	address := s.chainProvider.ShowKey(keyName)
	s.Require().Equal(keyInfo[keyName], address)
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

func (s *KeysTestSuite) TestGetKeyFromKeyName() {
	keyName := "testkeyname"
	privatekeyHex := privateKey

	// Add a key to test retrieval
	_, err := s.chainProvider.AddKeyWithPrivateKey(keyName, privatekeyHex, s.homePath, "")
	s.Require().NoError(err)

	// Retrieve the key using the key name
	key, err := s.chainProvider.getKeyFromKeyName(keyName, "")
	s.Require().NoError(err)
	s.Require().NotNil(key)

	// Verify that the retrieved private key matches the original private key
	s.Require().Equal(privateKey[2:], hex.EncodeToString(crypto.FromECDSA(key.PrivateKey))) // Remove "0x"
}
