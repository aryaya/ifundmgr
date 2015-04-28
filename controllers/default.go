package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

var filterUser = func(ctx *context.Context) {
	ssid, ok := ctx.Input.Session("uid").(string)
	if !ok && ctx.Request.RequestURI != "/signin" {
		ctx.Redirect(302, "/signin")
		return
	}

	// 根据ssid找出内存中对应的用户
	// 找不到, 302重新登录
	// 根据用户角色重定向
}

func init() {
	beego.SetStaticPath("/js", "static/js")
	beego.SetStaticPath("/css", "static/css")
	beego.SetStaticPath("/fonts", "static/fonts")

	beego.SessionOn = true
	beego.SessionName = "icloudsessionid"
	beego.InsertFilter("/*", beego.BeforeRouter, filterUser)
}

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Layout = "layout.html"
	c.TplNames = "info.html"
}

func (c *MainController) SigninGet() {
	c.Data["ShowSignin"] = true
	c.Layout = "layout.html"
	c.TplNames = "info.html"
}

func (c *MainController) SigninPost() {

}

func (c *MainController) IssuesGet() {
	c.Layout = "layout.html"
	c.TplNames = "info.html"

}

func (c *MainController) IssuesPost() {

}

func (c *MainController) DepositsGet() {

}

func (c *MainController) DepositsPost() {

}

func (c *MainController) RedeemsGet() {

}

func (c *MainController) RedeemsPost() {

}

func (c *MainController) WithdrawalsGet() {

}

func (c *MainController) WithdrawalsPost() {

}

func (c *MainController) AddIssue() {

}

func (c *MainController) AddDeposit() {

}

func (c *MainController) VerifyIssue() {

}

func (c *MainController) VerifyDeposit() {

}

func (c *MainController) VerifyRedeem() {

}

func (c *MainController) VerifyWithdrawal() {

}
