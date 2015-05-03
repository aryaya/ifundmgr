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
	c.Layout = "layout.html"
	c.TplNames = "view.html"
}

type ScIssueController struct {
	beego.Controller
}

func (c *ScIssueController) Get() {
	c.Layout = "layout.html"
	c.TplNames = "sc/form.html"
}

type ScDepositController struct {
	beego.Controller
}

func (c *ScDepositController) Get() {
	c.Layout = "layout.html"
	c.TplNames = "sc/form.html"
}
