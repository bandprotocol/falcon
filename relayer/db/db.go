package db

type Database interface {
	AddOrUpdateTransaction(transaction *Transaction) error
}
