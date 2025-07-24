package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
)

var _ Database = &SQL{}

type SQL struct {
	db *gorm.DB
}

func NewSQL(driverName, dbPath string) (SQL, error) {
	var db *gorm.DB
	var err error

	switch driverName {
	case "postgresql":
		db, err = gorm.Open(postgres.Open(dbPath), &gorm.Config{
			SkipDefaultTransaction: true,
			Logger: logger.Default.
				LogMode(logger.Silent),
		})
		if err != nil {
			return SQL{}, err
		}
	case "mysql":
		db, err = gorm.Open(mysql.Open(dbPath), &gorm.Config{
			SkipDefaultTransaction: true,
			Logger: logger.Default.
				LogMode(logger.Silent),
		})
		if err != nil {
			return SQL{}, err
		}

	default:
		return SQL{}, err
	}

	if err = db.AutoMigrate(&Transaction{}, &SignalPrice{}); err != nil {
		return SQL{}, err
	}

	return SQL{db: db}, nil
}

func (sql SQL) AddOrUpdateTransaction(transaction *Transaction) error {
	var queryTransaction Transaction

	err := sql.db.Where(&Transaction{TxHash: transaction.TxHash}).First(&queryTransaction).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Insert transaction if not found
			return sql.db.Create(transaction).Error
		}
		return err
	}

	// If the existing transaction is pending and the new status is not pending, update it
	if queryTransaction.Status == chainstypes.TX_STATUS_PENDING &&
		transaction.Status != chainstypes.TX_STATUS_PENDING {
		if err := sql.db.
			Where(&Transaction{TxHash: transaction.TxHash}).
			Updates(&Transaction{
				Status:            transaction.Status,
				GasUsed:           transaction.GasUsed,
				EffectiveGasPrice: transaction.EffectiveGasPrice,
				BalanceDelta:      transaction.BalanceDelta,
				Timestamp:         transaction.Timestamp,
			}).Error; err != nil {
			return err
		}
	}

	return nil
}
