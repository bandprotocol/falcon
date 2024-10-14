package chains

// ChainType defines the type of the target chain.
type ChainType int

const (
	ChainTypeUndefined ChainType = iota
	ChainTypeEVM
	Cosmwasm
)

// Config defines the common configuration for the target chain client.
type Config struct{}
