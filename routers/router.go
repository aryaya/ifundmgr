package routers

import (
	"github.com/astaxie/beego"
	"ifundmgr/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "get:Get")
	beego.Router("/signin", &controllers.MainController{}, "get:SigninGet;post:SigninPost")
	beego.Router("/signout", &controllers.MainController{}, "get:SigninGet")

	beego.Router("/issue", &controllers.MainController{}, "get:IssuesGet;post:IssuesPost")
	beego.Router("/deposit", &controllers.MainController{}, "get:DepositsGet;post:DepositsPost")
	beego.Router("/redeem", &controllers.MainController{}, "get:RedeemsGet;post:RedeemsPost")
	beego.Router("/withdrawal", &controllers.MainController{}, "get:WithdrawalsGet;post:WithdrawalsPost")

	beego.Router("/issue/add", &controllers.MainController{}, "get:AddIssueGet;post:AddIssuePost")
	beego.Router("/deposit/add", &controllers.MainController{}, "get:AddDepositGet;post:AddDepositPost")

	beego.Router("/issue/verify", &controllers.MainController{}, "get:VerifyIssue")
	beego.Router("/deposit/verify", &controllers.MainController{}, "get:VerifyDeposit")
	beego.Router("/redeem/verify", &controllers.MainController{}, "get:VerifyIssue")
	beego.Router("/withdrawal/verify", &controllers.MainController{}, "get:VerifyDeposit")

	beego.Router("/issue/detailes", &controllers.MainController{}, "get:DetaileIssue")
	beego.Router("/deposit/detailes", &controllers.MainController{}, "get:DetaileDeposit")
	beego.Router("/redeem/detailes", &controllers.MainController{}, "get:DetaileIssue")
	beego.Router("/withdrawal/detailes", &controllers.MainController{}, "get:DetaileDeposit")
}
