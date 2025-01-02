package relayer_test

import (
	"context"
	"os"
	"path"
	"testing"
	"time"

	cmbytes "github.com/cometbft/cometbft/libs/bytes"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/internal/relayertest"
	"github.com/bandprotocol/falcon/internal/relayertest/mocks"
	"github.com/bandprotocol/falcon/relayer"
	"github.com/bandprotocol/falcon/relayer/band"
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	"github.com/bandprotocol/falcon/relayer/chains"
	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/types"
)

type AppTestSuite struct {
	suite.Suite

	app                 *relayer.App
	ctx                 context.Context
	chainProviderConfig *mocks.MockChainProviderConfig
	chainProvider       *mocks.MockChainProvider
	client              *mocks.MockClient
}

// SetupTest sets up the test suite by creating a temporary directory and declare mock objects.
func (s *AppTestSuite) SetupTest() {
	tmpDir := s.T().TempDir()
	ctrl := gomock.NewController(s.T())

	log, err := zap.NewDevelopment()
	s.Require().NoError(err)

	// mock objects.
	s.chainProviderConfig = mocks.NewMockChainProviderConfig(ctrl)
	s.chainProvider = mocks.NewMockChainProvider(ctrl)
	s.client = mocks.NewMockClient(ctrl)

	s.chainProviderConfig.EXPECT().
		NewChainProvider("testnet_evm", log, tmpDir, false).
		Return(s.chainProvider, nil).
		AnyTimes()

	s.chainProvider.EXPECT().Init(gomock.Any()).Return(nil).AnyTimes()

	cfg := relayer.Config{
		BandChain: band.Config{},
		TargetChains: map[string]chains.ChainProviderConfig{
			"testnet_evm": s.chainProviderConfig,
		},
		Global: relayer.GlobalConfig{},
	}
	s.ctx = context.Background()

	s.app = relayer.NewApp(log, tmpDir, false, &cfg)

	err = s.app.Init(s.ctx, "", "")
	s.app.BandClient = s.client
	s.Require().NoError(err)
}

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}

func (s *AppTestSuite) TestInitConfig() {
	s.app.Config = nil
	customCfgPath := ""

	err := s.app.InitConfigFile(s.app.HomePath, customCfgPath)
	s.Require().NoError(err)

	cfgPath := path.Join(s.app.HomePath, "config", "config.toml")
	s.Require().FileExists(cfgPath)

	// read the file
	actualByte, err := os.ReadFile(cfgPath)
	s.Require().NoError(err)

	// marshal default config
	expect := relayer.DefaultConfig()
	expectBytes, err := toml.Marshal(expect)
	s.Require().NoError(err)

	s.Require().Equal(string(expectBytes), string(actualByte))
}

func (s *AppTestSuite) TestInitExistingConfig() {
	s.app.Config = nil
	customCfgPath := ""

	err := s.app.InitConfigFile(s.app.HomePath, customCfgPath)
	s.Require().NoError(err)

	// second time should fail
	err = s.app.InitConfigFile(s.app.HomePath, customCfgPath)
	s.Require().ErrorContains(err, "config already exists:")
}

func (s *AppTestSuite) TestInitCustomConfig() {
	s.app.Config = nil
	customCfgPath := path.Join(s.app.HomePath, "custom.toml")

	// Create custom config file
	cfg := `
		[target_chains]

		[global]
		checking_packet_interval = 60000000000
	
		[bandchain]
		rpc_endpoints = ['http://localhost:26659']
		timeout = 50
	`
	// write file
	err := os.WriteFile(customCfgPath, []byte(cfg), 0o600)
	s.Require().NoError(err)

	err = s.app.InitConfigFile(s.app.HomePath, customCfgPath)
	s.Require().NoError(err)

	s.Require().FileExists(path.Join(s.app.HomePath, "config", "config.toml"))

	// read the file
	b, err := os.ReadFile(path.Join(s.app.HomePath, "config", "config.toml"))
	s.Require().NoError(err)

	// unmarshal data
	actual := relayer.Config{}
	err = toml.Unmarshal(b, &actual)
	s.Require().NoError(err)

	expect := relayer.Config{
		BandChain: band.Config{
			RpcEndpoints: []string{"http://localhost:26659"},
			Timeout:      50,
		},
		TargetChains: nil,
		Global: relayer.GlobalConfig{
			CheckingPacketInterval: time.Minute,
		},
	}

	s.Require().Equal(expect, actual)
}

