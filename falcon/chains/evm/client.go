package evm

import (
	"math/big"

	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/falcon/chains"
)

var _ chains.Client = &Client{}

// Client is the struct that handles interactions with the EVM chain.
type Client struct {
	Log      *zap.Logger
	Config   *Config
	KeyStore *KeyStore
	Keys     []Key
}

// NewClient creates a new EVM client from config file and load keys.
func NewClient(log *zap.Logger, cfgPath string) *Client {
	// TODO: implement this
	_ = cfgPath

	return &Client{
		Log: log,
	}
}

func (c *Client) GetNonce(address string) (uint64, error) {
	return 0, nil
}

func (c *Client) BroadcastTx(rawTx string) (string, error) {
	return "", nil
}

func (c *Client) GetBalances(accounts []string) ([]*big.Int, error) {
	return nil, nil
}

func (c *Client) GetTunnelNonce(targetAddress string, tunnelID uint64) (uint64, error) {
	return 0, nil
}
