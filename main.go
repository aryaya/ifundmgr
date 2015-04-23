package main

import (
	"github.com/astaxie/beego"
	_ "ifundmgr/routers"
)

func main() {
<<<<<<< HEAD
	beego.SetStaticPath("/js", "static/js")
	beego.SetStaticPath("/css", "static/css")
	beego.SetStaticPath("/fonts", "static/fonts")
=======
<<<<<<< HEAD
	beego.SetStaticPath("/js", "static/js")
	beego.SetStaticPath("/css", "static/css")
	beego.SetStaticPath("/fonts", "static/fonts")
=======
>>>>>>> 754bd7dc12e6c787998b7161e58b9a989f5b53a6
>>>>>>> 29e087821ca7217599d8822c845d7770821100d0
	beego.Run()
}
