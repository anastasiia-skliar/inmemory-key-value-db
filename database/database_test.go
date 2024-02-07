package database

import (
	"testing"
)

func TestDatabase(t *testing.T) {
	tests := []struct {
		name       string
		operations func(db *InMemoryDatabase)
		assertions func(db *InMemoryDatabase, t *testing.T)
	}{
		{
			name: "Set and Get",
			operations: func(db *InMemoryDatabase) {
				db.StartTransaction()
				db.Set("a", "1")
			},
			assertions: func(db *InMemoryDatabase, t *testing.T) {
				if val := db.Get("a"); val != "1" {
					t.Errorf("Expected value for key 'a' to be '1', got '%s'", val)
				}
			},
		},
		{
			name: "Rollback",
			operations: func(db *InMemoryDatabase) {
				db.StartTransaction()
				db.Set("a", "1")
				db.Rollback()
			},
			assertions: func(db *InMemoryDatabase, t *testing.T) {
				if val := db.Get("a"); val != "" {
					t.Errorf("Expected value for key 'a' to be empty after rollback, got '%s'", val)
				}
			},
		},
		{
			name: "Commit",
			operations: func(db *InMemoryDatabase) {
				db.StartTransaction()
				db.Set("a", "1")
				db.Commit()
			},
			assertions: func(db *InMemoryDatabase, t *testing.T) {
				if val := db.Get("a"); val != "1" {
					t.Errorf("Expected value for key 'a' to be '1' after commit, got '%s'", val)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := NewInMemoryDatabase()
			test.operations(db)
			test.assertions(db, t)
		})
	}
}

func TestNestedTransactions(t *testing.T) {
	tests := []struct {
		name       string
		operations func(db *InMemoryDatabase)
		assertions func(db *InMemoryDatabase, t *testing.T)
	}{
		{
			name: "Nested Commit",
			operations: func(db *InMemoryDatabase) {
				db.StartTransaction()
				db.Set("a", "1")
				db.StartTransaction()
				db.Set("a", "2")
				db.Commit()
			},
			assertions: func(db *InMemoryDatabase, t *testing.T) {
				if val := db.Get("a"); val != "2" {
					t.Errorf("Expected value for key 'a' to be '2' after nested commit, got '%s'", val)
				}
			},
		},
		{
			name: "Nested Rollback",
			operations: func(db *InMemoryDatabase) {
				db.StartTransaction()
				db.Set("a", "1")
				db.StartTransaction()
				db.Set("a", "2")
				db.Rollback()
			},
			assertions: func(db *InMemoryDatabase, t *testing.T) {
				if val := db.Get("a"); val != "1" {
					t.Errorf("Expected value for key 'a' to be '1' after nested rollback, got '%s'", val)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := NewInMemoryDatabase()
			test.operations(db)
			test.assertions(db, t)
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name       string
		operations func(db *InMemoryDatabase)
		assertions func(db *InMemoryDatabase, t *testing.T)
	}{
		{
			name: "Delete",
			operations: func(db *InMemoryDatabase) {
				db.StartTransaction()
				db.Set("a", "1")
				db.Delete("a")
			},
			assertions: func(db *InMemoryDatabase, t *testing.T) {
				if val := db.Get("a"); val != "" {
					t.Errorf("Expected value for key 'a' to be empty after deletion, got '%s'", val)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := NewInMemoryDatabase()
			test.operations(db)
			test.assertions(db, t)
		})
	}
}

func TestNonexistentKey(t *testing.T) {
	db := NewInMemoryDatabase()

	if val := db.Get("nonexistent"); val != "" {
		t.Errorf("Expected value for nonexistent key to be empty, got '%s'", val)
	}
}

func TestGetAfterCommit(t *testing.T) {
	db := NewInMemoryDatabase()

	db.StartTransaction()
	db.Set("a", "1")
	db.Commit()

	if val := db.Get("a"); val != "1" {
		t.Errorf("Expected value for key 'a' to be '1' after commit, got '%s'", val)
	}
}

func TestGetAfterRollback(t *testing.T) {
	db := NewInMemoryDatabase()

	db.StartTransaction()
	db.Set("a", "1")
	db.Rollback()

	if val := db.Get("a"); val != "" {
		t.Errorf("Expected value for key 'a' to be empty after rollback, got '%s'", val)
	}
}
