package v1

import (
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/superconsensus/matrix-sdk-go/v2/account"
	"github.com/superconsensus/matrix-sdk-go/v2/xuper"
	"github.com/superconsensus/open-api-server/common"
	"github.com/superconsensus/open-api-server/conf"
	"github.com/superconsensus/open-api-server/controllers"
	"log"
	"net/http"
)

// Test 测试接口，购买治理代币
func Test(c *gin.Context) {
	// go channel
	req := new(controllers.Req)
	_ = c.ShouldBind(req)

	args := make(map[string][]byte)
	for s, s2 := range req.Args {
		args[s] = []byte(s2)
	}

	acc, _ := account.RetrieveAccount(req.Mnemonic, conf.Req.Language)
	if acc.Address == "" {
		acc.Address = req.Account
	}

	request, err := xuper.NewRequest(acc, "xkernel", "$govern_token", "Buy", args, "testa", req.Args["money"], xuper.WithBcname(req.BcName))
	if err != nil {
		common.Record(c, "test failed for new request", err.Error())
		log.Println("test new request error", err)
		return
	}

	xclient, err := xuper.New(req.Node, xuper.WithConfigFile("./conf/sdk.yaml"))
	if err != nil {
		common.Record(c, "调用失败", err.Error())
		log.Println("new xclient failed, error=", err)
		return
	}
	defer func() {
		closeErr := xclient.Close()
		if closeErr != nil {
			log.Println("create chain: close xclient failed, error=", closeErr)
		}
	}()

	// 包括preExec（主要校验参数和计算gas）跟complete（预执行没有问题之后sign）两部分
	genTx, err := xclient.GenerateTx(request)
	if err != nil {
		common.Record(c, "gen tx failed", err.Error())
		return
	} else {
		// 改了下sdk，这里能获取到预执行的gas消耗了
		fmt.Println("genTx: ", genTx.GasUsed, hex.EncodeToString(genTx.Tx.Txid), hex.EncodeToString(genTx.DigestHash))
	}
	// 前面complete已经使用私钥签名了，这里重新sign画蛇添足——sign方法是开放给多签名的
	//err = genTx.Sign(acc)
	// 提交已经sign的交易
	postTx, err := xclient.PostTx(genTx)

	// way 2
	//tx, err := xclient.PreExecTx(request)
	//postTx, err := xclient.Do(request)

	/*////// ------ 如果交易构建跟签名分开的话，这里需要改成这样
	configFile, err := config.GetConfig("./conf/sdk.yaml")
	if err != nil {
		common.Record(c, "获取配置文件失败", err.Error())
		log.Println("get config failed", err)
		return
	}
	prop, err := xuper.NewProposal(xclient, request, configFile)
	if err != nil {
		common.Record(c, "new proposal failed", err.Error())
		log.Println("new proposal failed", err)
		return
	}
	err = prop.PreExecWithSelectUtxo()
	if err != nil {
		common.Record(c, "预执行失败", err.Error())
		log.Println("pre exec select utxo error=", err)
		return
	}
	// 这里构建交易，但不签名，等待发送
	preTx, err := prop.GenCompleteTx()
	// send to Tee 签名
	response, err := SendToTee(preTx) // preTx.Sign()

	postTx, err := xclient.PostTx(preTx)
	if err != nil {
		common.Record(c, "提交交易失败", err.Error())
		log.Println("post tx failed, error=", err)
		return
	}////// ------*/

	if err != nil {
		common.Record(c, "test failed for post tx", err.Error())
		log.Println("post tx failed, error=", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "调用成功",
		"resp": controllers.Result{
			Txid:    hex.EncodeToString(postTx.Tx.Txid),
			GasUsed: postTx.GasUsed,
			Data:    string(postTx.ContractResponse.Body),
		},
	})
}
