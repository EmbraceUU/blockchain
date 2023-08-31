package bo

type ERC20Transfer struct {
	BlockNumber  uint64 `json:"block_number"`
	TxHash       string `json:"tx_hash"`
	Index        uint   `json:"index"`
	TokenAddress string `json:"token_address"`
	From         string `json:"from"`
	To           string `json:"to"`
	Amount       string `json:"amount"`
}
