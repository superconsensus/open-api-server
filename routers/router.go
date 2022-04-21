package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/superconsensus/open-api-server/controllers/v1"
	"github.com/superconsensus/open-api-server/middlewares"
)

func NewRouter() *gin.Engine {

	//记录到文件
	//dir := path.Base(conf.Log.FilePath)
	//err := os.MkdirAll(dir, os.ModePerm)
	//if err != nil {
	//	panic(err)
	//}
	//logFilePath := conf.Log.FilePath
	//logFileName := conf.Log.RouterFile
	//logfile := path.Join(logFilePath, logFileName)
	//f, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	//if err != nil {
	//	panic(err)
	//}
	//
	//gin.DefaultWriter = io.MultiWriter(f) //输出到文件
	//gin.DefaultWriter = io.MultiWriter(f, os.Stdout) //输出到文件和控制台

	r := gin.Default() //不要使用logger
	//r := gin.New()
	//r.Use(gin.Recovery())

	// 中间件
	r.Use(middlewares.Cors())
	r.Use(middlewares.Logs())

	// 路由
	//r.POST("upload", v0.Upload)
	//r.GET("download/:filename", v0.Download)

	rv1 := r.Group("v1")
	rv1.POST("contract_invoke", v1.ContractInvoke)
	//rv1.POST("contract_query", v1.ContractQuery)
	rv1.POST("test", v1.Test)
	//rv1.POST("contract_deploy", v1.ContractDeploy)
	//rv1.POST("create_account", v1.CreateAccount)
	//rv1.POST("create_contract_account", v1.CreateContractAccount)
	//rv1.POST("balance", v1.Balance)
	//rv1.POST("transfer", v1.Transfer)
	//rv1.POST("query_tx", v1.QueryTx)
	//rv1.POST("method_acl", v1.MethodAcl)
	//rv1.POST("account_acl", v1.AccountAcl)
	//rv1.POST("status", v1.Status)
	////rv1.POST("group_chain", v1.GroupChain)
	////rv1.POST("group_node", v1.GroupNode) // 后续改造这两个接口为群组申请/邀请以及审核
	//rv1.POST("query_acl", v1.QueryAcl)
	//rv1.POST("query_block", v1.QueryBlock)
	//rv1.POST("query_list", v1.QueryLists)
	//rv1.POST("create_chain", v2.CreateChain)
	//rv1.POST("get_netURL", v2.GetNetURL)
	//rv1.POST("query_miners", v2.QueryMiners)
	//rv1.POST("addr_trans", v2.AddressTransfer)

	// 新版
	// post: https://url.example.com/v1/transaction	转账、查询
	// post: https://url.example.com/v1/block/(id\height)	块相关，根据id/height查询

	// todo
	// post: https://url.example.com/v1/contract/(acl)	合约相关（invoke/query）
	// post https://url.example.com/v1/contract/kernel/(governToken/proposal/parachain...)	系统合约相关

	// post: https://url.example.com/v1/consensus/(status)	共识相关
	// post https://url.example.com/v1/contract-account	合约账户相关
	// post: https://url.example.com/v1/account/(balance)	链上地址账户相关
	// post: https://url.example.com/v1/chain/(status)	链相关

	// post: https://url.example.com/v1/user/tx-history	需要配合中心化db实现
	return r
}
