package logic

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/big"
	"strings"
	"time"
	"web-blockchain/internal/model/bo"
	"web-blockchain/internal/model/request"
	"web-blockchain/internal/service"
)

type BTCLogic struct {
	btcService service.BTCService

	log *logrus.Logger
}

func NewBTCLogic(btcService service.BTCService, log *logrus.Logger) *BTCLogic {
	logic := new(BTCLogic)
	logic.btcService = btcService
	logic.log = log
	return logic
}

// GetBlockTransactions 查询指定块的交易记录，并返回inputs、outputs、fee
// 查询 block hash
// 查询 block
// 遍历 block txs
// 1. 组装 outputs
// 2. 组装 inputs
// 3. 计算 fee
// 返回
func (l *BTCLogic) GetBlockTransactions(blockNumber int64) ([]bo.BTCTransaction, error) {
	blockHash := new(request.RespBlockHash)
	err := l.btcService.GetBlockHash(blockNumber, blockHash)
	if err != nil {
		return nil, fmt.Errorf("get block hash failed, %s. ", err.Error())
	}

	blockInfo := new(request.RespBlockInfo)
	err = l.btcService.GetBlock(blockHash.Result, 2, blockInfo)
	if err != nil {
		return nil, fmt.Errorf("get block failed, %s. ", err.Error())
	}

	return l.blockInfoHandler(&blockInfo.Result)
}

func (l *BTCLogic) blockInfoHandler(blockInfo *request.BlockInfo) ([]bo.BTCTransaction, error) {
	result := make([]bo.BTCTransaction, 0, blockInfo.NTx)

	l.log.Infof("%d tx in block, begin handle. ", len(blockInfo.Tx))

	for idx, tx := range blockInfo.Tx {
		btx, err := l.transactionHandler(&tx)
		if err != nil {
			l.log.Infof("transaction handler failed, %s. ", err.Error())
			continue
		}

		if btx.CoinBase {
			l.log.Infof("coinbase filter %s. ", tx.TxId)
			continue
		}

		str, _ := json.Marshal(btx)

		l.log.Infof("print %d tx. ", idx+1)
		l.log.Info(string(str))
		
		result = append(result, btx)

		// 需要查询 input, 控制一下频率
		time.Sleep(time.Millisecond * 100)
	}
	return result, nil
}

func (l *BTCLogic) transactionHandler(tx *request.TransactionInfo) (bo.BTCTransaction, error) {
	var btx bo.BTCTransaction
	outAmount := big.NewFloat(0)
	for _, out := range tx.VOut {
		if strings.Contains(out.ScriptPubKey.ASM, "OP_RETURN") {
			continue
		}

		var bout bo.BTCOutput
		// index
		bout.Index = out.N
		// public key
		bout.Address = out.ScriptPubKey.Address
		// amount in float64
		bout.Amount = out.Value

		btx.Outputs = append(btx.Outputs, bout)
		// out amount sum
		outAmount = new(big.Float).Add(outAmount, big.NewFloat(out.Value))
	}

	inAmount := big.NewFloat(0)
	// 没有本地交易记录, 直接查询 rpc node
	for _, in := range tx.VIn {
		if in.Coinbase != "" {
			// filter coinbase input
			btx.CoinBase = true
			return btx, nil
		}

		rawTx := new(request.RespRawTransaction)
		err := l.btcService.GetRawTransaction(in.TxId, true, rawTx)
		if err != nil {
			// 出现错误直接返回，避免计算错误
			return btx, fmt.Errorf("get input detail failed, txId: %s, %s. ", in.TxId, err.Error())
		}

		// 检查 input vout 数量
		if len(rawTx.Result.VOut) < in.VOut+1 {
			return btx, fmt.Errorf("get input detail failed, txId: %s, input vout number invalid. ", in.TxId)
		}

		vout := rawTx.Result.VOut[in.VOut]

		var bin bo.BTCInput
		// index
		bin.Index = vout.N
		// address
		bin.Address = vout.ScriptPubKey.Address
		// amount
		bin.Amount = vout.Value

		btx.Inputs = append(btx.Inputs, bin)

		// input amount sum
		inAmount = new(big.Float).Add(inAmount, big.NewFloat(vout.Value))
	}

	// 检查 fee
	if fee := new(big.Float).Sub(inAmount, outAmount); fee.Cmp(big.NewFloat(0)) < 0 {
		return btx, fmt.Errorf("invalid fee, input < output, txId: %s. ", tx.TxId)
	} else {
		btx.Fee.Amount, _ = fee.Float64()
		return btx, nil
	}
}
