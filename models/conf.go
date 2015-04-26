//
//
//

package models

import (
	"encoding/json"
)

type Recipient struct {
	Name       string
	BankName   string
	BankId     string
	Currencies []string // 支持的货币种类
}

type Config struct {
	IssueRecipients   []Recipient // 收款人信息
	DepositRecipients []Recipient // 收款人信息
	Currencies        []string    // 支持的货币种类
	Wallet            []string    // 网关钱包地址
	Timeout           int         // 超过Timeout小时请求没有审批, 则超时关闭
}
