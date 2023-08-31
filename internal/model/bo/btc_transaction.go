package bo

type BTCTransaction struct {
	CoinBase bool        `json:"coinbase"`
	Inputs   []BTCInput  `json:"inputs"`
	Outputs  []BTCOutput `json:"outputs"`
	Fee      Fee         `json:"fee"`
}

type BTCOutput struct {
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
	Index   int     `json:"index"`
}

type BTCInput struct {
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
	Index   int     `json:"index"`
}

type Fee struct {
	Amount float64 `json:"amount"`
}
