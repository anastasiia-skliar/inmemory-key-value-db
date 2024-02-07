package database

import (
	"log"
	"os"
)

const (
	TransactionStartLog    = "Transaction started"
	TransactionCommitLog   = "Transaction committed"
	TransactionRollbackLog = "Transaction rolled back"
	SetKeyValueLog         = "Key-value pair set: %s=%s"
	GetValueLog            = "Value retrieved for key: %s"
	ValueNotFoundLog       = "Value not found for key: %s"
	KeyDeletedLog          = "Key deleted: %s"
)

type Transaction struct {
	data       map[string]string
	parent     *Transaction
	commited   bool
	rollbacked bool
}

type InMemoryDatabase struct {
	currentTransaction *Transaction
	logger             *log.Logger
}

func NewInMemoryDatabase() *InMemoryDatabase {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	return &InMemoryDatabase{currentTransaction: nil, logger: logger}
}

// StartTransaction is starting a new transaction. All operations within this transaction are isolated from others.
func (db *InMemoryDatabase) StartTransaction() {
	db.currentTransaction = &Transaction{data: make(map[string]string), parent: db.currentTransaction}
	db.logger.Println(TransactionStartLog)
}

// Rollback all changes made within the current transaction and discard them.
func (db *InMemoryDatabase) Rollback() {
	if db.currentTransaction != nil && !db.currentTransaction.commited && !db.currentTransaction.rollbacked {
		db.currentTransaction.rollbacked = true
		db.currentTransaction = db.currentTransaction.parent
		db.logger.Println(TransactionRollbackLog)
	}
}

// Commit all changes made within the current transaction to the database.
func (db *InMemoryDatabase) Commit() {
	if db.currentTransaction != nil && !db.currentTransaction.commited && !db.currentTransaction.rollbacked {
		if db.currentTransaction.parent != nil {
			for k, v := range db.currentTransaction.data {
				db.currentTransaction.parent.data[k] = v
			}
		}
		db.currentTransaction.commited = true
		db.currentTransaction = db.currentTransaction.parent
		db.logger.Println(TransactionCommitLog)
	}
}

// Set is storing a key-value pair in the database.
func (db *InMemoryDatabase) Set(key, value string) {
	if db.currentTransaction != nil && !db.currentTransaction.commited && !db.currentTransaction.rollbacked {
		db.currentTransaction.data[key] = value
		db.logger.Printf(SetKeyValueLog, key, value)
	}
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
	db.logger.Printf(ValueNotFoundLog, key)
	return ""
}

// Delete the key-value pair associated with the given key.
func (db *InMemoryDatabase) Delete(key string) {
	if db.currentTransaction != nil && !db.currentTransaction.commited && !db.currentTransaction.rollbacked {
		delete(db.currentTransaction.data, key)
		db.logger.Printf(KeyDeletedLog, key)
	}
}
