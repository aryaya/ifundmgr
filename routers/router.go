package routers

import (
	"github.com/astaxie/beego"
	"ifundmgr/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
}
