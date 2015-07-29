// Copyright 2015 iCloudFund. All Rights Reserved.

package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego/context"
	"github.com/wangch/glog"
	"github.com/wangch/ifundmgr/models"
	"github.com/wangch/ripple/data"
)

type QuoteResp struct {
	Result       string `json:"result"`
	Error        string `json:"error,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	QuoteJson    *Quote `json:"quote,omitempty"`
}

func quoteErrorResp(msg string) *QuoteResp {
	return &QuoteResp{
		Result:       "error",
		Error:        "-1",
		ErrorMessage: msg,
	}
}

type Quote struct {
	Address        string        `json:"address"`
	DestinationTag uint          `json:"destination_tag"`
	InvoiceID      string        `json:"invoice_id"`
	Send           []data.Amount `json:"send"`
}

func sendResp(resp interface{}, ctx *context.Context) error {
	b, err := json.Marshal(resp)
	if err != nil {
		glog.Error(err)
		return err
	}
	glog.Info(string(b))
	ctx.ResponseWriter.Write(b)
	return nil
}

// https://ripplecn.com/bridge?type=quote&amount=1%2FCNY&destination=z&address=ra5tSyQ2cvJUHfAvEdmC89HKSKZTn7xXMw&alipay_account=aa&full_name=bb&contact_info=cc
func (c *MainController) Quote() {
	sa := c.Ctx.Request.URL.Query().Get("amount")
	address := c.Ctx.Request.URL.Query().Get("address")
	bank_name := c.Ctx.Request.URL.Query().Get("bank_name")
	card_number := c.Ctx.Request.URL.Query().Get("card_number")
	full_name := c.Ctx.Request.URL.Query().Get("full_name")
	contact_info := c.Ctx.Request.URL.Query().Get("contact_info")

	a, err := data.NewAmount(sa)
	if err != nil {
		glog.Error(err)
		resp := quoteErrorResp("the query amount err")
		sendResp(resp, c.Ctx)
		return
	}
	glog.Info("Quote:", address, bank_name, card_number, full_name, contact_info, a, a.IsNative())
	sv := a.Value.String()
	am, err := strconv.ParseFloat(sv, 64)
	if err != nil {
		glog.Error(err)
		resp := quoteErrorResp("the query amount err2")
		sendResp(resp, c.Ctx)
		return
	}

	if a.IsNative() {
		pv, _ := data.NewNativeValue(1e6)
		a.Value, _ = a.Value.Divide(*pv)
	}

	currency := a.Currency.String()
	fee := am * models.Gconf.Fees.Rate
	min := models.Gconf.Fees.FeeMap[currency][0]
	max := models.Gconf.Fees.FeeMap[currency][1]
	if fee < min {
		fee = min
	}
	if fee > max {
		fee = max
	}
	acc, err := data.NewAccountFromAddress(models.Gconf.ColdWallet)
	if err != nil {
		glog.Fatal(err)
	}
	a.Issuer = *acc

	u := &models.User{
		UName:     full_name,
		UWallet:   address,
		UBankName: bank_name,
		UBankId:   card_number,
		UContact:  contact_info,
	}

	t := models.Withdrawal
	if a.IsNative() {
		t = models.Redeem
	}
	req, err := models.AddReq("", "", models.Gconf.ColdWallet, t, u, a.Currency.String(), am/1e6, fee)
	if err != nil {
		// glog.Error(err)
		return
	}

	fv, err := data.NewValue(fmt.Sprintf("%f", fee), a.IsNative())
	if err != nil {
		glog.Fatal(err)
	}

	a.Value, err = a.Value.Add(*fv)
	if err != nil {
		glog.Fatal(err)
	}

	if a.IsNative() {
		pv, _ := data.NewNativeValue(1e6)
		a.Value, _ = a.Value.Divide(*pv)
		a.Currency, _ = data.NewCurrency("ICC")
	}

	quote := &Quote{
		Address:        models.Gconf.ColdWallet,
		DestinationTag: 2147483647,
		Send:           []data.Amount{*a},
		InvoiceID:      req.InvoiceId,
	}
	resp := &QuoteResp{
		Result:    "success",
		QuoteJson: quote,
	}
	glog.Info("Quote:", address, full_name, "OK")
	sendResp(resp, c.Ctx)
}
