// Copyright 2015 iCloudFund. All Rights Reserved.

package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type GateBankAccount struct {
	Name       string   // 开户人姓名
	BankName   string   // 银行名字
	BankId     string   // 银行账号
	Currencies []string // 支持的货币种类
}

type HortWallet struct {
	Name, AccountId, Secret string
}

type Fees struct {
	Min, Max float64 // 最低, 最高费率
	Rate     float64 // 费率比率
}

type Config struct {
	GBAs       []GateBankAccount // 收款人信息
	Currencies []string          // 支持的货币种类
	ColdWallet string            // 网关钱包地址 用于发行
	HoltWallet []HortWallet      // 网关钱包地址 用于支付
	ServerAddr string            // Server 地址
	Timeout    int               // 超过Timeout小时请求没有审批, 则超时关闭
	Roles      []Role            //
	Fees       Fees
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

// m01 1qaz2wsx iLyHPoNvsdN7LFFvpy8GGkDJGV1xo3V5We
// m02 1qaz2wsx iMxKojv7vNYyca7YdsVBSRKCmvGciHUNap
// m03 1qaz2wsx iLwUZfEo8pB9VTxzDtjwJBBuiWuVpLxz9m

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
	ColdWallet: "iN8sGowQCg1qptWcJG1WyTmymKX7y9cpmr", // ss1TCkz333t3t2J5eobcEMkMY3bXk // w01
	HoltWallet: []HortWallet{{"w02", "iwsZ7gxHFzu2xbj8YMf4UvK1PrDEMuxGkf", "ss9qoFiFNkokVfgrb3YkKHido6a1q"}, {"w03", "ine3v1DStiLfncJiCEgyfFct1i9t6M7z9E", "snwh9xAzpVoh9WxRc3pVENBJV44fj"}},
	ServerAddr: "wss://local.icloud.com:19528",
	Timeout:    24,
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

func initConf() {
	conf, err := loadConf()
	if err != nil {
		conf = defaultConf
		b, err := json.MarshalIndent(conf, "", " ")
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
