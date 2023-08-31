package svc

import (
	"github.com/sirupsen/logrus"
	"web-blockchain/internal/config"
	"web-blockchain/internal/logic"
	"web-blockchain/internal/service"
)

type ServiceContext struct {
	BtcLogic *logic.BTCLogic
	EvmLogic *logic.EVMLogic
}

func NewServiceContext(conf config.Config, log *logrus.Logger) *ServiceContext {
	btcService := service.NewBtcService(conf.BTC)
	evmService := service.NewEvmService(conf.EVM, log)

	btcLogic := logic.NewBTCLogic(btcService, log)
	evmLogic := logic.NewEVMLogic(evmService)

	svc := new(ServiceContext)
	svc.BtcLogic = btcLogic
	svc.EvmLogic = evmLogic
	return svc
}
