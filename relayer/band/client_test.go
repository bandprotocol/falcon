package band_test

import (
	"context"
	"testing"

	"github.com/bandprotocol/falcon/internal/relayertest/mocks"
	"github.com/bandprotocol/falcon/relayer"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type AppTestSuite struct {
	suite.Suite

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
}


func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}

