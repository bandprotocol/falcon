package alert

const (
	ConnectSingleBandClientErrorMsg  = "Failed to connect BandChain client"
	ConnectAllBandClientErrorMsg     = "Failed to connect BandChain client on all endpoints"
	ConnectSingleChainClientErrorMsg = "Failed to connect chain client"
	ConnectAllChainClientErrorMsg    = "Failed to connect chain client on all endpoints"
	EstimateGasFeeErrorMsg           = "Failed to estimate gas fee"
	RelayTxErrorMsg                  = "Failed to relay transaction"
	GetTunnelErrorMsg                = "Failed to get tunnel from BandChain"
	GetTunnelPacketErrorMsg          = "Failed to get tunnel packet from BandChain"
	GetContractTunnelInfoErrorMsg    = "Failed to get tunnel info from contract"
	PacketSigningStatusErrorMsg      = "Failed tunnel packet signing status"
	GetBlockErrorMsg                 = "Failed to get block from chain"
	GetBalanceErrorMsg               = "Failed to prepare database transaction"
	SaveDatabaseErrorMsg             = "Failed to save to database"
)

// Alert represents an object that triggers and resets alerts.
type Alert interface {
	Trigger(topic, detail string)
	Reset(topic string)
}

// HandleAlert handles triggering alert with the given topic and detail, including tunnel ID and chain name.
func HandleAlert(alert Alert, topic *Topic, detail string) {
	if alert == nil {
		return
	}
	alert.Trigger(topic.GetFullTopic(), detail)
}

// HandleReset handles resetting alert with the given topic, including tunnel ID and chain name.
func HandleReset(alert Alert, topic *Topic) {
	if alert == nil {
		return
	}
	alert.Reset(topic.GetFullTopic())
}
