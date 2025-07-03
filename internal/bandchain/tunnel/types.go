package tunnel

// events
const (
	EventTypeUpdateParams             = "update_params"
	EventTypeCreateTunnel             = "create_tunnel"
	EventTypeUpdateSignalsAndInterval = "update_signals_and_interval"
	EventTypeActivateTunnel           = "activate_tunnel"
	EventTypeDeactivateTunnel         = "deactivate_tunnel"
	EventTypeTriggerTunnel            = "trigger_tunnel"
	EventTypeProducePacketFail        = "produce_packet_fail"
	EventTypeProducePacketSuccess     = "produce_packet_success"
	EventTypeDepositToTunnel          = "deposit_to_tunnel"
	EventTypeWithdrawFromTunnel       = "withdraw_from_tunnel"

	AttributeKeyParams           = "params"
	AttributeKeyTunnelID         = "tunnel_id"
	AttributeKeySequence         = "sequence"
	AttributeKeyInterval         = "interval"
	AttributeKeyRoute            = "route"
	AttributeKeyEncoder          = "encoder"
	AttributeKeyInitialDeposit   = "initial_deposit"
	AttributeKeyFeePayer         = "fee_payer"
	AttributeKeySignalID         = "signal_id"
	AttributeKeySoftDeviationBPS = "soft_deviation_bps"
	AttributeKeyHardDeviationBPS = "hard_deviation_bps"
	AttributeKeyIsActive         = "is_active"
	AttributeKeyCreatedAt        = "created_at"
	AttributeKeyCreator          = "creator"
	AttributeKeyDepositor        = "depositor"
	AttributeKeyWithdrawer       = "withdrawer"
	AttributeKeyAmount           = "amount"
	AttributeKeyReason           = "reason"
)
