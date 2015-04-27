package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.TplNames = "index.html"
}

type SignInController struct {
	beego.Controller
}

func (c *SignInController) Get() {

}

type SignOutController struct {
	beego.Controller
}

func (c *SignOutController) Get() {
	c.TplNames = "index.html"
}
