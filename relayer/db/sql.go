package db

import (
	"fmt"
	"strings"

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
// Supported drivers: "postgresql", "sqlite".
func NewSQL(dbPath string) (SQL, error) {
	var db *gorm.DB
	var err error

	cfg := &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	}

	driverName, path, err := splitDbPath(dbPath)
	if err != nil {
		return SQL{}, err
	}

	fmt.Println(driverName, path)

	switch driverName {
	case "postgres":
		db, err = gorm.Open(postgres.Open(dbPath), cfg)
		if err != nil {
			return SQL{}, err
		}
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(path), cfg)
		if err != nil {
			return SQL{}, err
		}

	default:
		return SQL{}, gorm.ErrUnsupportedDriver
	}

	return SQL{db: db}, nil
}

// splitDbPath splits "<driver>:<dsn>" into driver and DSN.
// Keeps colons inside DSN (uses SplitN). For SQLite
// Example: "postgresql:postgres://u:p@host:5432/db" -> ("postgresql", "postgres://u:p@host:5432/db")
// Example: "sqlite:///myfile.db" -> myfile.db
func splitDbPath(dbPath string) (string, string, error) {
	parts := strings.SplitN(dbPath, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid db path")
	}

	driver := parts[0]
	path := parts[1]
	if driver == "sqlite" {
		path = strings.TrimPrefix(path, "///")
	}

	return driver, path, nil
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

// AddOrUpdateSenderBalance inserts a new sender record or updates the existing balance.
// If the sender already exists, its balance and updated_at are updated.
func (sql SQL) AddOrUpdateSenderBalance(sender *Sender) error {
	up := sql.db.Model(&Sender{}).
		Where(&Sender{Address: sender.Address}).
		Updates(&Sender{Balance: sender.Balance})
	if up.RowsAffected == 0 {
		return sql.db.Create(&sender).Error
	}
	return up.Error
}
