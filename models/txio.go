// Copyright 2015 iCloudFund. All Rights Reserved.

package models

import (
	"crypto/sha256"
	"encoding/hex"
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
					err = Gorm.Read(r)
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
func Payment(r *Request, secret, sender string) error {
	return payment(secret, sender, r.UWallet, r.Currency, r.Amount)
}

// secret is sender's secret
func payment(secret, sender, recipient, currency string, amount float64) error {
	srcAcccount, err := data.NewAccountFromAddress(sender)
	if err != nil {
		return err
	}
	ar, err := gws.AccountInfo(*srcAcccount)
	if err != nil {
		return err
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
		return err
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
		sam = fmt.Sprintf("%u/ICC", uint64(amount*1e6))
	} else {
		sam = fmt.Sprintf("%f/%s/%s", amount, currency, Gconf.ColdWallet)
	}
	a, err := data.NewAmount(sam)
	if err != nil {
		return err
	}

	h, err := data.NewHash256(GetInvoiceID(a))
	if err != nil {
		return err
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
	key, err := crypto.NewEd25519Key(seed.Payload())
	if err != nil {
		return err
	}
	// key, err := crypto.NewECDSAKey(seed)
	// if err != nil {
	// 	return err
	// }

	keySeq := uint32(0)
	err = data.Sign(ptx, key, &keySeq)
	if err != nil {
		return err
	}
	r, err := gws.Submit(ptx)
	if err != nil {
		return err
	}
	log.Println(r)
	return nil
}

func GetInvoiceID(a *data.Amount) string {
	hash := sha256.Sum256(a.Bytes())
	invoiceID := hex.EncodeToString(hash[:])
	return invoiceID
}
