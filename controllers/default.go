package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
<<<<<<< HEAD
	c.TplNames = "index.html"
=======
<<<<<<< HEAD
	c.TplNames = "index.html"
=======
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplNames = "index.tpl"
>>>>>>> 754bd7dc12e6c787998b7161e58b9a989f5b53a6
>>>>>>> 29e087821ca7217599d8822c845d7770821100d0
}
