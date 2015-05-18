//
//
//

package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type GateBankAccount struct {
	Name       string
	BankName   string
	BankId     string
	Currencies []string // 支持的货币种类
}

type Config struct {
	GBAs       []GateBankAccount // 收款人信息
	Currencies []string          // 支持的货币种类
	Wallet     []string          // 网关钱包地址
	Timeout    int               // 超过Timeout小时请求没有审批, 则超时关闭
	Roles      []Role            //
}

var configFile = "./conf.json"

func loadConf() (*Config, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var conf Config
	err = json.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

var defaultConf = &Config{
	GBAs: []GateBankAccount{
		{
			Name:       "王春晖",
			BankName:   "中国工商银行武汉支行",
			BankId:     "1234 5678 1379 2468",
			Currencies: []string{"USD", "CNY", "HKD", "EUR", "JPY"},
		},
		{
			Name:       "马军",
			BankName:   "香港汇丰银行",
			BankId:     "5555 5555 5555 5555",
			Currencies: []string{"USD", "CNY", "HKD", "EUR", "JPY"},
		},
	},
	Currencies: []string{"USD", "CNY", "HKD", "EUR", "JPY"},
	// Wallet:     []string{""},
	Timeout: 24,
	Roles: []Role{
		{
			Username: "wangchC",
			Password: "passwordC",
			Type:     RoleC,
		},
		{
			Username: "wangchF",
			Password: "passwordF",
			Type:     RoleF,
		},
		{
			Username: "wangchM",
			Password: "passwordM",
			Type:     RoleM,
		},
		{
			Username: "wangchA",
			Password: "passwordA",
			Type:     RoleA,
		},
	},
}

func init() {
	conf, err := loadConf()
	if err != nil {
		conf = defaultConf
		b, err := json.Marshal(conf)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(configFile, b, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	Gconf = conf
}

var Gconf *Config
