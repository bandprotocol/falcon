package evm

import (
	_ "embed"
	"math/big"
)

//go:embed abi/TunnelRouter.json
var gasPriceTunnelRouterABI string

// TunnelInfoOutput defines the output parameters for the TunnelRouter.isActive method.
type TunnelInfoOutput struct {
	IsActive       bool     `json:"is_active"`
	LatestSequence uint64   `json:"latest_sequence"`
	Balance        *big.Int `json:"balance"`
}

type TunnelInfoOutputRaw struct {
	Info TunnelInfoOutput `json:"tunnel_info"`
}
