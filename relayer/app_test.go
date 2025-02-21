package relayer_test

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	cmbytes "github.com/cometbft/cometbft/libs/bytes"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	tsstypes "github.com/bandprotocol/falcon/internal/bandchain/tss"
	"github.com/bandprotocol/falcon/internal/relayertest"
	"github.com/bandprotocol/falcon/internal/relayertest/mocks"
	"github.com/bandprotocol/falcon/relayer"
	"github.com/bandprotocol/falcon/relayer/band"
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/evm"
	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/config"
	"github.com/bandprotocol/falcon/relayer/types"
)

type AppTestSuite struct {
	suite.Suite

	app                 *relayer.App
	chainProviderConfig *mocks.MockChainProviderConfig
	chainProvider       *mocks.MockChainProvider
	client              *mocks.MockClient
	mockStore           *mocks.MockStore

	passphrase       string
	hashedPassphrase []byte
}

// SetupTest sets up the test suite by creating a temporary directory and declare mock objects.
func (s *AppTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	log := zap.NewNop()

	// mock objects.
	s.chainProviderConfig = mocks.NewMockChainProviderConfig(ctrl)
	s.chainProvider = mocks.NewMockChainProvider(ctrl)
	s.client = mocks.NewMockClient(ctrl)
	s.mockStore = mocks.NewMockStore(ctrl)

	cfg := config.Config{
		BandChain: band.Config{
			RpcEndpoints:               []string{"http://localhost:26659"},
			LivelinessCheckingInterval: 5 * time.Minute,
		},
		TargetChains: map[string]chains.ChainProviderConfig{
			"testnet_evm": s.chainProviderConfig,
		},
		Global: config.GlobalConfig{},
	}

	s.passphrase = "secret"

	h := sha256.New()
	h.Write([]byte(s.passphrase))
	s.hashedPassphrase = h.Sum(nil)
	s.mockStore.EXPECT().GetHashedPassphrase().Return(s.hashedPassphrase, nil).AnyTimes()

	s.app = &relayer.App{
		Log:    log,
		Config: &cfg,
		Store:  s.mockStore,
		TargetChains: map[string]chains.ChainProvider{
			"testnet_evm": s.chainProvider,
		},
		BandClient: s.client,
		Passphrase: s.passphrase,
	}
}

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}

