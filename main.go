package main

import (
	"github.com/astaxie/beego"
	_ "ifundmgr/routers"
)

func main() {
	beego.SetStaticPath("/js", "static/js")
	beego.SetStaticPath("/css", "static/css")
	beego.SetStaticPath("/fonts", "static/fonts")

	beego.Run()
}
