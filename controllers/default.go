package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	"ifundmgr/models"
	"log"
	"strconv"
	"time"
)

var filterUser = func(ctx *context.Context) {
	_, ok := ctx.Input.Session("uid").(string)
	if !ok && ctx.Request.RequestURI != "/signin" {
		ctx.Redirect(302, "/signin")
	}
}

func init() {
	beego.SetStaticPath("/js", "static/js")
	beego.SetStaticPath("/css", "static/css")
	beego.SetStaticPath("/fonts", "static/fonts")

	beego.SessionOn = true
	beego.SessionName = "icloudsessionid"
	beego.InsertFilter("/*", beego.BeforeRouter, filterUser)

	beego.AddFuncMap("showVerify", showVerify)
	beego.AddFuncMap("fmtStatus", fmtStatus)
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
	c.Layout = "layout.html"
	c.TplNames = "info.html"
}

func (c *MainController) SigninGet() {
	c.Data["ShowSignin"] = true
	c.Layout = "layout.html"
	c.TplNames = "info.html"
	c.Data["Token"] = "someToken"
}

func (c *MainController) passHash(password string) string {
	return password
}

func (c *MainController) SigninPost() {
	uname := c.GetString("Name")
	upass := c.GetString("Password")
	token := c.GetString("Token")
	if token != c.Data["Token"] {
		c.Redirect("/login", 302)
		return
	}
	r := &models.Role{Username: uname}
	err := models.Gorm.Read(r)
	if err != nil {
		c.Redirect("/login", 302)
		return
	}
	if r.PasswordHash != c.passHash(upass) {
		c.Redirect("/login", 302)
		return
	}
	c.SetSession("Role", r)
	c.Data["Role"] = r
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

func (c *MainController) addReq(isDeposit bool, role *models.Role) error {
	bankId := c.GetString("BankId")
	gbas := models.Gconf.GBAs
	var gba models.GateBankAccount
	for _, g := range gbas {
		if g.BankId == bankId {
			gba = g
			break
		}
	}

	amount, err := c.GetFloat("amount")
	if err != nil {
		return err
	}
	fees, err := c.GetFloat("fees")
	if err != nil {
		return err
	}

	req := &models.Request{
		CsId:      role.Username,
		CsTime:    time.Now(),
		CName:     c.GetString("depositorName"),
		CWallet:   c.GetString("depositorWallet"),
		CBankName: c.GetString("depositorBankName"),
		CBankId:   c.GetString("depositorBankId"),
		GName:     gba.Name,
		GBankName: gba.BankName,
		GBankId:   gba.BankId,
		Currence:  c.GetString("currency"),
		Amount:    amount,
		Fees:      fees,
	}
	rec := &models.Recoder{
		Status: models.COK,
		R:      req,
	}
	req.R = rec
	models.Gorm.Insert(req)
	models.Gorm.Insert(rec)
	return nil
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
	err := c.addReq(false, r)
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
	err := c.addReq(true, r)
	if err != nil {
		return
	}
	c.csHtml(true, true, r)
}

var tableHeaders = []string{
	// "ID",
	// "客服ID",
	// "客服提交时间",
	"存款人真实姓名",
	"存款人钱包地址",
	"存款人银行名称",
	"存款人银行账号",
	// "收款人真实姓名",
	// "收款人银行名称",
	"收款人银行账号",
	"货币",
	"金额",
	"费用",
	"状态",
	"详情",
	"审核",
}

type HtmlReq struct {
	*models.Request
	Rec   *models.Recoder
	Role  *models.Role
	Tname string
}

func (c *MainController) getReqs(role *models.Role, tname string, status int, st, et *time.Time) []HtmlReq {

	var reqs []*models.Request
	qs := models.Gorm.QueryTable(tname).Filter("CsTime__gte", st).Filter("CsTime__lte", et)
	if status != -1 {
		qs.Filter("Recoder__Status", status).All(reqs)
	} else {
		qs.All(reqs)
	}
	hreqs := make([]HtmlReq, len(reqs))
	for i, r := range reqs {
		hreqs[i] = HtmlReq{
			Request: r,
			Rec:     r.R,
			Role:    role,
			Tname:   tname,
		}
	}
	return hreqs
}

// func (c *MainController) reqHtml(st, et *time.Time, reqs []HtmlReq) {
// 	c.Data["StartDate"] = st.Format("2015-05-02")
// 	c.Data["EndDate"] = et.Format("2015-05-02")
// 	c.Data["Requests"] = reqs
// 	c.Data["TableHeaders"] = tableHeaders
// 	c.Layout = "layout.html"
// 	c.TplNames = "reqtable.html"
// }

func (c *MainController) queryTable(tname string) {
	r := c.getNoScRole()
	if r == nil {
		return
	}
	sst := c.GetString("stime")
	set := c.GetString("etime")
	st, err := time.Parse("2015-05-02", sst)
	if err != nil {
		log.Fatal(err)
	}
	et, err := time.Parse("2015-05-02", set)
	if err != nil {
		log.Fatal(err)
	}
	status, ok := models.StatusMap[c.GetString("status")]
	if !ok {
		status = -1
	}
	c.SetSession("StartDate", st)
	c.SetSession("StartDate", et)
	c.SetSession("Status", status)
}

func (c *MainController) getTable(tname string) {
	r := c.getNoScRole()
	if r == nil {
		return
	}

	st, ok := c.GetSession("StartDate").(*time.Time)
	if !ok {
		tst := time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local)
		st = &tst
	}
	et, ok := c.GetSession("EndTime").(*time.Time)
	if !ok {
		tet := time.Now()
		et = &tet
	}
	status, ok := c.GetSession("Status").(int)
	reqs := c.getReqs(r, tname, status, st, et)

	c.Data["StartDate"] = st.Format("2015-05-02")
	c.Data["EndDate"] = et.Format("2015-05-02")
	c.Data["Requests"] = reqs
	c.Data["TableHeaders"] = tableHeaders
	c.Layout = "layout.html"
	c.TplNames = "reqtable.html"
}

