package icon

import (
	iconclient "github.com/icon-project/goloop/client"
	"github.com/icon-project/goloop/server/jsonrpc"
	"github.com/shopspring/decimal"
)

// ClientConnectionResult is the struct that contains the result of connecting to the specific endpoint.
type ClientConnectionResult struct {
	Endpoint    string
	Client      *iconclient.ClientV3
	BlockHeight uint64
}

type ContractData struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type GetRefDataBulkParams struct {
	Bases  []string `json:"bases"`
	Quotes []string `json:"quotes"`
}

type ContractOutput struct {
	Rate            jsonrpc.HexInt `json:"rate"`
	LastUpdateBase  jsonrpc.HexInt `json:"last_update_base"`
	LastUpdateQuote jsonrpc.HexInt `json:"last_update_quote"`
}

func (c ContractOutput) Parse() (decimal.Decimal, error) {
	rb, err := c.Rate.BigInt()
	if err != nil {
		return decimal.Decimal{}, err
	}
	return decimal.NewFromBigInt(rb, 0), nil
}
