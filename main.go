package main

import (
	"github.com/astaxie/beego"
	_ "ifundmgr/routers"
)

func main() {
	println("go here 0")
	beego.Run()
}
