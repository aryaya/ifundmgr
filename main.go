// Copyright 2015 iCloudFund. All Rights Reserved.

package main

import (
	"github.com/wangch/ifundmgr/routers"

	"github.com/astaxie/beego"
)

func main() {
	routers.Init()
	beego.Run()
}
