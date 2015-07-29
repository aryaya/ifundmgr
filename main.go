// Copyright 2015 iCloudFund. All Rights Reserved.

package main

import (
	"flag"

	"github.com/astaxie/beego"
	"github.com/wangch/ifundmgr/routers"
)

func main() {
	flag.Parse()

	routers.Init()
	beego.Run()
}
