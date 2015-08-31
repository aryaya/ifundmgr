// Copyright 2015 iCloudFund. All Rights Reserved.

package main

import (
	"github.com/astaxie/beego"
	"github.com/wangch/glog"
	"github.com/wangch/ifundmgr/routers"
)

func main() {
	glog.SetLogDirs(".")
	glog.SetLogToStderr(true)

	routers.Init()
	beego.Run()
}