func (s *AppTestSuite) TestQueryTunnelInfo() {
	tunnelBandInfo := bandtypes.NewTunnel(1, 1, "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2", "testnet_evm", false)
	tunnelChainInfo := chainstypes.NewTunnel(1, "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2", false)

	s.client.EXPECT().
		GetTunnel(s.ctx, uint64(1)).
		Return(tunnelBandInfo, nil)

	s.chainProvider.EXPECT().
		QueryTunnelInfo(s.ctx, uint64(1), "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2").
		Return(tunnelChainInfo, nil)

	tunnel, err := s.app.QueryTunnelInfo(s.ctx, 1)
	bandChainInfo := bandtypes.NewTunnel(1, 1, "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2", "testnet_evm", false)

	expected := types.NewTunnel(
		bandChainInfo,
		tunnelChainInfo,
	)
	s.Require().NoError(err)
	s.Require().Equal(expected, tunnel)
}

func (s *AppTestSuite) TestQueryTunnelInfoNotSupportedChain() {
	s.app.Config.TargetChains = nil
	err := s.app.Init(s.ctx, "", "")

	s.Require().NoError(err)

	tunnelBandInfo := bandtypes.NewTunnel(1, 1, "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2", "testnet_evm", false)
	s.client.EXPECT().
		GetTunnel(s.ctx, uint64(1)).
		Return(tunnelBandInfo, nil)
	s.app.BandClient = s.client

	tunnel, err := s.app.QueryTunnelInfo(s.ctx, 1)

	expected := types.NewTunnel(
		tunnelBandInfo,
		nil,
	)
	s.Require().NoError(err)
	s.Require().Equal(expected, tunnel)
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
		"SIGNING_STATUS_SUCCESS",
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
		GetTunnelPacket(s.ctx, uint64(1), uint64(1)).
		Return(tunnelPacketBandInfo, nil)

	// Call the function under test
	packet, err := s.app.QueryTunnelPacketInfo(s.ctx, 1, 1)

	// Create the expected packet structure for comparison
	expected := bandtypes.NewPacket(1, 1, signalPrices, signingInfo, nil)

	// Assertions
	s.Require().NoError(err)
	s.Require().Equal(expected, packet)
}

func (s *AppTestSuite) TestAddChainConfig() {
	s.app.Config = nil
	// create new chain config file
	chainCfgPath := path.Join(s.app.HomePath, "chain_config.toml")
	chainName := "testnet"

	// write chain config file
	err := os.WriteFile(chainCfgPath, []byte(relayertest.ChainCfgText), 0o600)
	s.Require().NoError(err)

	s.Require().FileExists(chainCfgPath)

	// init chain config file
	customCfgPath := ""
	err = s.app.InitConfigFile(s.app.HomePath, customCfgPath)
	s.Require().NoError(err)

	s.Require().FileExists(path.Join(s.app.HomePath, "config", "config.toml"))

	// load config
	err = s.app.LoadConfigFile()
	s.Require().NoError(err)

	err = s.app.AddChainConfig(chainName, chainCfgPath)
	s.Require().NoError(err)

	expectedBytes := []byte(relayertest.DefaultCfgTextWithChainCfg)
	actualBytes, err := os.ReadFile(path.Join(s.app.HomePath, "config", "config.toml"))

	s.Require().NoError(err)
	s.Require().Equal(relayertest.DefaultCfgTextWithChainCfg, string(actualBytes))

	s.Require().Equal(expectedBytes, actualBytes)
}

func (s *AppTestSuite) TestDeleteChainConfig() {
	s.app.Config = nil
	customCfgPath := path.Join(s.app.HomePath, "custom.toml")

	// write file
	err := os.WriteFile(customCfgPath, []byte(relayertest.DefaultCfgTextWithChainCfg), 0o600)
	s.Require().NoError(err)

	err = s.app.InitConfigFile(s.app.HomePath, customCfgPath)
	s.Require().NoError(err)

	// load config file
	err = s.app.LoadConfigFile()
	s.Require().NoError(err)

	// delete chain config by given chain name
	chainName := "testnet"
	err = s.app.DeleteChainConfig(chainName)
	s.Require().NoError(err)

	expectedBytes := []byte(relayertest.DefaultCfgText)

	actualBytes, err := os.ReadFile(path.Join(s.app.HomePath, "config", "config.toml"))
	s.Require().NoError(err)

	s.Require().Equal(expectedBytes, actualBytes)
}

func (s *AppTestSuite) TestGetChainConfig() {
	s.app.Config = nil
	customCfgPath := path.Join(s.app.HomePath, "custom.toml")

	// write file
	err := os.WriteFile(customCfgPath, []byte(relayertest.DefaultCfgTextWithChainCfg), 0o600)
	s.Require().NoError(err)

	err = s.app.InitConfigFile(s.app.HomePath, customCfgPath)
	s.Require().NoError(err)

	// load config file
	err = s.app.LoadConfigFile()
	s.Require().NoError(err)

	chainName := "testnet"
	actual, err := s.app.GetChainConfig(chainName)
	s.Require().NoError(err)

	expect := relayertest.CustomCfg.TargetChains[chainName]

	s.Require().Equal(expect, actual)
}
