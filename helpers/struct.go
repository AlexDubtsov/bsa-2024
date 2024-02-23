package helpers

type Transaction struct {
	ID        string
	Amount    float64
	Spent     bool
	CreatedAt string
}

type Wallet struct {
	UnspentTransactions []Transaction
}
