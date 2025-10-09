package db

// Database defines the interface for the db interaction with the target chain.
type Database interface {
	AddOrUpdateTransaction(transaction *Transaction) error
	AddOrUpdateSenderBalance(sender *Sender) error
}
