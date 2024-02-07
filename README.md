# In-Memory Database in Go

This project implements a simple in-memory database in Go with support for transactions and key-value storage. In addition, the database supports nested transactions, which means you are able to create a transaction within a transaction. If parent transaction is commited, all uncommited child transaction data is lost.

## Features

- Supports transactions
- Key-value storage
- Operations: get, set, delete, start_transaction, commit, roll_back
- Nested transactions
- Not supposed to be used for concurrent requests

## Getting Started

### Prerequisites

- Go 1.21 or higher installed

### Installation

```bash
git clone https://github.com/anastasiia-skliar/inmemory-key-value-db.git
cd inmemory-key-value-db
go run main.go
```
### Usage
Import the package in your Go code:
```go
import "github.com/yourusername/in-memory-database"
```
Create a new instance of the database:
```go
db := inmemory.NewInMemoryDatabase()
```
Use the provided methods to interact with the database:
```go
db.startTransaction()
db.set("key", "value")
value := db.get("key")
fmt.Println(value) // Output: value
db.commit()
```
### Testing
To run tests, use the following command:

```bash
make test
```
### Test coverage
To check test coverage, use the following command:

```bash
make coverage
```
### Linting
To lint the code, use the following command:

```bash
make lint
```







