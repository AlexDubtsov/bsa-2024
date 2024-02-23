package helpers

type Transaction struct {
	ID        string  `json:"id"`
	Amount    float64 `json:"amount"`
	Spent     bool    `json:"spent"`
	CreatedAt string  `json:"created"`
}

type Wallet struct {
	AllTransactions []Transaction `json:"transactions"`
}

type RateAPI struct {
	Symbol     string `json:"symbol"`
	Value      string `json:"value"`
	Sources    int    `json:"sources"`
	Updated_at string `json:"updated_at"`
}
