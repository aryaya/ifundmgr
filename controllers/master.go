//
//
//

package controllers

import (
	"github.com/astaxie/beego"
)

type MasterController struct {
	beego.Controller
}

func (c *MasterController) Get() {
	c.Layout = "master.html"
	c.TplNames = "view.html"
}
