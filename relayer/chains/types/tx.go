package types

import (
	"database/sql/driver"
	"fmt"
)

// TxStatus is the status of the transaction
type TxStatus int

const (
	TX_STATUS_PENDING TxStatus = iota + 1
	TX_STATUS_SUCCESS
	TX_STATUS_FAILED
	TX_STATUS_TIMEOUT
)

var txStatusNameMap = map[TxStatus]string{
	TX_STATUS_PENDING: "Pending",
	TX_STATUS_SUCCESS: "Success",
	TX_STATUS_FAILED:  "Failed",
	TX_STATUS_TIMEOUT: "Timeout",
}

func (t TxStatus) String() string {
	return txStatusNameMap[t]
}

var txStatusFromString = map[string]TxStatus{
	"Pending": TX_STATUS_PENDING,
	"Success": TX_STATUS_SUCCESS,
	"Failed":  TX_STATUS_FAILED,
	"Timeout": TX_STATUS_TIMEOUT,
}

// need to manually create `tx_status` type in a database first
// by "CREATE TYPE tx_status AS ENUM ('Pending', 'Success', 'Failed', 'Timeout')"
func (t *TxStatus) Scan(value interface{}) error {
	tx, ok := txStatusFromString[value.(string)]
	if !ok {
		return fmt.Errorf("invalid tx status")
	}
	*t = tx
	return nil
}

func (t TxStatus) Value() (driver.Value, error) { return t.String(), nil }
