// Copyright 2015 iCloudFund. All Rights Reserved.

package models

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/wangch/ripple/crypto"
	"github.com/wangch/ripple/data"
	"github.com/wangch/ripple/websockets"
)

var gws *websockets.Remote

// 监控网关账号的存款deposit和ICC发行issue 取款withdrawal和ICC赎回redeem

func monitor(serverAddr string, wallets []string) error {
	for {
		ws, err := websockets.NewRemote(serverAddr)
		gws = ws
		if err != nil {
			log.Println("@@@ 0:", err)
			time.Sleep(time.Second * 1)
			continue
		}

		_, err = ws.Subscribe(false, false, false, false, wallets)
		if err != nil {
			log.Println("@@@ 1:", err)
			time.Sleep(time.Second * 1)
			continue
		}

		for {
			msg, ok := <-ws.Incoming
			if !ok {
				log.Println("@@@ 2:", "ws.Incoming closed")
				break
			}

			switch msg := msg.(type) {
			case *websockets.TransactionStreamMsg:
				// the transaction must be validated
				// and only watch payments
				if msg.Transaction.GetType() == "Payment" {
					paymentTx := msg.Transaction.Transaction.(*data.Payment)
					log.Println(paymentTx)
					if paymentTx.InvoiceID == nil {
						break
					}
					// query the paymen tx InvoiceId in database and update tx hansh
					r := &Request{InvoiceId: paymentTx.InvoiceID.String()}
					err = Gorm.Read(r, "invoice_id")
					if err != nil {
						// have error in database
						// must report the error msg on web
						log.Println("@@@ 3:", err)
						break
					}
					r.R.TxHash = paymentTx.Hash.String()
					if isOut(paymentTx.Account.String(), wallets) {
						r.R.Status = OKC
					}
					_, err = Gorm.Update(r.R)
					if err != nil {
						// have error in database
						// must report the error msg on web
						log.Println("@@@ 4:", err)
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
		return errors.New(errMsg)
	}
	return payment(gws, secret, sender, Gconf.ColdWallet, r.UWallet, r.Currency, r.InvoiceId, r.Amount)
}

// secret is sender's secret
func payment(ws *websockets.Remote, secret, sender, issuer, recipient, currency, invoiceID string, amount float64) error {
	log.Println("payment:", sender, secret, issuer, recipient, currency, amount)
	srcAcccount, err := data.NewAccountFromAddress(sender)
	if err != nil {
		return errors.New("NewAccountFromAddress error: " + err.Error())
	}
	ar, err := ws.AccountInfo(*srcAcccount)
	if err != nil {
		return errors.New("AccountInfo error: " + err.Error())
	}
	destAccount, err := data.NewAccountFromAddress(recipient)
	if err != nil {
		return err
	}

	lls := ar.LedgerSequence
	if lls == 0 {
		lls = ar.AccountData.Ledger()
	}
	if lls == 0 {
		return errors.New("last ledger sequence can't get")
	}
	lls += 4

	fee, err := data.NewNativeValue(int64(100))
	if err != nil {
		return errors.New("NewNativeValue error: " + err.Error())
	}

	tb := data.TxBase{
		TransactionType:    data.PAYMENT,
		Account:            *srcAcccount,
		Sequence:           *ar.AccountData.Sequence,
		LastLedgerSequence: &lls,
		Fee:                *fee,
	}

	sam := ""
	if currency == "ICC" {
		sam = fmt.Sprintf("%d/ICC", uint64(amount*1e6))
	} else {
		sam = fmt.Sprintf("%f/%s/%s", amount, currency, issuer)
	}
	a, err := data.NewAmount(sam)
	if err != nil {
		return errors.New("NewAmount error: " + sam + err.Error())
	}

	h, err := data.NewHash256(invoiceID)
	if err != nil {
		return errors.New("NewHash256 error: " + err.Error())
	}

	ptx := &data.Payment{
		TxBase:      tb,
		Destination: *destAccount,
		Amount:      *a,
		InvoiceID:   h,
	}

	seed, err := crypto.NewRippleHashCheck(secret, crypto.RIPPLE_FAMILY_SEED)
	if err != nil {
		return err
	}

	// Ed25519 NOT surport because client use ECDSA
	// key, err := crypto.NewEd25519Key(seed.Payload())
	// if err != nil {
	// 	return err
	// }

	key, err := crypto.NewECDSAKey(seed.Payload())
	if err != nil {
		return err
	}
	var pseq uint32 = 0

	err = data.Sign(ptx, key, &pseq)
	if err != nil {
		return err
	}
	r, err := ws.Submit(ptx)
	if err != nil {
		return err
	}
	log.Println(r)
	return nil
}