package database

import (
	"log"
	"os"
)

const (
	TransactionStartLog      = "Transaction started"
	TransactionCommitLog     = "Transaction committed"
	TransactionRollbackLog   = "Transaction rolled back"
	SetKeyValueLog           = "Key-value pair set: %s=%s"
	GetValueLog              = "Value retrieved for key: %s"
	ValueNotFoundLog         = "Value not found for key: %s"
	KeyDeletedLog            = "Key deleted: %s"
	TransactionNotStartedLog = "Transaction has not been started"
)

type Transaction struct {
	data   map[string]string
	parent *Transaction
}

type InMemoryDatabase struct {
	currentTransaction *Transaction
	logger             *log.Logger
	data               map[string]string
}

func NewInMemoryDatabase() *InMemoryDatabase {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	return &InMemoryDatabase{
		currentTransaction: nil,
		logger:             logger,
		data:               make(map[string]string)}
}

// StartTransaction is starting a new transaction.
func (db *InMemoryDatabase) StartTransaction() {
	db.currentTransaction = &Transaction{
		data:   make(map[string]string),
		parent: db.currentTransaction, // will be nil is db does not have currentTransaction
	}
	db.logger.Println(TransactionStartLog)
}

// Rollback all changes made within the current transaction and discard them.
func (db *InMemoryDatabase) Rollback() {
	if db.currentTransaction != nil {
		db.currentTransaction = db.currentTransaction.parent
		db.logger.Println(TransactionRollbackLog)
	}
	db.logger.Println(TransactionNotStartedLog)
}

// Commit all changes made within the current transaction to the database.
func (db *InMemoryDatabase) Commit() {
	if db.currentTransaction != nil {
		if db.currentTransaction.parent != nil {
			// copy data to parent transaction
			for k, v := range db.currentTransaction.data {
				db.currentTransaction.parent.data[k] = v
			}
		} else {
			// if this is a main transaction - copy data to db storage
			for k, v := range db.currentTransaction.data {
				db.data[k] = v
			}
		}
		db.currentTransaction = db.currentTransaction.parent
		db.logger.Println(TransactionCommitLog)
	}
	db.logger.Println(TransactionNotStartedLog)
}

// Set is storing a key-value pair in the database.
func (db *InMemoryDatabase) Set(key, value string) {
	if db.currentTransaction != nil {
		db.currentTransaction.data[key] = value
	} else {
		db.data[key] = value
	}
	db.logger.Printf(SetKeyValueLog, key, value)
}

// Get the value associated with the given key. Returns the value associated with the key or "" if the key does not exist.
func (db *InMemoryDatabase) Get(key string) string {
	current := db.currentTransaction
	for current != nil {
		if val, ok := current.data[key]; ok {
			db.logger.Printf(GetValueLog, key)
			return val
		}
		current = current.parent
	}
	if val, ok := db.data[key]; ok {
		db.logger.Printf(GetValueLog, key)
		return val
	}
	db.logger.Printf(ValueNotFoundLog, key)
	return ""
}

// Delete the key-value pair associated with the given key.
func (db *InMemoryDatabase) Delete(key string) {
	if db.currentTransaction != nil {
		delete(db.currentTransaction.data, key)
		db.logger.Printf(KeyDeletedLog, key)
	}
	// fix that. Value should be deleted in main db and only after commit
}
