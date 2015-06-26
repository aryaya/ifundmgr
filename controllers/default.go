// Copyright 2015 iCloudFund. All Rights Reserved.

package controllers

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/wangch/ifundmgr/models"
)

func Init() {
	models.Init()

	beego.SetStaticPath("/js", "static/js")
	beego.SetStaticPath("/css", "static/css")
	beego.SetStaticPath("/fonts", "static/fonts")

	beego.SessionOn = true
	beego.SessionName = "icloudsessionid"
	beego.InsertFilter("/", beego.BeforeRouter, filterUser)
	beego.InsertFilter("/*", beego.BeforeRouter, filterUser)

	beego.AddFuncMap("showVerify", showVerify)
	beego.AddFuncMap("getGbas", getGbas)
	beego.AddFuncMap("getHoltWallets", getHoltWallets)
	beego.AddFuncMap("canModifyGBankId", canModifyGBankId)
	beego.AddFuncMap("canModifyGWallet", canModifyGWallet)
	beego.AddFuncMap("fmtStatus", fmtStatus)
	beego.AddFuncMap("issue", issue)
}

var filterUser = func(ctx *context.Context) {
	_, ok := ctx.Input.Session("Role").(*models.Role)
	if !ok && ctx.Request.RequestURI != "/signin" {
		ctx.Redirect(302, "/signin")
	}
}

func issue(currency string) string {
	return ""
}

func fmtStatus(status int) string {
	for k, v := range models.StatusMap {
		if v == status {
			return k
		}
	}
	return "Unkown Status"
}

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Role"] = c.GetSession("Role").(*models.Role)
	c.Layout = "layout.html"
	c.TplNames = "index.html"
}

func (c *MainController) SigninGet() {
	c.Data["ShowSignin"] = true
	c.Layout = "layout.html"
	c.TplNames = "index.html"
	h := md5.New()
	io.WriteString(h, "wangch"+time.Now().String())
	token := fmt.Sprintf("%x", h.Sum(nil))
	c.Data["Token"] = token
	c.SetSession("Token", token)
}

func (c *MainController) SignoutPost() {
	c.Redirect("/signin", 302)
}

func (c *MainController) SigninPost() {
	uname := c.GetString("Name")
	upass := c.GetString("Password")
	token := c.GetString("Token")
	log.Println(token)
	if token != c.GetSession("Token").(string) {
		c.Redirect("/signin", 302)
		return
	}
	r := &models.Role{Username: uname}
	err := models.Gorm.Read(r, "Username")
	if err != nil {
		c.Redirect("/signin", 302)
		return
	}
	if r.Password != models.PassHash(upass) {
		c.Redirect("/signin", 302)
		return
	}
	c.SetSession("Role", r)
	c.Redirect("/", 302)
}

func (c *MainController) getScRole() *models.Role {
	r := c.GetSession("Role").(*models.Role)
	if r.Type != models.RoleC {
		c.Redirect("/", 302)
		return nil
	}
	c.Data["Role"] = r
	return r
}

func (c *MainController) getNoScRole() *models.Role {
	r := c.GetSession("Role").(*models.Role)
	if r.Type == models.RoleC {
		c.Redirect("/", 302)
		return nil
	}
	c.Data["Role"] = r
	return r
}

func (c *MainController) csHtml(isDeposit, ok bool, r *models.Role) {
	c.Data["ShowSignin"] = false
	c.Data["IsDeposit"] = isDeposit
	c.Data["OK"] = ok
	c.Data["Gbas"] = models.Gconf.GBAs
	c.Layout = "layout.html"
	c.TplNames = "form.html"
}

