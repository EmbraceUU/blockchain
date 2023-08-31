package service

import (
	"os"
	"testing"
	"web-blockchain/internal/common/log"
	"web-blockchain/internal/config"
	"web-blockchain/internal/consts"
	"web-blockchain/internal/model/request"
)

var s BTCService
var es EVMService

func TestMain(m *testing.M) {
	conf := config.NewConfig()
	serviceLogger := log.WithLoggerName(consts.ServiceLoggerName)
	s = NewBtcService(conf.BTC)
	es = NewEvmService(conf.EVM, serviceLogger)
	os.Exit(m.Run())
}

func TestBtcService_GetBlock(t *testing.T) {
	result := new(request.RespBlockInfo)
	err := s.GetBlock("000000000000000000025441256481aa1d90ff3d8a6dd87f41ea53555fe21585", 2, result)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%v", result)
	}
}

func TestBtcService_GetRawTransaction(t *testing.T) {
	result := new(request.RespRawTransaction)
	err := s.GetRawTransaction("8143b3b341f665b22adcb8489158356c03f7c93cf4e4fa673d8518fa0fed95e4", true, result)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%v", result)
	}
}

func TestBtcService_GetBlockHash(t *testing.T) {
	result := new(request.RespBlockHash)
	err := s.GetBlockHash(805310, result)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%v", result)
	}
}
