package server

import (
	"github.com/gin-gonic/gin"
	"web-blockchain/api/controller"
)

func InitRouter(handler *controller.Handler) *gin.Engine {
	r := gin.New()
	api := r.Group("/api")
	{
		api.GET("/btc/transactions", handler.GetBTCTransactions)
		api.GET("/evm/transactions", handler.GetERC20Transactions)
		api.POST("/evm/transactions", handler.SendTransactions)
	}

	return r
}
