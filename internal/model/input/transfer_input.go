package input

import "math/big"

type TransferInput struct {
	PrivateKey   string
	ToAddress    string
	TokenAddress string
	Amount       *big.Float
}