func AddReq(roleName, gid string, t int, u *models.User, currency string, amount, fees float64) error {
	gbas := models.Gconf.GBAs
	var gba *models.GateBankAccount
	for _, g := range gbas {
		if g.BankId == gid {
			gba = &g
			break
		}
	}

	req := &models.Request{
		CsId:      roleName,
		CsTime:    time.Now(),
		UName:     u.UName,
		UWallet:   u.UWallet,
		UBankName: u.UBankName,
		UBankId:   u.UBankId,
		UContact:  u.UContact,
		Currency:  currency,
		Amount:    amount,
		Fees:      fees,
		Type:      t,
	}

	if gba == nil && len(gbas) > 0 {
		gba = &gbas[0]
	}

	tm := time.Unix(0, 0)

	rec := &models.Recoder{
		Status: models.COK,
		R:      req,
		FTime:  tm,
		MTime:  tm,
		ATime:  tm,
	}
	if gba != nil {
		// req.GName = gba.Name
		// req.GBankName = gba.BankName
		rec.GBankId = gba.BankId
	}
	// log.Printf("%#v", rec)
	req.R = rec
	_, err := models.Gorm.Insert(rec)
	if err != nil {
		log.Fatal(err)
	}
	_, err = models.Gorm.Insert(req)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (c *MainController) addReq(role *models.Role, t int) error {
	amount, err := c.GetFloat("amount")
	if err != nil {
		return err
	}
	fees, err := c.GetFloat("fees")
	if err != nil {
		return err
	}
	u := &models.User{
		UName:     c.GetString("name"),
		UWallet:   c.GetString("iccWallet"),
		UBankName: c.GetString("bankName"),
		UBankId:   c.GetString("bankId"),
		UContact:  c.GetString("contact"),
	}

	return AddReq(role.Username, c.GetString("Gba"), t, u, c.GetString("currency"), amount, fees)
}

func (c *MainController) AddIssueGet() {
	r := c.getScRole()
	if r == nil {
		return
	}
	c.Data["Role"] = r
	c.csHtml(false, false, r)
}

func (c *MainController) AddIssuePost() {
	r := c.getScRole()
	if r == nil {
		return
	}
	err := c.addReq(r, models.Issue)
	if err != nil {
		return
	}
	c.csHtml(false, true, r)
}

func (c *MainController) AddDepositGet() {
	r := c.getScRole()
	if r == nil {
		return
	}
	c.csHtml(true, false, r)
}

func (c *MainController) AddDepositPost() {
	r := c.getScRole()
	if r == nil {
		return
	}
	err := c.addReq(r, models.Deposit)
	if err != nil {
		return
	}
	c.csHtml(true, true, r)
}

var tableHeaders = []string{
	"时间",
	"用户姓名",
	"用户钱包",
	"用户银行",
	"用户银行账号",
	"货币",
	"金额",
	"费用",
	"状态",
	"网关银行账号",
	"网关钱包",
	"详情",
	"审核",
}

type HtmlReq struct {
	*models.Request
	Rec     *models.Recoder
	Role    *models.Role
	Status  string
	GBankId string
}

func (c *MainController) getReqs(role *models.Role, tname string, status int, st, et *time.Time, t int) []HtmlReq {
	var reqs []*models.Request
	log.Println("getReqs:", tname, status, st, et)
	qs := models.Gorm.QueryTable(tname).Filter("Type", t).Filter("CsTime__gte", st).Filter("CsTime__lte", et).RelatedSel()
	if status != -1 {
		qs.Filter("R__Status", status).All(&reqs)
	} else {
		qs.All(&reqs)
	}
	hreqs := make([]HtmlReq, len(reqs))
	for i, r := range reqs {
		log.Println(r.R)
		hreqs[i] = HtmlReq{
			Request: r,
			Rec:     r.R,
			Role:    role,
			Status:  models.RStatusMap[r.R.Status],
		}
	}
	return hreqs
}

func (c *MainController) queryTable() {
	r := c.getNoScRole()
	if r == nil {
		return
	}
	sst := c.GetString("stime")
	set := c.GetString("etime")

	st, err := time.Parse("2006-01-02", sst)
	if err != nil {
		log.Fatal(err)
	}
	et, err := time.Parse("2006-01-02", set)
	if err != nil {
		log.Fatal(err)
	}
	c.SetSession("StartDate", &st)
	c.SetSession("EndDate", &et)
	c.SetSession("Status", c.GetString("status"))
}

func (c *MainController) getTable(t int) {
	tname := "request"
	r := c.getNoScRole()
	if r == nil {
		return
	}

	st, ok := c.GetSession("StartDate").(*time.Time)
	if !ok {
		tst := time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local)
		st = &tst
	}
	et, ok := c.GetSession("EndDate").(*time.Time)
	if !ok {
		tet := time.Now()
		et = &tet
	}
	status := -1
	ss, ok := c.GetSession("Status").(string)
	if ok {
		status = models.StatusMap[ss]
	}
	reqs := c.getReqs(r, tname, status, st, et, t)

	c.Data["StartDate"] = st.Format("2006-01-02")
	c.Data["EndDate"] = et.Format("2006-01-02")
	c.Data["Requests"] = reqs
	c.Data["TableHeaders"] = tableHeaders
	c.Data["StatusSlice"] = models.StatusSlice
	c.Data["Status"] = ss
	c.Data["Gbas"] = models.Gconf.GBAs

	c.Layout = "layout.html"
	c.TplNames = "reqtable.html"
}

func (c *MainController) IssuesGet() {
	c.getTable(models.Issue)
}

func (c *MainController) IssuesPost() {
	c.queryTable()
	c.Redirect("/issue/", 302)
}

func (c *MainController) DepositsGet() {
	c.getTable(models.Deposit)
}

func (c *MainController) DepositsPost() {
	c.queryTable()
	c.Redirect("/deposit/", 302)
}

func (c *MainController) RedeemsGet() {
	c.getTable(models.Redeem)
}

func (c *MainController) RedeemsPost() {
	c.queryTable()
	c.Redirect("/redeem/", 302)
}

func (c *MainController) WithdrawalsGet() {
	c.getTable(models.Withdrawal)
}

func (c *MainController) WithdrawalsPost() {
	c.queryTable()
	c.Redirect("/withdrawal/", 302)
}

func canVerify(rtype, status, typ int) int {
	if rtype == models.RoleC {
		return -1
	}
	if rtype == models.RoleF {
		if status == models.COK {
			return models.FOK
		}
	}
	if rtype == models.RoleM {
		if status == models.FOK {
			return models.MOK
		}
	}
	if rtype == models.RoleA {
		if status == models.MOK {
			return models.AOK
		}
	}
	return -1
}

