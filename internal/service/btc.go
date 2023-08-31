package service

import (
	"web-blockchain/internal/config"
	"web-blockchain/internal/model/request"
)

type BTCService interface {
	// GetRawTransaction 根据tx查询 raw transaction
	GetRawTransaction(txId string, verbose bool, result *request.RespRawTransaction) error

	// GetBlockHash 查询 block hash
	GetBlockHash(blockNumber int64, result *request.RespBlockHash) error

	// GetBlock 查询 block info, 里面可以包含 tx 信息
	GetBlock(blockHash string, verbosity int, result *request.RespBlockInfo) error
}

func NewBtcService(conf config.BTC) BTCService {
	service := new(btcService)
	service.endpoint = conf.EndPoint
	service.id = conf.ID
	return service
}
