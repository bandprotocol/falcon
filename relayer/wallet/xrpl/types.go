package xrpl

import (
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
)

type SignerPayload struct {
	SignalPrices    []bandtypes.SignalPrice
	Account         string
	OracleId        uint64
	Fee             string
	Sequence        uint64
	LastUpdatedTime uint64
}

func NewSignerPayload(signalPrices []bandtypes.SignalPrice, account string, oracleId uint64, fee string, sequence uint64, lastUpdatedTime uint64) SignerPayload {
	return SignerPayload{
		SignalPrices:    signalPrices,
		Account:         account,
		OracleId:        oracleId,
		Fee:             fee,
		Sequence:        sequence,
		LastUpdatedTime: lastUpdatedTime,
	}
}
