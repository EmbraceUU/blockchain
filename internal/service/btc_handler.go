package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"web-blockchain/internal/model/request"
)

const (
	GetBlockHashMethod      = "getblockhash"
	GetBlockMethod          = "getblock"
	GetRawTransactionMethod = "getrawtransaction"
	ContentType             = "application/json"
)

type btcService struct {
	id       string
	endpoint string
}

func (s *btcService) GetRawTransaction(txId string, verbose bool, result *request.RespRawTransaction) error {
	params := make([]interface{}, 0)
	params = append(params, txId, verbose)
	err := s.request(GetRawTransactionMethod, params, result)
	if err != nil {
		return err
	}

	if result.Error != "" {
		return fmt.Errorf(result.Error)
	}
	return nil
}

func (s *btcService) GetBlock(blockHash string, verbosity int, result *request.RespBlockInfo) error {
	params := make([]interface{}, 0)
	params = append(params, blockHash, verbosity)
	err := s.request(GetBlockMethod, params, result)
	if err != nil {
		return err
	}

	if result.Error != "" {
		return fmt.Errorf(result.Error)
	}
	return nil
}

func (s *btcService) GetBlockHash(blockNumber int64, result *request.RespBlockHash) error {
	params := make([]interface{}, 0)
	params = append(params, blockNumber)
	err := s.request(GetBlockHashMethod, params, result)
	if err != nil {
		return err
	}

	if result.Error != "" {
		return fmt.Errorf(result.Error)
	}
	return nil
}

func (s *btcService) request(method string, param []interface{}, result interface{}) (err error) {
	var body request.Param
	body.ID = s.id
	body.Method = method
	body.Params = append(body.Params, param...)
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		err = fmt.Errorf("param marshal failed, %s. ", err.Error())
		return
	}

	response, err := http.Post(s.endpoint, ContentType, bytes.NewReader(bodyBytes))
	if err != nil {
		err = fmt.Errorf("post failed, %s. ", err.Error())
		return
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("post failed, status is %d. ", response.StatusCode)
		return
	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = fmt.Errorf("read response body failed, %s. ", err.Error())
		return
	}

	_ = response.Body.Close()
	err = json.Unmarshal(respBody, result)
	if err != nil {
		err = fmt.Errorf("unmarshal result failed, %s. ", err.Error())
	}
	return
}
