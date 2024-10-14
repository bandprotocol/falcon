package band

// Config defines the configuration for the BandChain client.
type Config struct {
	RpcEndpoints []string `toml:"rpc_endpoints"`
	Timeout      int      `toml:"timeout"`
}
