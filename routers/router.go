package routers

import (
	"github.com/astaxie/beego"
	"ifundmgr/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/sc", &controllers.ScController{})
	beego.Router("/fin", &controllers.FinController{})
	beego.Router("/master", &controllers.MasterController{})
}
