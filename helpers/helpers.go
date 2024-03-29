package helpers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"os"
	"time"
)

func CalcBalance() float64 {
	var precision float64 = 0.00001 // Precision in calculating balance for skipping precision errors of float64
	var result float64

	// Initialize the database connection
	DBopen()

	// Query the database to get not spent transactions
	records, err := DB.Query("SELECT ID, AMOUNT, SPENT, CREATED FROM TRANSACTIONS where SPENT = ?", 0)
	if err != nil {
		fmt.Println("Error", err.Error())
		os.Exit(1)
	}

	defer records.Close()

	// Iterate over the records; Add each record to Wallet
	for records.Next() {
		var transaction Transaction
		err := records.Scan(&transaction.ID, &transaction.Amount, &transaction.Spent, &transaction.CreatedAt)
		if err != nil {
			fmt.Println("Error", err.Error())
			os.Exit(1)
		}
		result += math.Round(transaction.Amount/precision) * precision
	}

	// Defer the closing of the database connection
	defer DB.Close()

	return result
}

func DBreadAndWrite(Amount float64, supply bool) error {
	// Generating set of newIDs
	newId, err := randHexStr()
	if err != nil {
		fmt.Println("Error creating random strings:", err)
		return err
	}

	// Starting Amount below 0
	RemainingAmount := Amount * (-1)
	// Continue to add amount?
	addMoreAmount := true

	// Initialize the database connection
	DBopen()

	// Start a transaction
	tx, err := DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	// ### TO FIX THIS: in case of many DB records not to load all records? ###

	// Query the database to get all transactions
	records, err := tx.Query("SELECT ID, AMOUNT, SPENT, CREATED FROM TRANSACTIONS")
	if err != nil {
		fmt.Println("Error", err.Error())
		return err
	}

	// Iterate through the records, check if ID is unique and update the SPENT column (if !supply)
	for records.Next() {
		var record Transaction
		err := records.Scan(&record.ID, &record.Amount, &record.Spent, &record.CreatedAt)
		if err != nil {
			fmt.Println("Error scanning record:", err.Error())
			continue
		}

		// If not to supply, but REQUEST funds
		if !supply {
			if record.Spent == 0 && addMoreAmount { // If current DB record is not spent and all previous have not covered requested amount
				// Update the SPENT column
				_, err = tx.Exec("UPDATE TRANSACTIONS SET SPENT = 1 WHERE ID = ?", record.ID)
				if err != nil {
					fmt.Println("Error updating record:", err.Error())
				}

				// Update the remaining amount
				RemainingAmount += record.Amount
			}
			// Stop on Remaining Amount is positive
			if RemainingAmount >= 0 {
				addMoreAmount = false
			}
		} else { // If SUPPLY
			RemainingAmount = Amount
		}

		// If newID is same as name of record => remove newID from list of available IDs
		for record.ID == newId[0] {
			newId = newId[1:]
		}
	}

	defer records.Close()

	// Timestamp calc
	now := time.Now()
	// Convert time variable to string using Format method
	created := now.Format("2006-01-02 15:04:05")

	// Insert record to DB
	err = DBInsert(tx, newId[0], RemainingAmount, 0, created)
	if err != nil {
		return err
	}

	// Commit the DB transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	// Defer the closing of the database connection
	defer DB.Close()

	return nil
}

func DBInsert(tx *sql.Tx, id string, amount float64, spent int, created string) error {

	// Insert into DB
	_, err := tx.Exec("INSERT INTO TRANSACTIONS (ID, AMOUNT, SPENT, CREATED) VALUES (?, ?, ?, ?)", id, amount, spent, created)
	if err != nil {
		fmt.Println("Error inserting record:", err)
		return err
	}

	return nil
}

func randHexStr() ([]string, error) {
	// Number of bytes used to convert to string
	numBytes := 16
	// Number of IDs that are generated
	numStrings := 5

	var randomHex []string

	for i := 0; i < numStrings; i++ {
		// Generate random bytes
		randomBytes := make([]byte, numBytes)
		_, err := rand.Read(randomBytes)
		if err != nil {
			return make([]string, 0), err
		}

		// Convert bytes to hexadecimal string
		randomHex = append(randomHex, hex.EncodeToString(randomBytes))
	}

	return randomHex, nil
}
