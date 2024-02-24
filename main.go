package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AlexDubtsov/bsa-2024/helpers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create DB (in case of absence)
	if helpers.DB == nil {
		helpers.DBcreate()
	}

	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// Endpoint for fetching all transactions
	router.GET("/transactions", getTransactions)
	// Endpoint for fetchinc current balance
	router.GET("/balance", getBalance)
	// Endpoint for making new deposit
	router.POST("/supply", postSupply)
	// Endpoint for making new transaction
	router.POST("/transfer", postTransfer)
	// PrintOut
	fmt.Println("\t✔ SUCCESS")
	fmt.Println("The server operates on http://localhost:8080/")
	fmt.Println("\n\tUSAGE:\n✖ Check not spent funds amount:\t\thttp://localhost:8080/balance\n✖ List of all transactions:\t\thttp://localhost:8080/transactions")
	fmt.Println("\n✖ Endpoint for the transfer request:\thttp://localhost:8080/transfer\nJSON payload example:\n\t\t{\n\t\t\t\"requested_amount\": \"0.001\"\n\t\t}")
	fmt.Println("\n✖ Endpoint for the supply request:\thttp://localhost:8080/supply\nJSON payload example:\n\t\t{\n\t\t\t\"supplied_amount\": \"0.001\"\n\t\t}")
	fmt.Println("REQUEST LOG:")
	// Server run
	router.Run(":8080")
}

// getTransactions = handler to get a list of all transactions
func getTransactions(c *gin.Context) {

	var wallet = helpers.Wallet{
		AllTransactions: []helpers.Transaction{},
	}

	// Initialize the database connection
	helpers.DBopen()

	// Query the database to get all transactions
	records, err := helpers.DB.Query("SELECT ID, AMOUNT, SPENT, CREATED FROM TRANSACTIONS")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error getTransactions1": err.Error()})
		return
	}
	defer records.Close()

	// Defer the closing of the database connection
	defer helpers.DB.Close()

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

	c.JSON(http.StatusOK, gin.H{"transactions": wallet.AllTransactions})
}

// getBalance = handler to get current balance in BTC and EUR
func getBalance(c *gin.Context) {
	currencyPair := "BTC/EUR"
	urlRateAPI := "http://api-cryptopia.adca.sh/v1/prices/ticker"

	balBTC := helpers.CalcBalance()

	balEur := balBTC * helpers.GetRate(urlRateAPI, currencyPair)

	c.JSON(http.StatusOK, gin.H{"balance_BTC": balBTC, "balance_EUR": balEur})
}

// postTransfer = handler for creating a new transaction
func postSupply(c *gin.Context) {
	// Define minimum supply amount
	var minSupply float64 = 0.00001
	var err error

	// Struct init
	var supplyRequest helpers.CreateSupplyRequest
	// Trying to read JSON-request
	err = c.BindJSON(&supplyRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error postSupply1": err.Error()})
		return
	}

	// Convert Requested Amount value string to float64
	SuppliedAmount, err := strconv.ParseFloat(supplyRequest.SuppliedAmount, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error postSupply2": err.Error()})
		return
	}

	// Check if supplied amount is less then permitted
	if SuppliedAmount < minSupply {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Supplied amount is less then permitted"})
		return
	}

	// Write changes to DB; supply flag = true means this transaction is supply
	err = helpers.DBreadAndWrite(SuppliedAmount, true)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error postSupply3": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Supply accepted successfully"})
}

// postTransfer = handler for creating a new transaction
func postTransfer(c *gin.Context) {
	// Define minimum request amount
	var minRequest float64 = 0.00001
	var err error

	// Struct init
	var transferRequest helpers.CreateTransferRequest

	// Trying to read JSON-request
	err = c.BindJSON(&transferRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error postTransfer1": err.Error()})
		return
	}

	// Convert Requested Amount value string to float64
	RequestedAmount, err := strconv.ParseFloat(transferRequest.RequestedAmount, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error postTransfer2": err.Error()})
		return
	}

	// Check if requested amount is less then permitted
	if RequestedAmount < minRequest {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requested amount is less then permitted"})
		return
	}

	// Check if enough not spent amount available
	if RequestedAmount > helpers.CalcBalance() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough balance for the transfer"})
		return
	}

	// Write changes to DB; supply flag = false means this transaction is request
	err = helpers.DBreadAndWrite(RequestedAmount, false)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error postTransfer3": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer created successfully"})
}
