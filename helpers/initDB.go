package helpers

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Initialize the database connection
func DBopen() {
	var err error
	DB, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
}

func DBcreate() {

	// Initialize the database connection
	DBopen()

	// Create sql table if it does not exist
	usersTable, err := DB.Prepare(`
    CREATE TABLE if not exists TRANSACTIONS(
        ID TEXT PRIMARY KEY,
        AMOUNT REAL,
        SPENT INTEGER,
        CREATED TEXT
    )
	`)
	if err != nil {
		log.Fatal(err)
	}
	usersTable.Exec()

	// Defer the closing of the database connection
	defer DB.Close()
}
