// Copyright 2015 iCloudFund. All Rights Reserved.

package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/astaxie/beego/context"
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

// func quoteSucessResp(a *data.Amount) *QuoteResp {

// 	quote := &Quote{
// 		Address:        models.Gconf.ColdWallet,
// 		DestinationTag: 2147483647,
// 		Send:           []data.Amount{*a},
// 		InvoiceID:      models.GetInvoiceID(a),
// 	}
// 	return &QuoteResp{
// 		Result:    "success",
// 		QuoteJson: quote,
// 	}
// }

func sendResp(resp interface{}, ctx *context.Context) error {
	b, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	ctx.ResponseWriter.Write(b)
	return nil
}

// https://ripplecn.com/bridge?type=quote&amount=1%2FCNY&destination=z&address=ra5tSyQ2cvJUHfAvEdmC89HKSKZTn7xXMw&alipay_account=aa&full_name=bb&contact_info=cc
func (c *MainController) Quote() {
	// destination := c.Ctx.Request.URL.Query().Get("destination")
	sa := c.Ctx.Request.URL.Query().Get("amount")
	address := c.Ctx.Request.URL.Query().Get("address")
	bank_name := c.Ctx.Request.URL.Query().Get("bank_name")
	card_number := c.Ctx.Request.URL.Query().Get("card_number")
	full_name := c.Ctx.Request.URL.Query().Get("full_name")
	// opening_branch := c.Ctx.Request.URL.Query().Get("opening_branch")
	contact_info := c.Ctx.Request.URL.Query().Get("contact_info")
	a, err := data.NewAmount(sa)
	if err != nil {
		log.Println(err)
		resp := quoteErrorResp("the query amount err")
		sendResp(resp, c.Ctx)
		return
	}
	t := models.Withdrawal
	if a.IsNative() {
		t = models.Redeem
	}
	sv := a.Value.String()
	am, err := strconv.ParseFloat(sv, 64)
	if err != nil {
		log.Println(err)
		resp := quoteErrorResp("the query amount err2")
		sendResp(resp, c.Ctx)
		return
	}
	fee := am * models.Gconf.Fees.Rate
	if fee < models.Gconf.Fees.Min {
		fee = models.Gconf.Fees.Min
	}
	if fee > models.Gconf.Fees.Max {
		fee = models.Gconf.Fees.Max
	}
	u := &models.User{
		UName:     full_name,
		UWallet:   address,
		UBankName: bank_name,
		UBankId:   card_number,
		UContact:  contact_info,
	}

	req, err := AddReq("", "", t, u, a.Currency.String(), am, fee)
	if err != nil {
		log.Println(err)
		return
	}

	fv, err := data.NewValue(fmt.Sprintf("%f", fee), a.IsNative())
	if err != nil {
		log.Fatal(err)
	}

	a.Value, err = a.Value.Add(*fv)
	if err != nil {
		log.Fatal(err)
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
	sendResp(resp, c.Ctx)
}
