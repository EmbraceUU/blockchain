package logic

import (
	"encoding/json"
	"math/big"
	"testing"
	"web-blockchain/internal/model/input"
)

func TestEVMLogic_TransferLogs(t *testing.T) {
	data, err := es.TransferLogs(18031348)
	if err != nil {
		t.Error(err)
		return
	}

	for _, tx := range data {
		str, _ := json.Marshal(tx)
		t.Logf("%s", string(str))
	}
}

func TestEVMLogic_SendTransaction(t *testing.T) {
	param := input.TransferInput{
		PrivateKey:   "xxxx",
		ToAddress:    "xxxxx",
		TokenAddress: "0x326C977E6efc84E512bB9C30f76E30c160eD06FB", // LINK
		Amount:       big.NewFloat(0.45),
	}

	hash, err := es.SendTransaction(param)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(hash)
}
