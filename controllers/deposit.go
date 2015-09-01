//
//
//

package controllers

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/wangch/glog"
	"github.com/wangch/ifundmgr/models"
)

func (c *MainController) addReq(t int) error {
	amt, ok := c.GetSession("Amt").(*AmountInfo)
	if !ok {
		c.Ctx.Request.ParseMultipartForm(1024 * 1024 * 10)
		glog.Infof("%#+v", c.Ctx.Request)
		gbankId := c.GetString("gbankId")
		currency := c.GetString("currency")
		fee, err := c.GetFloat("fees")
		if err != nil {
			glog.Errorln(err)
			return err
		}
		a, err := c.GetFloat("amount")
		if err != nil {
			glog.Errorln(err)
			return err
		}
		amt = &AmountInfo{
			BankId:   gbankId,
			Currency: currency,
			Amount:   a,
			Fees:     fee,
		}
	}
	uname := c.GetString("name")
	uwallet := c.GetString("iccWallet")
	ubankName := c.GetString("bankName")
	ubankId := c.GetString("bankId")
	ucontact := c.GetString("contact")

	rf, header, err := c.Ctx.Request.FormFile("certificate")
	if err != nil {
		glog.Error(err)
		return err
	}
	defer rf.Close()

	filePath := "./certificates/" + RandToken() + "-" + header.Filename

	wf, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		glog.Error(err)
		return err
	}
	io.Copy(wf, rf)
	wf.Close()

	u := &models.User{
		UName:        uname,
		UWallet:      uwallet,
		UBankName:    ubankName,
		UBankId:      ubankId,
		UContact:     ucontact,
		UCertificate: filePath[1:],
	}

	_, err = models.AddReq(amt.BankId, "", t, u, amt.Currency, amt.Amount, amt.Fees)
	if err != nil {
		glog.Errorln(err)
		return err
	}
	return nil
}

type AmountInfo struct {
	BankName string
	BankId   string
	Currency string
	Amount   float64
	Fees     float64
	Total    float64
}

func (c *MainController) ApiDepositGet() {
	currencys := models.Gconf.Currencies
	upath := c.Ctx.Request.URL.Path
	if strings.Contains(upath, "buyicc") {
		c.Data["BuyIcc"] = true
	} else {
		c.Data["BuyIcc"] = false
	}
	c.Data["Currencies"] = currencys
	c.Data["USDRate"] = models.Gconf.UsdRate
	c.TplNames = "deposit.html"
}

func getGba(currency string) (*models.GateBankAccount, error) {
	gbas := models.Gconf.GBAs
	for _, g := range gbas {
		for _, c := range g.Currencies {
			if currency == c {
				return &g, nil
			}
		}
	}
	return nil, errors.New("getGba error: currency is NOT in gbas of conf")
}

func (c *MainController) ApiDepositAmountPost() {
	glog.Infof("%#+v", c.Ctx.Request)
	upath := c.Ctx.Request.URL.Path
	buyIcc := false
	if strings.Contains(upath, "buyicc") {
		buyIcc = true
	}

	currency := "USD"
	if !buyIcc {
		s := c.GetString("currency")
		if s == "" {
			glog.Errorln("currency is nil")
			return
		}
		currency = getCurrencyID(s)
	}

	gba, err := getGba(currency)
	if err != nil {
		glog.Errorln(err)
		return
	}

	amount, err := c.GetFloat("amount")
	if err != nil {
		glog.Errorln(err)
		return
	}

	if buyIcc {
		amount *= models.Gconf.UsdRate
	}

	fees := models.Gconf.Fees

	min, max := fees.FeeMap[currency][0], fees.FeeMap[currency][1]
	fee := amount * fees.Rate
	if fee < min {
		fee = min
	}
	if fee > max {
		fee = max
	}

	amt := &AmountInfo{
		BankName: gba.BankName,
		BankId:   gba.BankId,
		Currency: currency,
		Amount:   amount,
		Fees:     fee,
		Total:    amount + fee,
	}
	c.SetSession("Amt", amt)
	currencys := models.Gconf.Currencies
	c.Data["Currencies"] = currencys
	c.Data["Amt"] = amt
	c.Data["BuyIcc"] = buyIcc
	c.TplNames = "deposit.html"
}

func (c *MainController) ApiDepositPost() {
	upath := c.Ctx.Request.URL.Path
	buyIcc := false
	if strings.Contains(upath, "buyicc") {
		buyIcc = true
	}
	if buyIcc {
		c.addReq(models.Issue)
		c.Redirect("/api/buyicc", 302)
	} else {
		c.addReq(models.Deposit)
		c.Redirect("/api/deposit", 302)
	}
}
