package alert

const (
	ConnectMutipleClientError  = "Failed to connect chain client on all endpoints"
	ConnectSingleClientError   = "Failed to connect chain client"
	EstimateGasFeeError        = "Failed to estimate gas fee"
	RelayTxError               = "Failed to relay transaction"
	GetTunnelError             = "Failed to get tunnel from BandChain"
	GetTunnelPacketError       = "Failed to get tunnel packet from BandChain"
	GetContractTunnelInfoError = "Failed to get tunnel info from contract"
	PacketSigningStatusError   = "Failed tunnel packet signing status"
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
