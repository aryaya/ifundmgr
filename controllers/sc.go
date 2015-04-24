//
//
//

package controllers

import (
	"github.com/astaxie/beego"
)

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

var issueRecipientHtml = `<div class="input-group">
   <input type="radio" name="" id=""> <h3>{{.Name}}</h3>
   <p>{{.BankName}} {{.BankId}} </p>
</div>
`
var depositRecipientHtml = `<div class="input-group">
   <input type="radio" name="" id=""> <h3>{{.Name}}</h3>
   <p>{{.BankName}} {{.BankId}} </p>
</div>
`

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
