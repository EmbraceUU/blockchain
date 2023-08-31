package response

type ReqSendTransaction struct {
	PrivateKey   string `json:"privateKey" form:"privateKey" binding:"required"`
	ToAddress    string `json:"toAddress" form:"toAddress" binding:"required"`
	TokenAddress string `json:"tokenAddress" form:"tokenAddress" binding:"required"`
	Amount       string `json:"amount" form:"amount" binding:"required"`
}
