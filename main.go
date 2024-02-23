package main

import (
	"net/http"
	"strconv"

	"github.com/AlexDubtsov/bsa-2024/helpers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create DB if not exist
	helpers.DBinit()
	helpers.DBcreate()
	// Scheduled close connection to DB
	defer helpers.DB.Close()

	router := gin.Default()
	// Endpoint for fetching all transactions
	router.GET("/transactions", getTransactions)
	// Endpoint for fetchinc current balance
	router.GET("/balance", getBalance)
	// Endpoint for making new transaction
	router.POST("/transfer", postTransfer)
	// Server run
	router.Run(":8080")
}

// getTransactions = handler to get a list of all transactions
func getTransactions(c *gin.Context) {

	var wallet = helpers.Wallet{
		AllTransactions: []helpers.Transaction{},
	}

	// Open DB connection
	helpers.DBinit()
	// Scheduled close connection to DB
	defer helpers.DB.Close()

	// Query the database to get all transactions
	records, err := helpers.DB.Query("SELECT ID, AMOUNT, SPENT, CREATED FROM TRANSACTIONS")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error getTransactions1": err.Error()})
		return
	}
	defer records.Close()

	// Iterate over the records; Add each record to Wallet
	for records.Next() {
		var transaction helpers.Transaction
		err := records.Scan(&transaction.ID, &transaction.Amount, &transaction.Spent, &transaction.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error getTransactions2": err.Error()})
			return
		}
		wallet.AllTransactions = append(wallet.AllTransactions, transaction)
	}

	// Check for errors during iteration
	if err := records.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error getTransactions3": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": wallet.AllTransactions})
}

// getBalance = handler to get current balance in BTC and EUR
func getBalance(c *gin.Context) {
	currencyPair := "BTC/EUR"
	var currencyRateStr string
	var currencyRateFloat float64
	var balBTC float64
	var balEur float64

	// Open external API
	urlRate := "http://api-cryptopia.adca.sh/v1/prices/ticker"
	var fetchedStruct struct {
		Data []helpers.RateAPI `json:"data"`
	}
	err := helpers.GetJson(urlRate, &fetchedStruct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error getBalance1": err.Error()})
		return
	}

	// Currency pair rate choice
	for _, word := range fetchedStruct.Data {
		if word.Symbol == currencyPair {
			currencyRateStr = word.Value
		}
	}

	// Convert string to float64
	currencyRateFloat, err = strconv.ParseFloat(currencyRateStr, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error getBalance2": err.Error()})
		return
	}

	// Open DB connection
	helpers.DBinit()
	// Scheduled close connection to DB
	defer helpers.DB.Close()
	// Query the database to get not spent transactions
	records, err := helpers.DB.Query("SELECT ID, AMOUNT, SPENT, CREATED FROM TRANSACTIONS where SPENT = ?", 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error getBalance1": err.Error()})
		return
	}
	defer records.Close()

	// Iterate over the records; Add each record to Wallet
	for records.Next() {
		var transaction helpers.Transaction
		err := records.Scan(&transaction.ID, &transaction.Amount, &transaction.Spent, &transaction.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error getTransactions2": err.Error()})
			return
		}
		balBTC += transaction.Amount
	}

	balEur = balBTC * currencyRateFloat

	c.JSON(http.StatusOK, gin.H{"balance_BTC": balBTC, "balance_EUR": balEur})
}

// postTransfer = handler for creating a new transaction
func postTransfer(c *gin.Context) {
	// ...

	c.JSON(http.StatusOK, gin.H{"message": "Transfer created successfully"})
}
