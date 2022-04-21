package v1

import (
	"encoding/hex"
	"errors"
	"github.com/superconsensus/open-api-server/common"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/superconsensus/matrix-sdk-go/v2/account"
	"github.com/superconsensus/open-api-server/conf"
	"github.com/superconsensus/open-api-server/controllers"
	//log "github.com/superconsensus/open-api-server/utils"
	"github.com/superconsensus/matrix-sdk-go/v2/xuper"
)

func ContractInvoke(c *gin.Context) {
	req := new(controllers.Req)
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数无效",
		})
		log.Printf("param invalid, err: %s", err.Error())
		return
	}

	// todo 账号地址从db拿
	acc, err := account.RetrieveAccount(req.Mnemonic, conf.Req.Language)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "助记词无效",
		})
		log.Printf("mnemonic can not retrieve account, err: %s", err.Error())
		return
	}

	xclient, err := xuper.New(req.Node)
	if err != nil {
		return
	}
	defer func() {
		closeErr := xclient.Close()
		if closeErr != nil {
			log.Println("contract invoke: close xclient failed, error=", closeErr)
		}
	}()

	// 第一次先用普通账户invoke
	tx, err := invoke(req, acc, xclient)
	if err != nil {
		// 失败的话再通过合约账户来调用本次操作
		if req.ContractAccount != "" {
			setContractE := acc.SetContractAccount(req.ContractAccount)
			if setContractE != nil {
				log.Printf("contract invoke: set contract account failed, error=", setContractE)
				common.Record(c, "调用失败", setContractE.Error())
				return
			}
			tx, err = invoke(req, acc, xclient)
			// 第二次通过合约账户调用仍然失败的话就真的是有错误了
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":  400,
					"msg":   "操作失败",
					"error": err.Error(),
				})
				log.Printf("if query [%v] contract fail, err: %s", req.Query, err.Error())
				return
			}
		} else {
			log.Println("contract invoke failed, error=", err)
			common.Record(c, "调用失败", err.Error())
		}
	}
	if !req.Query {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "调用成功",
			"resp": controllers.Result{
				Txid:    hex.EncodeToString(tx.Tx.Txid),
				Data:    string(tx.ContractResponse.Body),
				GasUsed: tx.GasUsed,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "查询成功",
			"resp": controllers.Result{
				Data: string(tx.ContractResponse.Body),
			},
		})
	}
}

// 封装合约调用/查询
func invoke(req *controllers.Req, acc *account.Account, xclient *xuper.XClient) (*xuper.Transaction, error) {
	tx := &xuper.Transaction{}
	var err error
	// 默认调用c/c++合约
	if req.ModuleName == "wasm" || req.ModuleName == "" {
		if req.Query {
			tx, err = xclient.QueryWasmContract(acc, req.ContractName, req.MethodName, req.Args, xuper.WithBcname(req.BcName))
		} else {
			tx, err = xclient.InvokeWasmContract(acc, req.ContractName, req.MethodName, req.Args, xuper.WithBcname(req.BcName))
		}
	} else if req.ModuleName == "native" {
		if req.Query {
			tx, err = xclient.QueryNativeContract(acc, req.ContractName, req.MethodName, req.Args, xuper.WithBcname(req.BcName))
		} else {
			tx, err = xclient.InvokeNativeContract(acc, req.ContractName, req.MethodName, req.Args, xuper.WithBcname(req.BcName))
		}
	} else if req.ModuleName == "evm" {
		if req.Query {
			tx, err = xclient.QueryEVMContract(acc, req.ContractName, req.MethodName, req.Args, xuper.WithBcname(req.BcName))
		} else {
			tx, err = xclient.InvokeEVMContract(acc, req.ContractName, req.MethodName, req.Args, xuper.WithBcname(req.BcName))
		}
	} else {
		return nil, errors.New("module_name参数缺失或格式错误")
	}
	return tx, err
}
