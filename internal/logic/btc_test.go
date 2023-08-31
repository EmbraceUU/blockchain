package logic

import (
	"encoding/json"
	"os"
	"testing"
	"web-blockchain/internal/common/log"
	"web-blockchain/internal/config"
	"web-blockchain/internal/consts"
	"web-blockchain/internal/service"
)

var s *BTCLogic
var es *EVMLogic

func TestMain(m *testing.M) {
	conf := config.NewConfig()
	serviceLogger := log.WithLoggerName(consts.ServiceLoggerName)

	btcService := service.NewBtcService(conf.BTC)
	evmService := service.NewEvmService(conf.EVM, serviceLogger)
	s = NewBTCLogic(btcService, serviceLogger)
	es = NewEVMLogic(evmService)
	os.Exit(m.Run())
}

func TestBTCLogic_GetBlockTransactions(t *testing.T) {
	data, err := s.GetBlockTransactions(805310)
	if err != nil {
		t.Error(err)
		return
	}

	for _, tx := range data {
		str, _ := json.Marshal(tx)
		t.Logf("%s", string(str))
	}
}
