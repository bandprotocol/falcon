package relayer

import (
	"context"
	"encoding/hex"
	"time"

	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/relayer/band"
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/types"
)

// TunnelRelayer is a relayer that listens to the tunnel and relays the packet
type TunnelRelayer struct {
	Log                    *zap.Logger
	TunnelID               uint64
	ContractAddress        string
	CheckingPacketInterval time.Duration
	BandClient             band.Client
	TargetChainProvider    chains.ChainProvider
}

// NewTunnelRelayer creates a new TunnelRelayer
func NewTunnelRelayer(
	log *zap.Logger,
	tunnelID uint64,
	contractAddress string,
	checkingPacketInterval time.Duration,
	bandClient band.Client,
	targetChainProvider chains.ChainProvider,
) TunnelRelayer {
	return TunnelRelayer{
		Log:                    log,
		TunnelID:               tunnelID,
		ContractAddress:        contractAddress,
		CheckingPacketInterval: checkingPacketInterval,
		BandClient:             bandClient,
		TargetChainProvider:    targetChainProvider,
	}
}

func (t *TunnelRelayer) CheckAndRelay(ctx context.Context) error {
	// TODO: Check if there is any new packet to be relayed in BandChain
	task, err := t.MockRelayerTask()
	if err != nil {
		t.Log.Error(
			"failed to mock relayer task",
			zap.Error(err),
			zap.Uint64("tunnel_id", t.TunnelID),
		)
		return err
	}

	if err := t.TargetChainProvider.RelayPacket(ctx, task); err != nil {
		t.Log.Error(
			"failed to relay packet",
			zap.Error(err),
			zap.Uint64("tunnel_id", t.TunnelID),
			zap.Uint64("sequence", task.Packet.Sequence),
		)
		return err
	}

	return nil
}

// TODO: remove this after the implementation is done
func (t *TunnelRelayer) MockRelayerTask() (*types.RelayerTask, error) {
	msgHex := "0E1AC2C4A50A82AA49717691FC1AE2E5FA68EFF45BD8576B0F2BE7A0850FA7C6" +
		"78512D24E95216DC140F557181A03631715A023424CBAD94601F3546CDFC3DE4" +
		"78512D24E95216DC140F557181A03631715A023424CBAD94601F3546CDFC3DE4" +
		"000000006705E8A00000000000000002D3813E0CCBA0AD5A" +
		"0000000000000000000000000000000000000000000000000000000000000020" +
		"0000000000000000000000000000000000000000000000000000000000000001" +
		"0000000000000000000000000000000000000000000000000000000000000002" +
		"00000000000000000000000000000000000000000000000000000000000000C0" +
		"0000000000000000000000000000000000000000000000000000000000000100" +
		"0000000000000000000000000000000000000000000000000000000000000160" +
		"000000000000000000000000000000000000000000000000000000006705E8A0" +
		"0000000000000000000000000000000000000000000000000000000000000003" +
		"6574680000000000000000000000000000000000000000000000000000000000" +
		"000000000000000000000000000000000000000000000000000000000000002A" +
		"307865303046316638356162444232614636373630373539353437643435306461363843453636426231" +
		"00000000000000000000000000000000000000000000" +
		"0000000000000000000000000000000000000000000000000000000000000002" +
		"0000000000000000000000000063727970746F5F70726963652E627463757364" +
		"0000000000000000000000000000000000000000000000000000000000000000" +
		"0000000000000000000000000063727970746F5F70726963652E657468757364" +
		"0000000000000000000000000000000000000000000000000000000000000000"

	msg, err := hex.DecodeString(msgHex)
	if err != nil {
		return nil, err
	}

	rAddr, err := hex.DecodeString("0b7754FD4545b561C1bc2E978922A5b7772F01D8")
	if err != nil {
		return nil, err
	}

	signature, err := hex.DecodeString("5A1B0A6ACD177D54D88E8CF18706C8ABB98EE3BBBC58A4AAA1351E3EA8AB9FC6")
	if err != nil {
		return nil, err
	}

	return &types.RelayerTask{
		Packet: &bandtypes.Packet{
			TunnelID:              t.TunnelID,
			TargetContractAddress: t.ContractAddress,
			TargetChainID:         "testnet_evm",
		},
		Signing: &bandtypes.Signing{
			Message: msg,
			EVMSignature: &bandtypes.EVMSignature{
				RAddress:  rAddr,
				Signature: signature,
			},
		},
	}, nil
}
