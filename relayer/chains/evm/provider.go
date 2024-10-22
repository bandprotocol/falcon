package evm

import (
	"context"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/relayer/chains/types"
)

// EVMChainProvider is the struct that handles interactions with the EVM chain.
type EVMChainProvider struct {
	Cfg *EVMChainProviderConfig

	ChainName string
	Client    Client

	TunnelRouterAddress gethcommon.Address
	TunnelRouterABI     abi.ABI

	Log *zap.Logger
}

// NewEVMChainProvider creates a new EVM chain provider.
func NewEVMChainProvider(
	chainName string,
	client Client,
	cfg *EVMChainProviderConfig,
	log *zap.Logger,
) (*EVMChainProvider, error) {
	// load abis here
	abi, err := abi.JSON(strings.NewReader(gasPriceTunnelRouterABI))
	if err != nil {
		return nil, err
	}

	addr, err := HexToAddress(cfg.TunnelRouterAddress)
	if err != nil {
		return nil, err
	}

	// check client connection
	if err := client.Connect(); err != nil {
		return nil, err
	}

	return &EVMChainProvider{
		Cfg:                 cfg,
		ChainName:           chainName,
		Client:              client,
		TunnelRouterAddress: addr,
		TunnelRouterABI:     abi,
		Log:                 log,
	}, nil
}

// Connect connects to the EVM chain.
func (cp *EVMChainProvider) Connect() error {
	return cp.Client.Connect()
}

// QueryTunnelInfo queries the tunnel info from the tunnel router contract.
func (cp *EVMChainProvider) QueryTunnelInfo(
	tunnelID uint64,
	tunnelDestinationAddr string,
) (*types.Tunnel, error) {
	addr, err := HexToAddress(tunnelDestinationAddr)
	if err != nil {
		return nil, err
	}

	isActive, err := cp.queryTargetContractIsActive(tunnelID, addr)
	if err != nil {
		return nil, err
	}

	return &types.Tunnel{
		ID:            tunnelID,
		TargetAddress: tunnelDestinationAddr,
		IsActive:      isActive,
	}, nil
}

func (cp *EVMChainProvider) queryTargetContractIsActive(
	tunnelID uint64,
	addr gethcommon.Address,
) (bool, error) {
	calldata, err := cp.TunnelRouterABI.Pack("isActive", tunnelID, addr)
	if err != nil {
		return false, err
	}

	b, err := cp.Client.Query(context.Background(), cp.TunnelRouterAddress, calldata)
	if err != nil {
		return false, err
	}

	var output TunnelRouterIsActiveOutput
	if err := cp.TunnelRouterABI.UnpackIntoInterface(&output, "isActive", b); err != nil {
		return false, err
	}

	return output.IsActive, nil
}