func (c *MainController) IssuesGet() {
	c.getTable("issue_req")
}

func (c *MainController) IssuesPost() {
	c.queryTable("issue_req")
	c.Redirect("/issue", 302)
}

func (c *MainController) DepositsGet() {
	c.getTable("deposit_req")
}

func (c *MainController) DepositsPost() {
	c.queryTable("deposit_req")
	c.Redirect("/deposit", 302)
}

func (c *MainController) RedeemsGet() {
	c.getTable("redeem_req")
}

func (c *MainController) RedeemsPost() {
	c.queryTable("redeem_req")
	c.Redirect("/redeem", 302)
}

func (c *MainController) WithdrawalsGet() {
	c.getTable("withdrawal_req")
}

func (c *MainController) WithdrawalsPost() {
	c.queryTable("withdrawal_req")
	c.Redirect("/withdrawal", 302)
}

func canVerify(rtype, status int, tname string) int {
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
	return canVerify(hr.Role.Type, hr.Rec.Status, hr.Tname) != -1
}

func (c *MainController) verify(tname string) {
	r := c.getNoScRole()
	if r == nil {
		return
	}
	sid := c.Ctx.Request.URL.Query().Get("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		log.Fatal(err)
	}
	rec := &models.Recoder{}
	qs := models.Gorm.QueryTable(tname).Filter("Id", id)
	err = qs.One(rec)
	if err != nil {
		log.Fatal(err, id)
	}
	newStatus := canVerify(r.Type, rec.Status, tname)
	if newStatus == -1 {
		return
	}
	n, err := qs.Update(orm.Params{
		"Status": newStatus,
	})
	if err != nil {
		log.Fatal(err)
	}
	if n != 1 {
		log.Fatal("update error")
	}
}

func (c *MainController) VerifyIssue() {
	c.verify("issue_rec")
	c.Redirect("/issue", 302)
}

func (c *MainController) VerifyDeposit() {
	c.verify("deposit_rec")
	c.Redirect("/deposit", 302)
}

func (c *MainController) VerifyRedeem() {
	c.verify("redeem_rec")
	c.Redirect("/redeem", 302)
}

func (c *MainController) VerifyWithdrawal() {
	c.verify("withdrawal_rec")
	c.Redirect("/withdrawal", 302)
}

func (c *MainController) DetaileIssue() {
	c.verify("issue_rec")
}

func (c *MainController) DetaileDeposit() {
	c.verify("issue_rec")
}

func (c *MainController) DetaileRedeem() {
	c.verify("issue_rec")
}

func (c *MainController) DetaileWithdrawal() {
	c.verify("issue_rec")
}
