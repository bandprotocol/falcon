package db

// Database defines the interface for the db interaction with the target chain.
type Database interface {
	AddOrUpdateTransaction(transaction *Transaction) error
	GetLatestSuccessSequence(tunnelID uint64) (uint64, error)
}