func canModifyGBankId(hr *HtmlReq) bool {
	if hr.Role.Type == models.RoleF && hr.Rec.Status == models.COK {
		return true
	}
	return false
}

func canModifyGWallet(hr *HtmlReq) bool {
	if hr.Role.Type == models.RoleF && hr.Rec.Status == models.COK {
		return true
	}
	return false
}

func getGbas() []models.GateBankAccount {
	return models.Gconf.GBAs
}

func getHoltWallets() []string {
	hws := make([]string, len(models.Gconf.HoltWallet))
	for i, x := range models.Gconf.HoltWallet {
		hws[i] = x.Name + ":" + x.AccountId
	}
	return hws
}

func showVerify(hr *HtmlReq) bool {
	return canVerify(hr.Role.Type, hr.Rec.Status, hr.Type) != -1
}

func (c *MainController) verify(isOut bool) {
	r := c.getNoScRole()
	if r == nil {
		return
	}
	sid := c.Ctx.Request.URL.Query().Get("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		log.Println(err)
		return
	}
	rec := &models.Recoder{Id: int64(id)}
	err = models.Gorm.Read(rec)
	if err != nil {
		log.Println(err)
		return
	}
	newStatus := canVerify(r.Type, rec.Status, 0 /*rec.R.Type*/)
	if newStatus == -1 {
		log.Println("CAN'T verify")
		return
	}
	if r.Type == models.RoleF {
		if rec.GWallet == "" || rec.GBankId == "" {
			return
		}
		rec.FId = r.Username
		rec.FTime = time.Now()
		rec.Status = models.FOK
	} else if r.Type == models.RoleM {
		rec.MId = r.Username
		rec.MTime = time.Now()
		rec.Status = models.MOK
	} else if r.Type == models.RoleA {
		rec.AId = r.Username
		rec.ATime = time.Now()
		rec.Status = models.AOK
		// 会计审批, 直接发送
		if isOut {
			err := models.Payment(rec.R, rec.GWallet)
			if err != nil {
				//
			}
		}
	} else {
		log.Println("r.Type error", r.Type)
		return
	}
	models.Gorm.Update(rec)
}

func deposit(req *models.Request) error {
	return nil
}

func (c *MainController) VerifyIssue() {
	c.verify()
	c.Redirect("/issue/", 302)
}

func (c *MainController) VerifyDeposit() {
	c.verify()
	c.Redirect("/deposit/", 302)
}

func (c *MainController) VerifyRedeem() {
	c.verify()
	c.Redirect("/redeem/", 302)
}

func (c *MainController) VerifyWithdrawal() {
	c.verify()
	c.Redirect("/withdrawal/", 302)
}

func (c *MainController) updateGbank() {
	r := c.getNoScRole()
	if r == nil {
		return
	}
	if r.Type != models.RoleF {
		return
	}
	sid := c.Ctx.Request.URL.Query().Get("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		log.Println(err)
		return
	}
	rec := &models.Recoder{Id: int64(id)}
	err = models.Gorm.Read(rec)
	if err != nil {
		log.Println(err)
		return
	}
	gbankid := c.GetString("Gba")
	if rec.GBankId == gbankid {
		return
	}
	rec.GBankId = gbankid
	models.Gorm.Update(rec)
}

func (c *MainController) IssueUpdateGbank() {
	c.updateGbank()
	c.Redirect("/issue/", 302)
}

func (c *MainController) DepositUpdateGbank() {
	c.updateGbank()
	c.Redirect("/deposit/", 302)
}

func (c *MainController) RedeemUpdateGbank() {
	c.updateGbank()
	c.Redirect("/redeem/", 302)
}

func (c *MainController) WithdrawalUpdateGbank() {
	c.updateGbank()
	c.Redirect("/withdrawal/", 302)
}

func (c *MainController) updateHotwallet() {
	r := c.getNoScRole()
	if r == nil {
		return
	}
	if r.Type != models.RoleF {
		return
	}
	sid := c.Ctx.Request.URL.Query().Get("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		log.Println(err)
		return
	}
	rec := &models.Recoder{Id: int64(id)}
	err = models.Gorm.Read(rec)
	if err != nil {
		log.Println(err)
		return
	}
	hw := c.GetString("HotWallet")
	if rec.GWallet == hw {
		return
	}
	rec.GWallet = hw
	models.Gorm.Update(rec)
}

func (c *MainController) IssueUpdateHotwallet() {
	c.updateHotwallet()
	c.Redirect("/issue/", 302)
}

func (c *MainController) DepositUpdateHotwallet() {
	c.updateHotwallet()
	c.Redirect("/deposit/", 302)
}

func (c *MainController) RedeemUpdateHotwallet() {
	c.updateHotwallet()
	c.Redirect("/redeem/", 302)
}

func (c *MainController) WithdrawalUpdateHotwallet() {
	c.updateHotwallet()
	c.Redirect("/withdrawal/", 302)
}
