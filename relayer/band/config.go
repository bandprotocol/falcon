package band

import "time"

// Config defines the configuration for the BandChain client.
type Config struct {
	RpcEndpoints []string      `toml:"rpc_endpoints"`
	Timeout      time.Duration `toml:"timeout"`
}
