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
	data    map[string]interface{} // Temporary storage for data changes
	deleted map[string]bool        // Track deleted keys
	parent  *Transaction
}

type InMemoryDatabase struct {
	currentTransaction *Transaction
	logger             *log.Logger
	data               map[string]interface{}
}

func NewInMemoryDatabase() *InMemoryDatabase {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	return &InMemoryDatabase{
		currentTransaction: nil,
		logger:             logger,
		data:               make(map[string]interface{})}
}

// StartTransaction is starting a new transaction.
func (db *InMemoryDatabase) StartTransaction() {
	db.currentTransaction = &Transaction{
		data:    make(map[string]interface{}),
		deleted: make(map[string]bool),
		parent:  db.currentTransaction, // will be nil is db does not have currentTransaction
	}
	db.logger.Println(TransactionStartLog)
}

// Rollback all changes made within the current transaction and discard them.
func (db *InMemoryDatabase) Rollback() {
	if db.currentTransaction == nil {
		db.logger.Println(TransactionNotStartedLog)
		return
	}
	db.currentTransaction = db.currentTransaction.parent
	db.logger.Println(TransactionRollbackLog)
}

// Commit all changes made within the current transaction to the database.
func (db *InMemoryDatabase) Commit() {
	if db.currentTransaction == nil {
		db.logger.Println(TransactionNotStartedLog)
		return
	}
	if db.currentTransaction.parent != nil {
		// copy data to parent transaction
		for k, v := range db.currentTransaction.data {
			db.currentTransaction.parent.data[k] = v
		}
		// Apply delete operations
		for k := range db.currentTransaction.deleted {
			db.currentTransaction.parent.deleted[k] = true
		}
	} else {
		// if this is a main transaction - copy data to db storage
		for k, v := range db.currentTransaction.data {
			db.data[k] = v
		}
		// Apply delete operations
		for k := range db.currentTransaction.deleted {
			delete(db.data, k)
		}
	}
	db.currentTransaction = db.currentTransaction.parent
	db.logger.Println(TransactionCommitLog)
}

// Set is storing a key-value pair in the database.
func (db *InMemoryDatabase) Set(key string, value interface{}) {
	if db.currentTransaction != nil {
		db.currentTransaction.data[key] = value
		// Remove key from deleted map if it was deleted before
		delete(db.currentTransaction.deleted, key)
	} else {
		db.data[key] = value
	}
	db.logger.Printf(SetKeyValueLog, key, value)
}

// Get the value associated with the given key. Returns the value associated with the key or nil if the key does not exist.
func (db *InMemoryDatabase) Get(key string) interface{} {
	current := db.currentTransaction
	if current != nil && current.deleted[key] {
		return nil
	}
	for current != nil {
		if val, ok := current.data[key]; ok {
			db.logger.Printf(GetValueLog, key)
			return val
		}
		if current.deleted[key] {
			return nil
		}
		current = current.parent
	}
	if val, ok := db.data[key]; ok {
		db.logger.Printf(GetValueLog, key)
		return val
	}
	db.logger.Printf(ValueNotFoundLog, key)
	return nil
}

// Delete the key-value pair associated with the given key.
func (db *InMemoryDatabase) Delete(key string) {
	if db.currentTransaction != nil {
		db.currentTransaction.deleted[key] = true
		db.logger.Printf(KeyDeletedLog, key)
	} else {
		delete(db.data, key)
	}
}
