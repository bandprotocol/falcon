package evm

import _ "embed"

//go:embed abi/TunnelRouter.json
var gasPriceTunnelRouterABI string

// TunnelRouterIsActiveOutput defines the output parameters for the TunnelRouter.isActive method.
type TunnelRouterIsActiveOutput struct {
	IsActive bool `json:"is_active"`
}
