package evm

import (
	"math/big"

	"github.com/bandprotocol/falcon/falcon/chains"
	"github.com/bandprotocol/falcon/falcon/keys"
)

var _ chains.Client = &Client{}

type Client struct {
	config *Config
	keys   []keys.Key
}

func NewClient(cfg *Config, relayers []keys.Key) *Client {
	return &Client{
		config: cfg,
		keys:   relayers,
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
