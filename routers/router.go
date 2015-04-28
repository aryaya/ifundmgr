package routers

import (
	"github.com/astaxie/beego"
	"ifundmgr/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "get:Get")
	beego.Router("/signin", &controllers.MainController{}, "get:SigninGet; post:SigninPost")
	beego.Router("/signout", &controllers.MainController{}, "get:SigninGet")

	beego.Router("/issue", &controllers.MainController{}, "get::IssuesGet; post:IssuesPost")
	beego.Router("/deposit", &controllers.MainController{}, "get::DepositsGet; post:DepositsPost")
	beego.Router("/redeem", &controllers.MainController{}, "get::RedeemsGet; post:RedeemsPost")
	beego.Router("/withdrawal", &controllers.MainController{}, "get::WithdrawalsGet; post: WithdrawalsPost")

	beego.Router("/issue/add.action", &controllers.MainController{}, "post:AddIssue")
	beego.Router("/deposit/add.action", &controllers.MainController{}, "post:AddDeposit")
	beego.Router("/issue/verify.action", &controllers.MainController{}, "post:VerifyIssue")
	beego.Router("/deposit/verify.action", &controllers.MainController{}, "post:VerifyDeposit")
	beego.Router("/redeem/verify.action", &controllers.MainController{}, "post:VerifyIssue")
	beego.Router("/withdrawal/verify.action", &controllers.MainController{}, "post:VerifyDeposit")
}
