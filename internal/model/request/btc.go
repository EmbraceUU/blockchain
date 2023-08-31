package request

type Param struct {
	ID     string        `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

type Base struct {
	Error string      `json:"error"`
	ID    interface{} `json:"id"`
}

type RespBlockInfo struct {
	Base
	Result BlockInfo `json:"result"`
}

type RespBlockHash struct {
	Base
	Result string `json:"result"`
}

type RespRawTransaction struct {
	Base
	Result TransactionInfo `json:"result"`
}

type BlockInfo struct {
	Hash          string            `json:"hash"`
	Confirmations uint              `json:"confirmations"`
	Size          uint              `json:"size"`
	Weight        uint              `json:"weight"`
	Version       uint              `json:"version"`
	Tx            []TransactionInfo `json:"tx"`
	Time          int64             `json:"time"`
	MedianTime    int64             `json:"mediantime"`
	Nonce         uint              `json:"nonce"`
	NTx           uint              `json:"nTx"`
}

type TransactionInfo struct {
	InActiveChain bool   `json:"in_active_chain"`
	TxId          string `json:"txid"`
	Hash          string `json:"hash"`
	Version       uint   `json:"version"`
	BlockHash     string `json:"blockhash"`
	Confirmations uint   `json:"confirmations"`
	Time          int64  `json:"time"`
	BlockTime     int64  `json:"blocktime"`
	Size          uint   `json:"size"`
	Weight        int    `json:"weight"`
	VIn           []In   `json:"vin"`
	VOut          []Out  `json:"vout"`
}

type In struct {
	TxId     string `json:"txid"`
	VOut     int    `json:"vout"`
	Coinbase string `json:"coinbase"`
}

type Out struct {
	Value        float64 `json:"value"`
	N            int     `json:"n"`
	ScriptPubKey struct {
		ASM     string `json:"asm"` // 判断是否是 OP_RETURN
		Address string `json:"address"`
	} `json:"scriptPubKey"`
}
