package service

import (
	"github.com/sirupsen/logrus"
	"math/big"
	"web-blockchain/internal/config"
	bo2 "web-blockchain/internal/model/bo"
	"web-blockchain/internal/model/input"
)

type EVMService interface {
	// SendERC20Transaction 转移ERC20, 并广播事件
	SendERC20Transaction(input input.TransferInput) (txHash string, err error)

	// ERC20TransferLogs 查询指定块的log
	ERC20TransferLogs(blockNumber int64) (transfers []bo2.ERC20Transfer, err error)
}

func NewEvmService(conf config.EVM, log *logrus.Logger) EVMService {
	service := new(evmService)
	service.endpoint = conf.EndPoint
	service.log = log
	service.cache = make(map[string]*big.Int)
	return service
}
