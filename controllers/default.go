// Copyright 2015 iCloudFund. All Rights Reserved.

package controllers

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/wangch/glog"
	"github.com/wangch/ifundmgr/models"
)

func Init() {
	models.Init()

	beego.SetStaticPath("/js", "static/js")
	beego.SetStaticPath("/css", "static/css")
	beego.SetStaticPath("/fonts", "static/fonts")
	beego.SetStaticPath("/certificates", "./certificates")

	beego.SessionOn = true
	beego.SessionName = "icloudsessionid"
	beego.InsertFilter("/*", beego.BeforeRouter, filterUser)

	beego.AddFuncMap("showVerify", showVerify)
	beego.AddFuncMap("getGbas", getGbas)
	beego.AddFuncMap("getGbaName", getGbas)
	beego.AddFuncMap("getHoltWallets", getHoltWallets)
	beego.AddFuncMap("canModifyGBankId", canModifyGBankId)
	beego.AddFuncMap("canModifyGWallet", canModifyGWallet)
	beego.AddFuncMap("fmtStatus", fmtStatus)
	beego.AddFuncMap("issue", issue)
}

var filterUser = func(ctx *context.Context) {
	us := ctx.Request.URL.String()
	if strings.Contains(us, "quote") || strings.Contains(us, "deposit") {
		return
	}
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
	r, ok := c.GetSession("Role").(*models.Role)
	if !ok {
		c.Redirect("/signin", 302)
		return
	}
	c.Data["Role"] = r
	c.SetSession("ErrMsg", "")
	c.Layout = "layout.html"
	c.TplNames = "index.html"
}

