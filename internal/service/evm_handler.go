package service

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
	"math"
	"math/big"
	"strings"
	"sync"
	"time"
	"web-blockchain/internal/common/abi/erc20"
	"web-blockchain/internal/model/bo"
	"web-blockchain/internal/model/input"
)

const (
	ERC20TransferSign   = "Transfer(address,address,uint256)"
	ERC20TransferFnSign = "transfer(address,uint256)"
)

type evmService struct {
	endpoint string

	cache  map[string]*big.Int
	cacheM sync.RWMutex

	log *logrus.Logger
}

func (s *evmService) getLocalDecimal(tokenAddress string) *big.Int {
	s.cacheM.RLock()
	defer s.cacheM.RUnlock()

	return s.cache[strings.ToLower(tokenAddress)]
}

func (s *evmService) setLocalDecimal(tokenAddress string, decimal *big.Int) {
	s.cacheM.Lock()
	s.cache[strings.ToLower(tokenAddress)] = decimal
	s.cacheM.Unlock()
}

func (s *evmService) ERC20TransferLogs(blockNumber int64) (transfers []bo.ERC20Transfer, err error) {
	// check param
	if !s.validBlockNumber(blockNumber) {
		err = fmt.Errorf("block number is invalid. ")
		return
	}

	// connect to rpc node
	client, err := ethclient.Dial(s.rpcNode())
	if err != nil {
		err = fmt.Errorf("dial client failed, %s. ", err.Error())
		return
	}

	// query param
	blockNumberBig := big.NewInt(blockNumber)
	q := ethereum.FilterQuery{
		FromBlock: blockNumberBig,
		ToBlock:   blockNumberBig,
		Topics: [][]common.Hash{
			{crypto.Keccak256Hash([]byte(ERC20TransferSign))},
		},
	}

	// filter logs
	logs, err := client.FilterLogs(context.Background(), q)
	if err != nil {
		err = fmt.Errorf("filter logs failed, %s. ", err.Error())
		return
	}

	// contract abi
	contractAbi, err := abi.JSON(strings.NewReader(erc20.Erc20MetaData.ABI))
	if err != nil {
		err = fmt.Errorf("abi reader failed, %s. ", err.Error())
		return
	}

	s.log.Infof("begin parse log. ")

	for _, vLog := range logs {
		var transferEvent bo.LogTransfer

		// 使用 abi 解析 data 数据
		err = contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
		if err != nil {
			// ERC721 Transfer param is equal to ERC20, but diff in log, ERC721 is 4 but ERC20 is 3.
			s.log.Errorf("This is not ERC20 transfer event. %d-%s-%d,  %s \n",
				vLog.BlockNumber, vLog.TxHash.Hex(), vLog.Index, err.Error())
			continue
		}

		transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
		transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

		decimal, err := s.erc20Decimal(vLog.Address, client)
		if err != nil {
			s.log.Errorf("get decimal failed, %s, %s. ", vLog.Address.Hex(), err.Error())
			continue
		}

		decimalBig := big.NewFloat(math.Pow(10, -float64(decimal)))
		valueBig, _ := new(big.Float).SetString(transferEvent.Value.String())
		amount := new(big.Float).Mul(decimalBig, valueBig)

		transfers = append(transfers, bo.ERC20Transfer{
			BlockNumber:  vLog.BlockNumber,
			TxHash:       vLog.TxHash.Hex(),
			Index:        vLog.Index,
			TokenAddress: vLog.Address.Hex(),
			From:         transferEvent.From.Hex(),
			To:           transferEvent.To.Hex(),
			Amount:       amount.Text('f', 6),
		})
	}
	return
}

func (s *evmService) transferKeys(key string) (*ecdsa.PrivateKey, common.Address, error) {
	var (
		fromAddress common.Address
		err         error
	)

	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		err = fmt.Errorf("hex to private key failed, %s", err.Error())
		return privateKey, fromAddress, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		err = fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return privateKey, fromAddress, err
	}

	// from address
	fromAddress = crypto.PubkeyToAddress(*publicKeyECDSA)
	return privateKey, fromAddress, nil
}

