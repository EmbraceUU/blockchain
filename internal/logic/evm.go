package logic

import (
	"web-blockchain/internal/model/bo"
	"web-blockchain/internal/model/input"
	"web-blockchain/internal/service"
)

type EVMLogic struct {
	evmService service.EVMService
}

func NewEVMLogic(evmService service.EVMService) *EVMLogic {
	logic := new(EVMLogic)
	logic.evmService = evmService
	return logic
}

func (l *EVMLogic) TransferLogs(blockNumber int64) (transfers []bo.ERC20Transfer, err error) {
	return l.evmService.ERC20TransferLogs(blockNumber)
}

func (l *EVMLogic) SendTransaction(input input.TransferInput) (txHash string, err error) {
	return l.evmService.SendERC20Transaction(input)
}
