package main

import (
	"net/http"

	"github.com/AlexDubtsov/bsa-2024/helpers"

	"github.com/gin-gonic/gin"
)

var wallet = helpers.Wallet{
	UnspentTransactions: []helpers.Transaction{},
}

func main() {
	// Create DB if not exist
	helpers.DBinit()

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
	c.JSON(http.StatusOK, gin.H{"transactions": wallet.UnspentTransactions})
}

// getBalance = handler to get current balance in BTC and EUR
func getBalance(c *gin.Context) {
	var balBTC float64
	var balEur float64

	// ...

	c.JSON(http.StatusOK, gin.H{"balance_BTC": balBTC, "balance_EUR": balEur})
}

// postTransfer = handler for creating a new transaction
func postTransfer(c *gin.Context) {
	// ...

	c.JSON(http.StatusOK, gin.H{"message": "Transfer created successfully"})
}
