package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"math/big"
	"net/http"
	"strconv"
	"web-blockchain/api/response"
	"web-blockchain/api/svc"
	input2 "web-blockchain/internal/model/input"
)

type Handler struct {
	svc *svc.ServiceContext

	log *logrus.Logger
	out *logrus.Logger
}

func NewHandler(svc *svc.ServiceContext, log *logrus.Logger, out *logrus.Logger) *Handler {
	var handler Handler
	handler.svc = svc
	handler.log = log
	handler.out = out
	return &handler
}

func (h *Handler) GetBTCTransactions(c *gin.Context) {
	blockNumber := c.Query("blockNumber")
	if blockNumber == "" {
		c.JSON(http.StatusOK, response.Result{
			Code: 1000, // 错误码
			Msg:  "block number is required. ",
		})
		return
	}

	blockNumberInt, err := strconv.Atoi(blockNumber)
	if err != nil {
		c.JSON(http.StatusOK, response.Result{
			Code: 1001, // 错误码
			Msg:  "block number is invalid. ",
		})
		return
	}

	go func(blockNumberInt int) {
		data, err := h.svc.BtcLogic.GetBlockTransactions(int64(blockNumberInt))
		if err != nil {
			h.log.Errorf("GetBlockTransactions process failed, %s. ", err.Error())
			return
		}

		for _, tx := range data {
			txBytes, _ := json.Marshal(tx)
			h.out.Infof(string(txBytes))
		}
	}(blockNumberInt)

	c.JSON(http.StatusOK, response.Result{
		Code: 0,
		Data: "success, data is in log. ",
	})
}

func (h *Handler) GetERC20Transactions(c *gin.Context) {
	blockNumber := c.Query("blockNumber")
	if blockNumber == "" {
		c.JSON(http.StatusOK, response.Result{
			Code: 1000, // 错误码
			Msg:  "block number is required. ",
		})
		return
	}

	blockNumberInt, err := strconv.Atoi(blockNumber)
	if err != nil {
		c.JSON(http.StatusOK, response.Result{
			Code: 1001, // 错误码
			Msg:  "block number is invalid. ",
		})
		return
	}

	go func(blockNumberInt int) {
		data, err := h.svc.EvmLogic.TransferLogs(int64(blockNumberInt))
		if err != nil {
			h.log.Errorf("ERC20TransferLogs process failed, %s. ", err.Error())
			return
		}

		for _, tx := range data {
			txBytes, _ := json.Marshal(tx)
			h.out.Infof(string(txBytes))
		}
	}(blockNumberInt)

	c.JSON(http.StatusOK, response.Result{
		Code: 0,
		Data: "success, data is in log. ",
	})
}

func (h *Handler) SendTransactions(c *gin.Context) {
	var fm response.ReqSendTransaction
	if err := c.ShouldBind(&fm); err != nil {
		c.JSON(http.StatusOK, response.Result{
			Code: 1000, // 错误码
			Msg:  err.Error(),
		})
		return
	}

	amount, _ := new(big.Float).SetString(fm.Amount)
	input := input2.TransferInput{
		PrivateKey:   fm.PrivateKey,
		ToAddress:    fm.ToAddress,
		TokenAddress: fm.TokenAddress,
		Amount:       amount,
	}
	data, err := h.svc.EvmLogic.SendTransaction(input)
	if err != nil {
		h.log.Errorf("ERC20TransferLogs process failed, %s. ", err.Error())
		c.JSON(http.StatusOK, response.Result{
			Code: 1002,
			Data: err.Error(),
		})
		return
	}

	h.log.Infof("hash is %s. ", data)

	c.JSON(http.StatusOK, response.Result{
		Code: 0,
		Data: data,
	})
}
