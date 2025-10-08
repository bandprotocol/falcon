package alert

const (
	ConnectMutipleClientErrorMsg  = "Failed to connect chain client on all endpoints"
	ConnectSingleClientErrorMsg   = "Failed to connect chain client"
	EstimateGasFeeErrorMsg        = "Failed to estimate gas fee"
	RelayTxErrorMsg               = "Failed to relay transaction"
	GetTunnelErrorMsg             = "Failed to get tunnel from BandChain"
	GetTunnelPacketErrorMsg       = "Failed to get tunnel packet from BandChain"
	GetContractTunnelInfoErrorMsg = "Failed to get tunnel info from contract"
	PacketSigningStatusErrorMsg   = "Failed tunnel packet signing status"
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
