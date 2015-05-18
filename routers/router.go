package routers

import (
	"github.com/astaxie/beego"
	"ifundmgr/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "get:Get")
	beego.Router("/signin", &controllers.MainController{}, "get:SigninGet;post:SigninPost")
	beego.Router("/signout", &controllers.MainController{}, "post:SignoutPost")

	beego.Router("/issue", &controllers.MainController{}, "get:IssuesGet;post:IssuesPost")
	beego.Router("/deposit", &controllers.MainController{}, "get:DepositsGet;post:DepositsPost")
	beego.Router("/redeem", &controllers.MainController{}, "get:RedeemsGet;post:RedeemsPost")
	beego.Router("/withdrawal", &controllers.MainController{}, "get:WithdrawalsGet;post:WithdrawalsPost")

	beego.Router("/issue/add", &controllers.MainController{}, "get:AddIssueGet;post:AddIssuePost")
	beego.Router("/deposit/add", &controllers.MainController{}, "get:AddDepositGet;post:AddDepositPost")

	beego.Router("/issue/verify?:id", &controllers.MainController{}, "post:VerifyIssue")
	beego.Router("/deposit/verify?:id", &controllers.MainController{}, "post:VerifyDeposit")
	beego.Router("/redeem/verify?:id", &controllers.MainController{}, "post:VerifyIssue")
	beego.Router("/withdrawal/verify?:id", &controllers.MainController{}, "post:VerifyDeposit")

	// beego.Router("/issue/detailes?:id", &controllers.MainController{}, "get:DetaileIssue")
	// beego.Router("/deposit/detailes?:id", &controllers.MainController{}, "get:DetaileDeposit")
	// beego.Router("/redeem/detailes?:id", &controllers.MainController{}, "get:DetaileIssue")
	// beego.Router("/withdrawal/detailes?:id", &controllers.MainController{}, "get:DetaileDeposit")
}
