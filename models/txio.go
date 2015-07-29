// Copyright 2015 iCloudFund. All Rights Reserved.

package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/wangch/glog"
	// "github.com/wangch/ripple/crypto"
	"github.com/wangch/ripple/data"
	"github.com/wangch/ripple/websockets"
)

// 监控网关账号的存款deposit和ICC发行issue 取款withdrawal和ICC赎回redeem
var gws *websockets.Remote

func monitor(serverAddr string, wallets []string) error {
	for {
		ws, err := websockets.NewRemote(serverAddr)
		gws = ws
		if err != nil {
			glog.Error(err)
			time.Sleep(time.Second * 1)
			continue
		}

		_, err = ws.Subscribe(false, false, false, false, wallets)
		if err != nil {
			glog.Error(err)
			time.Sleep(time.Second * 1)
			continue
		}

		for {
			msg, ok := <-ws.Incoming
			if !ok {
				glog.Warning("ws.Incoming chan closed.")
				break
			}

			switch msg := msg.(type) {
			case *websockets.TransactionStreamMsg:
				// the transaction must be validated
				// and only watch payments
				b, err := json.MarshalIndent(msg, "", "  ")
				if err != nil {
					glog.Error(err)
				}
				glog.Info(string(b))
				if !msg.EngineResult.Success() {
					glog.Warning("the transaction NOT success")
					break
				}

				if msg.Transaction.GetType() == "Payment" {
					paymentTx := msg.Transaction.Transaction.(*data.Payment)
					out := isOut(paymentTx.Account.String(), wallets)
					if paymentTx.InvoiceID == nil {
						glog.Warning("paymentTx.InvoiceID == nil")
						break
					}
					// query the paymen tx InvoiceId in database and update tx hash
					invid := paymentTx.InvoiceID.String()
					r := &Request{}
					Gorm.QueryTable("request").Filter("invoice_id", invid).RelatedSel().One(r)
					if r.R == nil {
						glog.Warning("the payment invoiceID " + invid + "is NOT in database")
						// must be cold wallet send to hotwallet
						break
					}

					r.R.TxHash = paymentTx.Hash.String()
					if out { // 存款 or 发行ICC
						r.R.Status = OKC
					} else { // 取款 or 回收ICC
						r.R.Status = COK
					}
					_, err = Gorm.Update(r.R)
					if err != nil {
						// have error in database
						// must report the error msg on web
						glog.Error(err)
					}
				}
			}
		}
	}
}

// acc是否向外支付
func isOut(acc string, accs []string) bool {
	for _, a := range accs {
		if a == acc {
			return true
		}
	}
	return false
}

// for _, w := range Gconf.HoltWallet {
// 	a, err := data.NewAccountFromAddress(w.AccountId)
// 	if err != nil {
// 		return err
// 	}
// 	r, err := gws.AccountInfo(a)
// 	if err != nil {
// 		return err
// 	}
// 	r.AccountData.ba
// }

// TODO: 获取Gconf中每个Hotwallet的每种货币的余额, 选择最多的那个作为sender
func Payment(r *Request, sender string) error {
	secret := ""
	for _, x := range Gconf.HoltWallet {
		if sender == x.AccountId {
			secret = x.Secret
			break
		}
	}
	if secret == "" {
		errMsg := fmt.Sprintf("Payment error: the Sender %s is NOT in config hotwallets", sender)
		err := errors.New(errMsg)
		glog.Error(err)
		return err
	}
	return payment(gws, secret, sender, Gconf.ColdWallet, r.UWallet, r.Currency, r.InvoiceId, r.Amount)
}

// func payment(ws *websockets.Remote, secret, sender, issuer, recipient, currency, invoiceID string, amount float64) error {

// secret is sender's secret
func payment(ws *websockets.Remote, secret, sender, issuer, recipient, currency, invoiceID string, amount float64) error {
	glog.Info("payment:", secret, sender, recipient, currency, amount, issuer, invoiceID)
	sam := ""
	if currency == "ICC" {
		sam = fmt.Sprintf("%d/ICC", uint64(amount))
	} else {
		sam = fmt.Sprintf("%f/%s/%s", amount, currency, issuer)
	}
	a, err := data.NewAmount(sam)
	if err != nil {
		err = errors.New("NewAmount error: " + sam + err.Error())
		glog.Error(err)
		return err
	}

	ptx := &websockets.PaymentTx{
		TransactionType: "Payment",
		Account:         sender,
		Destination:     recipient,
		Amount:          a,
		InvoiceID:       invoiceID,
	}

	// glog.Infof("payment: %+v", ptx)

	r, err := ws.SubmitWithSign(ptx, secret)
	if err != nil {
		glog.Error(err)
		return err
	}
	glog.Infof("pament result: %+v", r)
	if !r.EngineResult.Success() {
		return errors.New(r.EngineResultMessage)
	}
	return nil
}
