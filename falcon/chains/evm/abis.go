package evm

const gasPriceTunnelRouterABI = `[
  {
    "type": "function",
    "name": "isActive",
    "inputs": [
      { "name": "", "type": "uint64", "internalType": "uint64" },
      { "name": "", "type": "address", "internalType": "address" }
    ],
    "outputs": [{ "name": "is_active", "type": "bool", "internalType": "bool" }],
    "stateMutability": "view"
  }
]`

// TunnelRouterIsActiveOutput defines the output parameters for the TunnelRouter.isActive method.
type TunnelRouterIsActiveOutput struct {
	IsActive bool `json:"is_active"`
}
