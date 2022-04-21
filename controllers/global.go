package controllers

import (
	"github.com/superconsensus/matrix-sdk-go/v2/xuper"
)

type Req struct {
	RequestId             string            `json:"request_id,omitempty"`
	Node                  string            `json:"node,omitempty" binding:"required"`
	BcName                string            `json:"bc_name,omitempty"`
	Mnemonic              string            `json:"mnemonic,omitempty"`
	To                    string            `json:"to,omitempty"`
	Amount                int64             `json:"amount,omitempty"`
	Fee                   int64             `json:"fee,omitempty"`
	Desc                  string            `json:"desc,omitempty"`
	EndorseServiceHost    string            `json:"endorse_service_host,omitempty"`
	EndorseServiceFee     int               `json:"endorse_service_fee,omitempty"`
	EndorseServiceFeeAddr string            `json:"endorse_service_fee_addr,omitempty"`
	EndorseServiceAddr    string            `json:"endorse_service_addr,omitempty"`
	Crypto                string            `json:"crypto,omitempty"`
	ContractAccount       string            `json:"contract_account,omitempty"`
	ContractName          string            `json:"contract_name,omitempty"`
	ContractFile          string            `json:"-"`
	ContractCode          string            `json:"contract_code,omitempty"`
	ModuleName            string            `json:"module_name,omitempty"`
	MethodName            string            `json:"method_name,omitempty"`
	Args                  map[string]string `json:"args,omitempty"`
	Query                 bool              `json:"query,omitempty"`
	Upgrade               bool              `json:"upgrade,omitempty"`
	Runtime               string            `json:"runtime,omitempty"`
	Language              int               `json:"language,omitempty"`
	Address               []string          `json:"address,omitempty"`
	Txid                  string            `json:"txid,omitempty"`
	Method                string            `json:"method,omitempty"`
	BlockID               string            `json:"block_id,omitempty"`
	BlockHeight           int64             `json:"block_height,omitempty"`
	Account               string            `json:"account,omitempty"`
	Frozen                bool              `json:"frozen,omitempty"`
}

type Resp struct {
	RequestId string `json:"request_id,omitempty"`
	Code      int    `json:"code,omitempty"`
	Msg       string `json:"msg,omitempty"`
	Error     string `json:"error,omitempty"`
}

type Result struct {
	Txid            string     `json:"txid,omitempty"`
	Data            string     `json:"data,omitempty"`
	GasUsed         int64      `json:"gas_used,omitempty"`
	Mnemonic        string     `json:"mnemonic,omitempty"`
	Address         string     `json:"address,omitempty"`
	AccountAcl      *xuper.ACL `json:"account_acl,omitempty"`
	AccountBalance  string     `json:"account_balance,omitempty"`
	Tx              string     `json:"tx,omitempty"`
	ContractAccount string     `json:"contract_account,omitempty"`
}
