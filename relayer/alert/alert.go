package alert

import (
	"fmt"

	"github.com/bandprotocol/falcon/relayer/logger"
)

const (
	ConnectClientError         = "Failed to connect client"
	EstimateGasFeeError        = "Failed to estimate gas fee"
	CreateAndSignTxError       = "Failed to create and sign transaction"
	BroadcastTxError           = "Failed to broadcast transaction"
	BumpGasError               = "Failed to bump gas"
	ConfirmSuccessTxError      = "Failed to confirm success transaction"
	GetTunnelError             = "Failed to get tunnel from BandChain"
	GetTunnelPacketError       = "Failed to get tunnel packet from BandChain"
	GetContractTunnelInfoError = "Failed to get tunnel info from contract"
	PacketSigningStatusError   = "Failed tunnel packet signing status"
)

// Alert represents an object that can receive notifications when a new alert is fired or resolved.
type Alert interface {
	Trigger(topic, detail string) error
	Resolve(topic string) error
}

// HandleAlert sends an alert with the given topic and detail, including tunnel ID and chain name.
func HandleAlert(alert Alert, topic, detail string, tunnelID uint64, chainName string, log logger.Logger) {
	if alert == nil {
		return
	}
	if err := alert.Trigger(buildTopic(topic, tunnelID, chainName), detail); err != nil {
		log.Debug("Failed to send alert", err)
	}
}

// HandleResolve resolves an alert with the given topic, including tunnel ID and chain name.
func HandleResolve(alert Alert, topic string, tunnelID uint64, chainName string, log logger.Logger) {
	if alert == nil {
		return
	}
	if err := alert.Resolve(buildTopic(topic, tunnelID, chainName)); err != nil {
		log.Debug("Failed to resolve alert", err)
	}
}

// buildTopic append the topic string with tunnel ID and chain name.
func buildTopic(topic string, tunnelID uint64, chainName string) string {
	return fmt.Sprintf("%s TUNNEL_ID-%d CHAIN-%s", topic, tunnelID, chainName)
}
