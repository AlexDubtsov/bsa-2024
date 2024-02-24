package helpers

type Transaction struct {
	ID        string  `json:"id"`
	Amount    float64 `json:"amount"`
	Spent     int     `json:"spent"`
	CreatedAt string  `json:"created"`
}

type Wallet struct {
	AllTransactions []Transaction `json:"transactions"`
}

type RateApiSingleRecord struct {
	Symbol     string `json:"symbol"`
	Value      string `json:"value"`
	Sources    int    `json:"sources"`
	Updated_at string `json:"updated_at"`
}

type RateApiAllRecords struct {
	Data []RateApiSingleRecord `json:"data"`
}

type CreateTransferRequest struct {
	RequestedAmount string `json:"requested_amount"`
}

type CreateSupplyRequest struct {
	SuppliedAmount string `json:"supplied_amount"`
}
