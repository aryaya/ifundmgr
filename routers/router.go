// Copyright 2015 iCloudFund. All Rights Reserved.

package routers

import (
	"github.com/astaxie/beego"
	"github.com/wangch/ifundmgr/controllers"
)

func Init() {
	controllers.Init()

	beego.Router("/", &controllers.MainController{}, "get:Get")
	beego.Router("/signin", &controllers.MainController{}, "get:SigninGet;post:SigninPost")
	beego.Router("/signout", &controllers.MainController{}, "post:SignoutPost")

	beego.Router("/issue", &controllers.MainController{}, "get:IssuesGet;post:IssuesPost")
	beego.Router("/deposit", &controllers.MainController{}, "get:DepositsGet;post:DepositsPost")
	beego.Router("/redeem", &controllers.MainController{}, "get:RedeemsGet;post:RedeemsPost")
	beego.Router("/withdrawal", &controllers.MainController{}, "get:WithdrawalsGet;post:WithdrawalsPost")

	// beego.Router("/issue/add", &controllers.MainController{}, "get:AddIssueGet;post:AddIssuePost")
	// beego.Router("/deposit/add", &controllers.MainController{}, "get:AddDepositGet;post:AddDepositPost")

	beego.Router("/issue/verify?:id", &controllers.MainController{}, "post:VerifyIssue")
	beego.Router("/deposit/verify?:id", &controllers.MainController{}, "post:VerifyDeposit")
	beego.Router("/redeem/verify?:id", &controllers.MainController{}, "post:VerifyRedeem")
	beego.Router("/withdrawal/verify?:id", &controllers.MainController{}, "post:VerifyWithdrawal")

	beego.Router("/issue/gbankid?:id", &controllers.MainController{}, "post:IssueUpdateGbank")
	beego.Router("/deposit/gbankid?:id", &controllers.MainController{}, "post:DepositUpdateGbank")
	beego.Router("/redeem/gbankid?:id", &controllers.MainController{}, "post:RedeemUpdateGbank")
	beego.Router("/withdrawal/gbankid?:id", &controllers.MainController{}, "post:WithdrawalUpdateGbank")

	beego.Router("/issue/hotwallet?:id", &controllers.MainController{}, "post:IssueUpdateHotwallet")
	beego.Router("/deposit/hotwallet?:id", &controllers.MainController{}, "post:DepositUpdateHotwallet")
	beego.Router("/redeem/hotwallet?:id", &controllers.MainController{}, "post:RedeemUpdateHotwallet")
	beego.Router("/withdrawal/hotwallet?:id", &controllers.MainController{}, "post:WithdrawalUpdateHotwallet")

	beego.Router("/api/quote", &controllers.MainController{}, "get:ApiQuote")

	beego.Router("/api/deposit", &controllers.MainController{}, "get:ApiDepositGet")
	beego.Router("/api/deposit/amount", &controllers.MainController{}, "post:ApiDepositAmountPost")
	beego.Router("/api/deposit/add", &controllers.MainController{}, "post:ApiDepositPost")

	beego.Router("/api/buyicc", &controllers.MainController{}, "get:ApiDepositGet")
	beego.Router("/api/buyicc/amount", &controllers.MainController{}, "post:ApiDepositAmountPost")
	beego.Router("/api/buyicc/add", &controllers.MainController{}, "post:ApiDepositPost")
}
