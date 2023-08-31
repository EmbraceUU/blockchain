# blockchain

## 项目介绍

该 demo 提供了三个功能：
1. btc 交易记录打印
2. evm erc20 转账记录打印
3. evm erc20 转账并广播

## 功能入口

提供了两种方式检测：
1. 构建项目后，从 main 入口启动 http server，可以通过 REST 触发 task，默认端口 8080
   1. btc 交易记录打印：      GET /api/btc/transactions
   2. evm erc20 转账记录打印：GET /api/evm/transactions
   3. evm erc20 转账并广播：  POST /api/evm/transactions
2. 【推荐】/internal/logic/ 目录下的两个 test 文件中，有三个功能的测试用例，同样可以触发 task
   1. btc 交易记录打印：TestBTCLogic_GetBlockTransactions
   2. evm erc20 转账记录打印：TestEVMLogic_TransferLogs
   3. evm erc20 转账并广播：TestEVMLogic_SendTransaction

## 备注

1. 检测之前，需要在 /api/config.json 中配置 btc 和 evm rpc node
2. 【btc 交易记录打印功能】建议通过 debug 检测，因为本地没有 btc 历史数据，查询 inputs 时需要调用 rpc，一次完全的调用可能会消耗不少 request credit，可以只查看部分 tx 数据
3. 如果使用 REST 的方式触发 task，btc 和 evm 的交易记录打印功能是**异步处理**，可以在 log 中查看最终数据