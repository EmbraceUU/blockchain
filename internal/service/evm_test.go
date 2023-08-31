package service

import (
	"math/big"
	"testing"
	"web-blockchain/internal/model/input"
)

func TestService_TransferERC20(t *testing.T) {
	param := input.TransferInput{
		PrivateKey:   "xxxxx",
		ToAddress:    "xxxxxx",
		TokenAddress: "0x326C977E6efc84E512bB9C30f76E30c160eD06FB",
		Amount:       big.NewFloat(0.45),
	}

	txHash, err := es.SendERC20Transaction(param)
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log(txHash)
	}
}
