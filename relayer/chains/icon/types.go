package icon

import (
	iconclient "github.com/icon-project/goloop/client"
)

// ClientConnectionResult is the struct that contains the result of connecting to the specific endpoint.
type ClientConnectionResult struct {
	Endpoint    string
	Client      *iconclient.ClientV3
	BlockHeight uint64
}
