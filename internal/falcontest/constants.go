package falcontest

const DefaultCfgText = `[global]
checking_packet_interval = 60000000000

[bandchain]
rpc_endpoints = ['http://localhost:26657']
timeout = 5

[target_chains]
`

const CustomCfgText = `[global]
checking_packet_interval = 1

[bandchain]
rpc_endpoints = ['http://localhost:26657', 'http://localhost:26658']
timeout = 0

[target_chains]
`
