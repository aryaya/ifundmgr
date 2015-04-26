//
//
//

package controllers

import (
	"github.com/astaxie/beego"
)

var tableHeader = []string{
	"ID",
	"客服ID",
	"客服提交时间",
	"存款人真实姓名",
	"存款人钱包地址",
	"存款人银行名称",
	"存款人银行账号",
	"收款人真实姓名",
	"收款人银行名称",
	"收款人银行账号",
	"货币",
	"金额",
	"费用",
	"状态",
}

type ScController struct {
	beego.Controller
}

func (c *ScController) Get() {
	c.Layout = "sc.html"
	c.TplNames = "view.html"
}

type ScIssueController struct {
	beego.Controller
}

func (c *ScIssueController) Get() {
	c.Layout = "sc.html"
	c.TplNames = "sc/issue.html"
}

type ScDepositController struct {
	beego.Controller
}

func (c *ScDepositController) Get() {
	c.Layout = "sc.html"
	c.TplNames = "sc/deposit.html"
}
