package main

import (
	"fmt"

	"github.com/anastasiia-skliar/inmemory-key-value-db/database"
)

func main() {
	db := database.NewInMemoryDatabase()

	db.StartTransaction()
	db.Set("a", "1")
	fmt.Println(db.Get("a")) // Output: 1

	db.StartTransaction()
	db.Set("a", "2")
	fmt.Println(db.Get("a")) // Output: 2

	db.Rollback()
	fmt.Println(db.Get("a")) // Output: 1

	db.Commit()
	fmt.Println(db.Get("a")) // Output: 1
}