func (s *evmService) erc20Decimal(tokenAddress common.Address, client *ethclient.Client) (uint8, error) {
	if cacheDecimal := s.getLocalDecimal(tokenAddress.String()); cacheDecimal != nil {
		return uint8(cacheDecimal.Int64()), nil
	}

	// get decimal
	erc20Contract, err := erc20.NewErc20(tokenAddress, client)
	if err != nil {
		err = fmt.Errorf("new erc20 contract failed, %s. ", err.Error())
		return 0, err
	}
	// 生成合约会话对象
	erc20ContractSession := erc20.Erc20Session{
		Contract: erc20Contract,
	}
	decimal, err := erc20ContractSession.Decimals()
	if err != nil {
		err = fmt.Errorf("get token decimal failed, %s. ", err.Error())
		return 0, err
	}

	time.Sleep(time.Millisecond * 100)

	s.setLocalDecimal(tokenAddress.String(), big.NewInt(int64(decimal)))
	return decimal, nil
}

func (s *evmService) assembleData(toAddress, tokenAddress common.Address, inputAmount *big.Float, client *ethclient.Client) ([]byte, error) {
	// method id
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(ERC20TransferFnSign))
	methodID := hash.Sum(nil)[:4]

	// padded address
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)

	// get token decimal
	decimal, err := s.erc20Decimal(tokenAddress, client)
	if err != nil {
		return nil, err
	}

	// padded amount
	amount, _ := new(big.Float).Mul(big.NewFloat(math.Pow(10, float64(decimal))), inputAmount).Int(&big.Int{})
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)
	return data, nil
}

func (s *evmService) createUnSignTx(fromAddress common.Address, chainID *big.Int, input input.TransferInput, client *ethclient.Client) (*types.Transaction, error) {
	// get account nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		err = fmt.Errorf("pending nonce at failed, %s. ", err.Error())
		return nil, err
	}

	// suggest gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		err = fmt.Errorf("suggest gas price failed, %s. ", err.Error())
		return nil, err
	}

	// suggest gas tip cap
	tipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		err = fmt.Errorf("suggest gas tip cap failed, %s. ", err.Error())
		return nil, err
	}

	toAddress := common.HexToAddress(input.ToAddress)
	tokenAddress := common.HexToAddress(input.TokenAddress)

	// assemble transaction data
	data, err := s.assembleData(toAddress, tokenAddress, input.Amount, client)
	if err != nil {
		return nil, err
	}

	// gas limit
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:      fromAddress,
		To:        &tokenAddress,
		GasPrice:  gasPrice,
		GasTipCap: tipCap,
		Value:     big.NewInt(0),
		Data:      data,
	})
	if err != nil {
		err = fmt.Errorf("estimate gas failed, %d. ", gasLimit)
		return nil, err
	}

	// Create a new transaction
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: gasPrice,
		Gas:       gasLimit,
		To:        &tokenAddress,
		Value:     big.NewInt(0),
		Data:      data,
	})
	return tx, nil
}

func (s *evmService) SendERC20Transaction(input input.TransferInput) (txHash string, err error) {
	// connect to rpc node
	client, err := ethclient.Dial(s.rpcNode())
	if err != nil {
		err = fmt.Errorf("dial client failed, %s. ", err.Error())
		return
	}

	// get fromAddress by private key
	privateKey, fromAddress, err := s.transferKeys(input.PrivateKey)
	if err != nil {
		return
	}

	// get network id
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		err = fmt.Errorf("network id failed, %s. ", err.Error())
		return
	}

	// create un sign tx
	tx, err := s.createUnSignTx(fromAddress, chainID, input, client)
	if err != nil {
		return
	}

	// sign tx
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		err = fmt.Errorf("sign tx failed, %s. ", err.Error())
		return
	}

	// send transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		err = fmt.Errorf("send transaction failed, %s. ", err.Error())
		return
	}

	return signedTx.Hash().Hex(), nil
}

func (s *evmService) validBlockNumber(blockNumber int64) bool {
	if blockNumber < 0 {
		return false
	}
	return true
}

func (s *evmService) rpcNode() string {
	return s.endpoint
}
