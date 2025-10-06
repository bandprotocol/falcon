package alert

import (
	"fmt"
)

const (
	ConnectClientError         = "Failed to connect client"
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
func HandleAlert(alert Alert, topic string, tunnelID uint64, chainName string, detail string) {
	if alert == nil {
		return
	}
	alert.Trigger(buildTopic(topic, tunnelID, chainName), detail)
}

// HandleReset handles resetting alert with the given topic, including tunnel ID and chain name.
func HandleReset(alert Alert, topic string, tunnelID uint64, chainName string) {
	if alert == nil {
		return
	}
	alert.Reset(buildTopic(topic, tunnelID, chainName))
}

// buildTopic append the topic string with tunnel ID and chain name.
func buildTopic(topic string, tunnelID uint64, chainName string) string {
	return fmt.Sprintf("%s TUNNEL_ID-%d CHAIN-%s", topic, tunnelID, chainName)
}
