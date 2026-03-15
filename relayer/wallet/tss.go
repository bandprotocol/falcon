package wallet

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	tssPacketABI, _ = abi.NewType("tuple", "result", []abi.ArgumentMarshaling{
		{Name: "Sequence", Type: "uint64"},
		{
			Name:         "RelayPrices",
			Type:         "tuple[]",
			InternalType: "struct Prices[]",
			Components: []abi.ArgumentMarshaling{
				{Name: "SignalID", Type: "bytes32"},
				{Name: "Price", Type: "uint64"},
			},
		},
		{Name: "CreatedAt", Type: "int64"},
	})

	tssPacketArgs = abi.Arguments{
		{Type: tssPacketABI, Name: "packet"},
	}
)

const EncoderABIPrefixLength = 4

type TssPayload struct {
	TssMessage []byte
	RandomAddr []byte
	Signature  []byte
}

func NewTssPayload(tssMessage, randomAddr, signature []byte) TssPayload {
	return TssPayload{
		TssMessage: tssMessage,
		RandomAddr: randomAddr,
		Signature:  signature,
	}
}

// TSSPacket represents the Packet that will be used for encoding a tss message.
type TSSPacket struct {
	Sequence    uint64
	RelayPrices []RelayPrice
	CreatedAt   int64
}

// RelayPrice represents the price data for relaying to other chains.
type RelayPrice struct {
	SignalID [32]byte
	Price    uint64
}

func (p TssPayload) ToTSSPacket() (*TSSPacket, error) {
	if len(p.TssMessage) < EncoderABIPrefixLength {
		return nil, fmt.Errorf("tss message should have at least %d bytes", EncoderABIPrefixLength)
	}
	payload := p.TssMessage[EncoderABIPrefixLength:]

	values, err := tssPacketArgs.Unpack(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack ABI data: %w", err)
	}

	var result struct {
		TSSPacket TSSPacket
	}

	if err := tssPacketArgs.Copy(&result, values); err != nil {
		return nil, fmt.Errorf("failed to copy ABI data: %w", err)
	}

	return &result.TSSPacket, nil
}

func Bytes32ToString(byteArray [32]byte) string {
	trimmed := bytes.TrimLeft(byteArray[:], "\x00")
	return string(trimmed)
}
