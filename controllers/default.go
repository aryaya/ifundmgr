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

func init() {
	beego.SetStaticPath("/js", "static/js")
	beego.SetStaticPath("/css", "static/css")
	beego.SetStaticPath("/fonts", "static/fonts")

	beego.SessionOn = true
	beego.SessionName = "icloudsessionid"
	beego.InsertFilter("/", beego.BeforeRouter, filterUser)
	beego.InsertFilter("/*", beego.BeforeRouter, filterUser)

	beego.AddFuncMap("showVerify", showVerify)
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
	c.Data["Gbas"] = models.Gconf.GBAs
	c.Data["ShowSignin"] = false
	c.Data["IsDeposit"] = isDeposit
	c.Data["OK"] = ok
	c.Data["Gbas"] = models.Gconf.GBAs
	c.Layout = "layout.html"
	c.TplNames = "form.html"
}

func AddReq(roleName, gid string, t int, u *models.User, currency string, amount, fees float64) error {
	// gbas := models.Gconf.GBAs
	// var gba models.GateBankAccount
	// for _, g := range gbas {
	// 	if g.BankId == gid {
	// 		gba = g
	// 		break
	// 	}
	// }

	req := &models.Request{
		CsId:      roleName,
		CsTime:    time.Now(),
		UName:     u.UName,
		UWallet:   u.UWallet,
		UBankName: u.UBankName,
		UBankId:   u.UBankId,
		UContact:  u.UContact,
		// GName:     gba.Name,
		// GBankName: gba.BankName,
		// GBankId:   gba.BankId,
		Currency: currency,
		Amount:   amount,
		Fees:     fees,
		Type:     t,
	}
	rec := &models.Recoder{
		Status: models.COK,
		R:      req,
	}
	req.R = rec
	models.Gorm.Insert(rec)
	models.Gorm.Insert(req)
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
		UWallet:   c.GetString("icc_wallet"),
		UBankName: c.GetString("bank_name"),
		UBankId:   c.GetString("bank_id"),
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
	// "ID",
	// "客服ID",
	"时间",
	"用户姓名",
	"钱包地址",
	"银行名称",
	"银行账号",
	// "收款人真实姓名",
	// "收款人银行名称",
	// "收款人银行账号",
	"货币",
	"金额",
	"费用",
	"状态",
	"审核",
	"详情",
}

type HtmlReq struct {
	*models.Request
	Rec    *models.Recoder
	Role   *models.Role
	Status string
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

func showVerify(hr *HtmlReq) bool {
	return canVerify(hr.Role.Type, hr.Rec.Status, hr.Type) != -1
}

func (c *MainController) verify() {
	// tname := "recoder"
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
	log.Println("@@@:", rec.R)
	newStatus := canVerify(r.Type, rec.Status, 0 /*rec.R.Type*/)
	if newStatus == -1 {
		log.Println("CAN'T verify")
		return
	}
	if r.Type == models.RoleF {
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
