package band

import "time"

// Config defines the configuration for the BandChain client.
type Config struct {
	RpcEndpoints               []string      `mapstructure:"rpc_endpoints"                toml:"rpc_endpoints"`
	Timeout                    time.Duration `mapstructure:"timeout"                      toml:"timeout"`
	LivelinessCheckingInterval time.Duration `mapstructure:"liveliness_checking_interval" toml:"liveliness_checking_interval"`
}