func (s *AppTestSuite) TestInitConfig() {
	customCfg := &config.Config{
		BandChain: band.Config{
			RpcEndpoints: []string{"http://localhost:26659"},
			Timeout:      50,
		},
		TargetChains: map[string]chains.ChainProviderConfig{},
		Global: config.GlobalConfig{
			CheckingPacketInterval: time.Minute,
		},
	}

	testcases := []struct {
		name       string
		preprocess func()
		in         *config.Config
		err        error
	}{
		{
			name: "success - default",
			in:   nil,
			preprocess: func() {
				s.mockStore.EXPECT().HasConfig().Return(false, nil)
				s.mockStore.EXPECT().SaveConfig(config.DefaultConfig()).Return(nil)
			},
		},
		{
			name: "config already exists",
			preprocess: func() {
				s.mockStore.EXPECT().HasConfig().Return(true, nil)
			},
			in:  nil,
			err: fmt.Errorf("config already exists"),
		},
		{
			name: "init config from specific file",
			preprocess: func() {
				s.mockStore.EXPECT().HasConfig().Return(false, nil)
				s.mockStore.EXPECT().SaveConfig(customCfg).Return(nil)
			},
			in: customCfg,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			if tc.preprocess != nil {
				tc.preprocess()
			}

			err := s.app.SaveConfig(tc.in)

			if tc.err != nil {
				s.Require().ErrorContains(err, tc.err.Error())
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *AppTestSuite) TestAddChainConfig() {
	tmpDir := s.T().TempDir()
	newHomePath := path.Join(tmpDir, "new_folder")
	err := os.Mkdir(newHomePath, os.ModePerm)
	s.Require().NoError(err)

	type Input struct {
		chainName   string
		cfgPath     string
		existingCfg *config.Config
	}
	testcases := []struct {
		name       string
		preprocess func()
		in         Input
		err        error
	}{
		{
			name: "success",
			in: Input{
				chainName:   "testnet",
				cfgPath:     path.Join(newHomePath, "chain_config.toml"),
				existingCfg: config.DefaultConfig(),
			},
			preprocess: func() {
				chainCfgPath := path.Join(newHomePath, "chain_config.toml")
				err := os.WriteFile(chainCfgPath, []byte(relayertest.ChainCfgText), 0o600)
				s.Require().NoError(err)

				cfg, err := config.ParseConfig([]byte(relayertest.DefaultCfgTextWithChainCfg))
				s.Require().NoError(err)

				s.mockStore.EXPECT().SaveConfig(cfg).Return(nil)
			},
		},
		{
			name: "invalid chain type",
			in: Input{
				chainName:   "testnet",
				cfgPath:     path.Join(newHomePath, "chain_config.toml"),
				existingCfg: config.DefaultConfig(),
			},
			preprocess: func() {
				chainCfgPath := path.Join(newHomePath, "chain_config.toml")
				err := os.WriteFile(chainCfgPath, []byte(relayertest.ChainCfgInvalidChainTypeText), 0o600)
				s.Require().NoError(err)
			},
			err: fmt.Errorf("unsupported chain type"),
		},
		{
			name: "existing chain name",
			in: Input{
				chainName: "testnet",
				cfgPath:   path.Join(newHomePath, "chain_config.toml"),
				existingCfg: &config.Config{
					TargetChains: map[string]chains.ChainProviderConfig{
						"testnet": &evm.EVMChainProviderConfig{},
					},
				},
			},
			preprocess: func() {
				chainCfgPath := path.Join(newHomePath, "chain_config.toml")
				err := os.WriteFile(chainCfgPath, []byte(relayertest.ChainCfgText), 0o600)
				s.Require().NoError(err)
			},
			err: fmt.Errorf("existing chain name :"),
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			if tc.in.existingCfg != nil {
				s.app.Config = tc.in.existingCfg
			}

			if tc.preprocess != nil {
				tc.preprocess()
			}

			err = s.app.AddChainConfig(tc.in.chainName, tc.in.cfgPath)

			if tc.err != nil {
				s.Require().ErrorContains(err, tc.err.Error())
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *AppTestSuite) TestDeleteChainConfig() {
	var err error
	s.app.Config, err = config.ParseConfig([]byte(relayertest.DefaultCfgTextWithChainCfg))
	s.Require().NoError(err)

	testcases := []struct {
		name       string
		in         string
		preprocess func()
		out        string
		err        error
	}{
		{
			name: "success",
			in:   "testnet",
			preprocess: func() {
				s.mockStore.EXPECT().SaveConfig(config.DefaultConfig()).Return(nil)
			},
			out: relayertest.DefaultCfgText,
		},
		{
			name: "not existing chain name",
			in:   "testnet2",
			err:  fmt.Errorf("not existing chain name"),
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			if tc.preprocess != nil {
				tc.preprocess()
			}

			err = s.app.DeleteChainConfig(tc.in)
			if tc.err != nil {
				s.Require().ErrorContains(err, tc.err.Error())
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *AppTestSuite) TestGetChainConfig() {
	testcases := []struct {
		name string
		in   string
		err  error
		out  chains.ChainProviderConfig
	}{
		{
			name: "success",
			in:   "testnet_evm",
			out:  s.chainProviderConfig,
		},
		{
			name: "not existing chain name",
			in:   "testnet_evm2",
			err:  fmt.Errorf("not existing chain name"),
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			actual, err := s.app.GetChainConfig(tc.in)

			if tc.err != nil {
				s.Require().ErrorContains(err, tc.err.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.out, actual)
			}
		})
	}
}

func (s *AppTestSuite) TestQueryTunnelInfo() {
	mockTunnelBandInfo := bandtypes.NewTunnel(1, 1, "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2", "testnet_evm", false)
	mockTunnelBandInfoNoChain := bandtypes.NewTunnel(1, 1, "0xmock", "unknown_chain", false)
	mockTunnelChainInfo := chainstypes.NewTunnel(1, "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2", false)

	testcases := []struct {
		name       string
		preprocess func()
		in         uint64
		out        *types.Tunnel
		err        error
	}{
		{
			name: "success",
			preprocess: func() {
				s.client.EXPECT().
					GetTunnel(gomock.Any(), uint64(1)).
					Return(mockTunnelBandInfo, nil)
				s.chainProvider.EXPECT().
					QueryTunnelInfo(gomock.Any(), uint64(1), "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2").
					Return(mockTunnelChainInfo, nil)
			},
			in:  1,
			out: types.NewTunnel(mockTunnelBandInfo, mockTunnelChainInfo),
		},
		{
			name: "cannot query chain info",
			preprocess: func() {
				s.client.EXPECT().
					GetTunnel(gomock.Any(), uint64(1)).
					Return(mockTunnelBandInfo, nil)
				s.chainProvider.EXPECT().
					QueryTunnelInfo(gomock.Any(), uint64(1), "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2").
					Return(nil, fmt.Errorf("cannot connect to chain"))
			},
			in:  1,
			err: fmt.Errorf("cannot connect to chain"),
		},
		{
			name: "no chain provider",
			preprocess: func() {
				s.client.EXPECT().
					GetTunnel(gomock.Any(), uint64(1)).
					Return(mockTunnelBandInfoNoChain, nil)
			},
			in:  1,
			out: types.NewTunnel(mockTunnelBandInfoNoChain, nil),
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			if tc.preprocess != nil {
				tc.preprocess()
			}

			tunnel, err := s.app.QueryTunnelInfo(context.Background(), tc.in)

			if tc.err != nil {
				s.Require().ErrorContains(err, tc.err.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.out, tunnel)
			}
		})
	}
}

func (s *AppTestSuite) TestQueryTunnelPacketInfo() {
	signalPrices := []bandtypes.SignalPrice{
		{SignalID: "signal1", Price: 100},
		{SignalID: "signal2", Price: 200},
	}

	// Create a mock EVMSignature
	evmSignature := bandtypes.NewEVMSignature(
		cmbytes.HexBytes("0x1234"),
		cmbytes.HexBytes("0xabcd"),
	)

	// Create mock signing information
	signingInfo := bandtypes.NewSigning(
		1,
		cmbytes.HexBytes("0xdeadbeef"),
		evmSignature,
		tsstypes.SIGNING_STATUS_SUCCESS,
	)

	// Create the expected Packet object
	tunnelPacketBandInfo := bandtypes.NewPacket(
		1,
		1,
		signalPrices,
		signingInfo,
		nil,
	)

	// Set up the mock expectation
	s.client.EXPECT().
		GetTunnelPacket(gomock.Any(), uint64(1), uint64(1)).
		Return(tunnelPacketBandInfo, nil)

	// Call the function under test
	packet, err := s.app.QueryTunnelPacketInfo(context.Background(), 1, 1)

	// Create the expected packet structure for comparison
	expected := bandtypes.NewPacket(1, 1, signalPrices, signingInfo, nil)

	// Assertions
	s.Require().NoError(err)
	s.Require().Equal(expected, packet)
}

func (s *AppTestSuite) TestInitPassphrase() {
	ctrl := gomock.NewController(s.T())
	newStoreMock := mocks.NewMockStore(ctrl)

	s.app.Store = newStoreMock

	newStoreMock.EXPECT().
		SaveHashedPassphrase([]byte{
			194, 83, 183, 41, 238, 49, 98, 232, 230, 229, 194,
			192, 115, 133, 235, 215, 215, 206, 160, 68, 116,
			34, 59, 169, 179, 24, 231, 151, 191, 178, 90, 202,
		}).
		Return(nil)

	// Call InitPassphrase
	err := s.app.SavePassphrase("new_passphrase")
	s.Require().NoError(err)
	s.Require().Equal("new_passphrase", s.app.Passphrase)
}

func (s *AppTestSuite) TestAddKey() {
	testcases := []struct {
		name       string
		chainName  string
		keyName    string
		mnemonic   string
		privateKey string
		coinType   uint32
		account    uint
		index      uint
		err        error
		out        *chainstypes.Key
		preprocess func()
	}{
		{
			name:       "success - private key",
			chainName:  "testnet_evm",
			keyName:    "testkey",
			privateKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", // anvil
			coinType:   60,
			out:        chainstypes.NewKey("", "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", ""),
			preprocess: func() {
				s.chainProvider.EXPECT().
					AddKey(
						"testkey",
						"",
						"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
						uint32(60),
						uint(0),
						uint(0),
					).
					Return(chainstypes.NewKey("", "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", ""), nil)
			},
		},
		{
			name:       "error from AddKey",
			chainName:  "testnet_evm",
			keyName:    "testkey",
			privateKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", // anvil
			coinType:   60,
			preprocess: func() {
				s.chainProvider.EXPECT().
					AddKey(
						"testkey",
						"",
						"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
						uint32(60),
						uint(0),
						uint(0),
					).
					Return(nil, fmt.Errorf("add key error"))
			},
			err: fmt.Errorf("add key error"),
		},
		{
			name:       "chain name does not exist",
			chainName:  "testnet_evm2",
			keyName:    "testkey",
			privateKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", // anvil
			coinType:   60,
			err:        fmt.Errorf("chain name does not exist:"),
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			if tc.preprocess != nil {
				tc.preprocess()
			}

			actual, err := s.app.AddKey(
				tc.chainName,
				tc.keyName,
				tc.mnemonic,
				tc.privateKey,
				tc.coinType,
				tc.account,
				tc.index,
			)

			if tc.err != nil {
				s.Require().ErrorContains(err, tc.err.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.out, actual)
			}
		})
	}
}

func (s *AppTestSuite) TestDeleteKey() {
	testcases := []struct {
		name       string
		chainName  string
		keyName    string
		err        error
		preprocess func()
	}{
		{
			name:      "success",
			chainName: "testnet_evm",
			keyName:   "testkey",
			preprocess: func() {
				s.chainProvider.EXPECT().
					DeleteKey("testkey").
					Return(nil)
			},
		},
		{
			name:      "error delete key",
			chainName: "testnet_evm",
			keyName:   "testkey",
			preprocess: func() {
				s.chainProvider.EXPECT().
					DeleteKey("testkey").
					Return(fmt.Errorf("delete key error"))
			},
			err: fmt.Errorf("delete key error"),
		},
		{
			name:      "chain name does not exist",
			chainName: "testnet_evm2",
			keyName:   "testkey",
			err:       fmt.Errorf("chain name does not exist:"),
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			if tc.preprocess != nil {
				tc.preprocess()
			}

			err := s.app.DeleteKey(tc.chainName, tc.keyName)

			if tc.err != nil {
				s.Require().ErrorContains(err, tc.err.Error())
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *AppTestSuite) TestExportKey() {
	testcases := []struct {
		name       string
		chainName  string
		keyName    string
		out        string
		err        error
		preprocess func()
	}{
		{
			name:      "success",
			chainName: "testnet_evm",
			keyName:   "testkey",
			preprocess: func() {
				s.chainProvider.EXPECT().
					ExportPrivateKey("testkey").
					Return("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", nil)
			},
			out: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		},
		{
			name:      "error export private key",
			chainName: "testnet_evm",
			keyName:   "testkey",
			preprocess: func() {
				s.chainProvider.EXPECT().
					ExportPrivateKey("testkey").
					Return("", fmt.Errorf("export key error"))
			},
			err: fmt.Errorf("export key error"),
		},
		{
			name:      "chain name does not exist",
			chainName: "testnet_evm2",
			keyName:   "testkey",
			err:       fmt.Errorf("chain name does not exist:"),
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			if tc.preprocess != nil {
				tc.preprocess()
			}

			actual, err := s.app.ExportKey(tc.chainName, tc.keyName)
			if tc.err != nil {
				s.Require().ErrorContains(err, tc.err.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.out, actual)
			}
		})
	}
}

func (s *AppTestSuite) TestListKeys() {
	testcases := []struct {
		name       string
		in         string
		preprocess func()
		err        error
		out        []*chainstypes.Key
	}{
		{
			name: "success",
			in:   "testnet_evm",
			preprocess: func() {
				s.chainProvider.EXPECT().
					ListKeys().
					Return([]*chainstypes.Key{
						chainstypes.NewKey("", "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", "testkey1"),
						chainstypes.NewKey("", "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92267", "testkey2"),
					})
			},
			out: []*chainstypes.Key{
				chainstypes.NewKey("", "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", "testkey1"),
				chainstypes.NewKey("", "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92267", "testkey2"),
			},
		},
		{
			name: "chain name does not exist",
			in:   "testnet_evm2",
			err:  fmt.Errorf("chain name does not exist:"),
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			if tc.preprocess != nil {
				tc.preprocess()
			}

			actual, err := s.app.ListKeys(tc.in)
			if tc.err != nil {
				s.Require().ErrorContains(err, tc.err.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(actual, tc.out)
			}
		})
	}
}

func (s *AppTestSuite) TestShowKey() {
	testcases := []struct {
		name       string
		chainName  string
		keyName    string
		preprocess func()
		err        error
		out        string
	}{
		{
			name:      "success",
			chainName: "testnet_evm",
			keyName:   "testkey",
			preprocess: func() {
				s.chainProvider.EXPECT().
					ShowKey("testkey").
					Return("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92267", nil)
			},
			out: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92267",
		},
		{
			name:      "show key error",
			chainName: "testnet_evm",
			keyName:   "testkey",
			preprocess: func() {
				s.chainProvider.EXPECT().
					ShowKey("testkey").
					Return("", fmt.Errorf("show key error"))
			},
			out: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92267",
			err: fmt.Errorf("show key error"),
		},
		{
			name:      "chain name does not exist",
			chainName: "testnet_evm2",
			keyName:   "testkey",
			err:       fmt.Errorf("chain name does not exist:"),
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			if tc.preprocess != nil {
				tc.preprocess()
			}

			actual, err := s.app.ShowKey(tc.chainName, tc.keyName)
			if tc.err != nil {
				s.Require().ErrorContains(err, tc.err.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(actual, tc.out)
			}
		})
	}
}

func (s *AppTestSuite) TestValidatePassphraseInvalidPassphrase() {
	testcases := []struct {
		name          string
		envPassphrase string
		err           error
	}{
		{name: "valid", envPassphrase: "secret", err: nil},
		{name: "invalid", envPassphrase: "invalid", err: fmt.Errorf("invalid passphrase")},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			err := s.app.ValidatePassphrase(tc.envPassphrase)
			if tc.err != nil {
				s.Require().ErrorContains(err, tc.err.Error())
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
