//
//
//

package controllers

import (
	"github.com/astaxie/beego"
)

type FinController struct {
	beego.Controller
}

func (c *FinController) Get() {
	c.Layout = "fin.html"
	c.TplNames = "view.html"
}

type FinApproveController struct {
	beego.Controller
}

func (c *FinApproveController) Post() {
	c.Layout = "fin.html"
	// c.TplNames = "view.html"
}
