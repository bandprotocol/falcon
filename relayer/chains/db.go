package chains

import (
	"github.com/bandprotocol/falcon/relayer/alert"
	"github.com/bandprotocol/falcon/relayer/db"
	"github.com/bandprotocol/falcon/relayer/logger"
)

// HandleSaveTransaction saves the transaction to the database and triggers alert if any error occurs.
func HandleSaveTransaction(db db.Database, a alert.Alert, tx *db.Transaction, log logger.Logger) {
	if tx == nil {
		return
	}
	if db == nil {
		log.Debug("Database is not set; skipping saving transaction")
		return
	}
	if err := db.AddOrUpdateTransaction(tx); err != nil {
		log.Error("Save transaction error", err)
		alert.HandleAlert(a, alert.NewTopic(alert.SaveDatabaseErrorMsg).
			WithTunnelID(tx.TunnelID).
			WithChainName(tx.ChainName), err.Error())
	} else {
		alert.HandleReset(a, alert.NewTopic(alert.SaveDatabaseErrorMsg).
			WithTunnelID(tx.TunnelID).
			WithChainName(tx.ChainName))
	}
}
