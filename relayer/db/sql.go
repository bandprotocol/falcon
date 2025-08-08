package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"

	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
)

var _ Database = &SQL{}

// SQL is instance that wraps Gorm DB instance.
type SQL struct {
	db *gorm.DB
}

// NewSQL opens a new Gorm connection using the given driverName and dbPath.
// Supported drivers: "postgresql", "mysql", "sqlite", "sqlite3".
func NewSQL(driverName, dbPath string) (SQL, error) {
	var db *gorm.DB
	var err error

	cfg := &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	}

	switch driverName {
	case "postgresql":
		db, err = gorm.Open(postgres.Open(dbPath), cfg)
		if err != nil {
			return SQL{}, err
		}
	case "mysql":
		db, err = gorm.Open(mysql.Open(dbPath), cfg)
		if err != nil {
			return SQL{}, err
		}
	case "sqlite", "sqlite3":
		db, err = gorm.Open(sqlite.Open(dbPath), cfg)
		if err != nil {
			return SQL{}, err
		}

	default:
		return SQL{}, gorm.ErrUnsupportedDriver
	}

	if err = db.AutoMigrate(&Transaction{}, &SignalPrice{}); err != nil {
		return SQL{}, err
	}

	return SQL{db: db}, nil
}

// AddOrUpdateTransaction inserts a new Transaction record if none exists with the same TxHash.
// If an existing record is in PENDING state and the new transaction has progressed to a non-PENDING status.
func (sql SQL) AddOrUpdateTransaction(transaction *Transaction) error {
	return sql.db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "tx_hash"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"status", "gas_used", "effective_gas_price", "balance_delta", "block_timestamp", "updated_at",
			}),
			Where: clause.Where{
				Exprs: []clause.Expression{
					clause.Expr{
						SQL:  "transactions.status = ? AND EXCLUDED.status <> ?",
						Vars: []interface{}{chainstypes.TX_STATUS_PENDING, chainstypes.TX_STATUS_PENDING},
					},
				},
			},
		}).
		Create(transaction).
		Error
}