func RandToken() string {
	h := md5.New()
	io.WriteString(h, "wangch"+time.Now().String())
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (c *MainController) SigninGet() {
	c.Data["ShowSignin"] = true
	c.Layout = "layout.html"
	c.TplNames = "index.html"
	token := RandToken()
	c.Data["ErrMsg"] = c.GetSession("ErrMsg")
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
	if t, ok := c.GetSession("Token").(string); !ok || token != t {
		c.Redirect("/signin", 302)
		return
	}
	r := &models.Role{Username: uname}
	err := models.Gorm.Read(r, "Username")
	if err != nil {
		glog.Error("Sign in error: " + uname + " is NOT in database")
		c.SetSession("ErrMsg", "name OR password error")
		c.Redirect("/signin", 302)
		return
	}
	if r.Password != models.PassHash(upass) {
		glog.Error("Sign in error: password error")
		c.SetSession("ErrMsg", "name OR password error")
		c.Redirect("/signin", 302)
		return
	}
	c.SetSession("Role", r)
	c.Redirect("/", 302)
}

func (c *MainController) getRole() *models.Role {
	r := c.GetSession("Role").(*models.Role)
	c.Data["Role"] = r
	return r
}

const hextable = "0123456789ABCDEF"

func b2h(h []byte) []byte {
	b := make([]byte, len(h)*2)
	for i, v := range h {
		b[i*2] = hextable[v>>4]
		b[i*2+1] = hextable[v&0x0f]
	}
	return b
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
	qs := models.Gorm.QueryTable(tname).Filter("Type", t).Filter("CTime__gte", st).Filter("CTime__lte", et).RelatedSel()
	if status != -1 {
		qs.Filter("R__Status", status).All(&reqs)
	} else {
		qs.All(&reqs)
	}
	hreqs := make([]HtmlReq, len(reqs))
	for i, r := range reqs {
		glog.Infof("%+v", r.R)
		hreqs[i] = HtmlReq{
			Request: r,
			Rec:     r.R,
			Role:    role,
			Status:  models.RStatusMap[r.R.Status],
		}
	}
	return hreqs
}

var PDErr error = errors.New("Permission denied for the user")

func (c *MainController) queryTable() error {
	r := c.getRole()
	if r == nil {
		return PDErr
	}
	sst := c.GetString("stime")
	set := c.GetString("etime")

	st, err := time.Parse("2006-01-02", sst)
	if err != nil {
		glog.Error(err)
		return err
	}
	et, err := time.Parse("2006-01-02", set)
	if err != nil {
		glog.Error(err)
		return err
	}
	c.SetSession("StartDate", &st)
	c.SetSession("EndDate", &et)
	c.SetSession("Status", c.GetString("status"))
	return nil
}

func (c *MainController) getTable(t int) error {
	tname := "request"
	r := c.getRole()
	if r == nil {
		return PDErr
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
	c.Data["StatusSlice"] = models.StatusSlice
	c.Data["Status"] = ss
	c.Data["Gbas"] = models.Gconf.GBAs
	c.Data["Type"] = t

	c.Layout = "layout.html"
	c.TplNames = "reqtable.html"
	return nil
}

func (c *MainController) IssuesGet() {
	err := c.getTable(models.Issue)
	if err != nil {
		c.Data["ErrMsg"] = err.Error()
	}
}

func (c *MainController) IssuesPost() {
	err := c.queryTable()
	if err != nil {
		c.Data["ErrMsg"] = err.Error()
	}

	c.Redirect("/issue/", 302)
}

func (c *MainController) DepositsGet() {
	err := c.getTable(models.Deposit)
	if err != nil {
		c.Data["ErrMsg"] = err.Error()
	}
}

func (c *MainController) DepositsPost() {
	err := c.queryTable()
	if err != nil {
		c.Data["ErrMsg"] = err.Error()
	}

	c.Redirect("/deposit/", 302)
}

func (c *MainController) RedeemsGet() {
	err := c.getTable(models.Redeem)
	if err != nil {
		c.Data["ErrMsg"] = err.Error()
	}
}

func (c *MainController) RedeemsPost() {
	err := c.queryTable()
	if err != nil {
		c.Data["ErrMsg"] = err.Error()
	}

	c.Redirect("/redeem/", 302)
}

func (c *MainController) WithdrawalsGet() {
	err := c.getTable(models.Withdrawal)
	if err != nil {
		c.Data["ErrMsg"] = err.Error()
	}
}

func (c *MainController) WithdrawalsPost() {
	err := c.queryTable()
	if err != nil {
		c.Data["ErrMsg"] = err.Error()
	}

	c.Redirect("/withdrawal/", 302)
}

func canVerify(rtype, status, typ int) int {
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
	if hr.Role.Type == models.RoleF &&
		hr.Rec.Status == models.COK &&
		(hr.Type == models.Withdrawal || hr.Type == models.Redeem) {
		return true
	}
	return false
}

func canModifyGWallet(hr *HtmlReq) bool {
	glog.Infoln(hr.Role.Type == models.RoleF, hr.Rec.Status == models.COK, (hr.Type == models.Deposit || hr.Type == models.Issue))
	if hr.Role.Type == models.RoleF &&
		hr.Rec.Status == models.COK &&
		(hr.Type == models.Deposit || hr.Type == models.Issue) {
		return true
	}
	glog.Infoln("go here")
	return false
}

func getGbaName(g models.GateBankAccount) string {
	return g.Name + " " + g.BankName + " " + g.BankId
}

func getGbas() []models.GateBankAccount {
	return models.Gconf.GBAs
}

func getHoltWallets() []string {
	hws := make([]string, len(models.Gconf.HoltWallet))
	for i, x := range models.Gconf.HoltWallet {
		hws[i] = x.Name + ":" + x.AccountId
	}
	glog.Infoln(hws)
	return hws
}

func showVerify(hr *HtmlReq) bool {
	return canVerify(hr.Role.Type, hr.Rec.Status, hr.Type) != -1
}

func (c *MainController) verify(isOut bool) error {
	r := c.getRole()
	if r == nil {
		return PDErr
	}
	sid := c.Ctx.Request.URL.Query().Get("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		return err
	}
	// rec := &models.Recoder{Id: int64(id)}
	var req models.Request
	err = models.Gorm.QueryTable("Request").Filter("R__id", id).RelatedSel().One(&req)
	if err != nil {
		return err
	}
	rec := req.R
	newStatus := canVerify(r.Type, rec.Status, 0 /*rec.R.Type*/)
	if newStatus == -1 {
		return PDErr
	}
	if r.Type == models.RoleF {
		if rec.GWallet == "" || rec.GBankId == "" {
			err = errors.New("HotWallet OR GateBandId config error")
			glog.Error(err)
			return err
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
		// 会计审批, 直接发送
		if isOut {
			if rec.R == nil {
				return PDErr
			}
			sender := ""
			if strings.Contains(rec.GWallet, ":") {
				sender = strings.Split(rec.GWallet, ":")[1]
			}
			if sender == "" {
				return errors.New("HotWallet error")
			}
			rec.R.Currency = getCurrencyID(rec.R.Currency)
			err := models.Payment(rec.R, sender)
			if err != nil {
				return err
			}
			rec.Status = models.AOK
		} else { // 回收和取款
			rec.Status = models.OKC // 转账完成则整个记录完成
		}
	} else {
		glog.Fatal("can't go here")
	}
	_, err = models.Gorm.Update(rec)
	return err
}

func (c *MainController) VerifyIssue() {
	err := c.verify(true)
	if err != nil {
		c.Data["ErrMsg"] = err.Error()
	}
	c.Redirect("/issue/", 302)
}

func (c *MainController) VerifyDeposit() {
	err := c.verify(true)
	if err != nil {
		c.Data["ErrMsg"] = err.Error()
	}
	c.Redirect("/deposit/", 302)
}

func (c *MainController) VerifyRedeem() {
	err := c.verify(false)
	if err != nil {
		c.Data["ErrMsg"] = err.Error()
	}
	c.Redirect("/redeem/", 302)
}

func (c *MainController) VerifyWithdrawal() {
	err := c.verify(false)
	if err != nil {
		c.Data["ErrMsg"] = err.Error()
	}
	c.Redirect("/withdrawal/", 302)
}

func (c *MainController) updateGbank() error {
	r := c.getRole()
	if r == nil {
		return PDErr
	}
	if r.Type != models.RoleF {
		return PDErr
	}
	sid := c.Ctx.Request.URL.Query().Get("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		return err
	}
	rec := &models.Recoder{Id: int64(id)}
	err = models.Gorm.Read(rec)
	if err != nil {
		return err
	}
	gbankid := c.GetString("Gba")
	if rec.GBankId == gbankid {
		return errors.New("GateBankID config error")
	}
	rec.GBankId = gbankid
	_, err = models.Gorm.Update(rec)
	return err
}

func (c *MainController) IssueUpdateGbank() {
	err := c.updateGbank()
	if err != nil {
		glog.ErrorDepth(1, err)
		c.Data["ErrMsg"] = err.Error()
	}
	c.Redirect("/issue/", 302)
}

func (c *MainController) DepositUpdateGbank() {
	err := c.updateGbank()
	if err != nil {
		glog.ErrorDepth(1, err)
		c.Data["ErrMsg"] = err.Error()
	}
	c.Redirect("/deposit/", 302)
}

func (c *MainController) RedeemUpdateGbank() {
	err := c.updateGbank()
	if err != nil {
		glog.ErrorDepth(1, err)
		c.Data["ErrMsg"] = err.Error()
		return
	}
	c.Redirect("/redeem/", 302)
}

func (c *MainController) WithdrawalUpdateGbank() {
	err := c.updateGbank()
	if err != nil {
		glog.ErrorDepth(1, err)
		c.Data["ErrMsg"] = err.Error()
		return
	}
	c.Redirect("/withdrawal/", 302)
}

func (c *MainController) updateHotwallet() error {
	r := c.getRole()
	if r == nil {
		return PDErr
	}
	if r.Type != models.RoleF {
		return PDErr
	}
	sid := c.Ctx.Request.URL.Query().Get("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		return err
	}
	rec := &models.Recoder{Id: int64(id)}
	err = models.Gorm.Read(rec)
	if err != nil {
		return err
	}
	hw := c.GetString("HotWallet")
	if rec.GWallet == hw {
		return errors.New("HotWallet config error")
	}
	rec.GWallet = hw
	_, err = models.Gorm.Update(rec)
	return err
}

func (c *MainController) IssueUpdateHotwallet() {
	err := c.updateHotwallet()
	if err != nil {
		glog.ErrorDepth(1, err)
		c.Data["ErrMsg"] = err.Error()
	}
	c.Redirect("/issue/", 302)
}

func (c *MainController) DepositUpdateHotwallet() {
	err := c.updateHotwallet()
	if err != nil {
		glog.ErrorDepth(1, err)
		c.Data["ErrMsg"] = err.Error()
	}
	c.Redirect("/deposit/", 302)
}

func (c *MainController) RedeemUpdateHotwallet() {
	err := c.updateHotwallet()
	if err != nil {
		glog.ErrorDepth(1, err)
		c.Data["ErrMsg"] = err.Error()
	}
	c.Redirect("/redeem/", 302)
}

func (c *MainController) WithdrawalUpdateHotwallet() {
	err := c.updateHotwallet()
	if err != nil {
		glog.ErrorDepth(1, err)
		c.Data["ErrMsg"] = err.Error()
	}
	c.Redirect("/withdrawal/", 302)
}

func getCurrencyID(s string) string {
	if s == "ICC" {
		return "ICC"
	}
	switch s {
	case "港元", "HKD":
		return "HKD"
	case "美元元", "USD":
		return "USD"
	case "日元", "JPY":
		return "JPY"
	case "欧元", "EUR":
		return "EUR"
	case "人民币", "CNY":
		return "CNY"
	}
	return ""
}
