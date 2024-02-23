package helpers

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func DBinit() {
	var err error

	// Initialize the database connection
	DB, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
}

func DBcreate() {
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
}
