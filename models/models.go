// Copyright 2015 iCloudFund. All Rights Reserved.

package models

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wangch/glog"
	"github.com/wangch/ripple/data"
	"golang.org/x/crypto/scrypt"
)

var Gorm orm.Ormer

func PassHash(password string) string {
	key, err := scrypt.Key([]byte(password), []byte("wangch"), 16384, 8, 1, 32)
	if err != nil {
		glog.Error(err)
		return ""
	}
	return string(base64.StdEncoding.EncodeToString(key))
}

func Init() {
	orm.RegisterDataBase("default", "sqlite3", "icloud.db")
	orm.RegisterModel(new(Recoder), new(Request))
	orm.RegisterModel(new(Role), new(Log))
	orm.RunSyncdb("default", false, true)
	Gorm = orm.NewOrm()
	Gorm.Using("default")
	initConf()
	for _, r := range Gconf.Roles {
		r.Password = PassHash(r.Password)
		_, _, err := Gorm.ReadOrCreate(&r, "Username")
		if err != nil {
			glog.Fatal(err)
		}
	}
	accs := make([]string, len(Gconf.HoltWallet))
	for i, hw := range Gconf.HoltWallet {
		accs[i] = hw.AccountId
	}
	accs = append(accs, Gconf.ColdWallet)
	go monitor(Gconf.ServerAddr, accs)
}

// 人员类别
const (
	RoleC int = iota // 客服
	RoleF            // 财务
	RoleM            // 总监
	RoleA            // 会计
)

// 人员表
type Role struct {
	Id       int64
	Username string `orm:"unique"`
	Password string
	Type     int
}

const (
	Issue = iota
	Redeem
	Deposit
	Withdrawal
)

type User struct {
	UName        string // 客户真实姓名
	UWallet      string // 客户钱包地址
	UBankName    string // 客户银行名称
	UBankId      string // 客户银行账号
	UContact     string // 客户联系方式
	UCertificate string // 客户存款凭证 文件的url
}

// 请求表, 存款和取款
type Request struct {
	Id int64 // 请求ID, 唯一标识
	// CsId   string    // 客服 ID
	CTime time.Time `orm:"auto_now_add;type(date)"` // 客服提交时间

	UName        string   // 客户真实姓名
	UWallet      string   // 客户钱包地址
	UBankName    string   // 客户银行名称
	UBankId      string   // 客户银行账号
	UContact     string   // 客户联系方式
	UCertificate string   // 客户存款凭证 文件的url
	Currency     string   // 货币 USD,HKD,CNY,BTC等等
	Amount       float64  // 金额
	Fees         float64  // 费用 总金额 = Amount + Fees
	Type         int      // 类别 Issue | Redeem | Deposit | Withdrawal
	InvoiceId    string   // 标识此次请求
	R            *Recoder `orm:"rel(one)"`
}

// 状态
const (
	COK = iota + 1 // 提交请求 OK
	FOK            // 财务确认 OK
	MOK            // 总监确认 OK
	AOK            // 会计转账 OK
	TOC            // 超时关闭 - 财务不确认, 总监不确认的情况下都将引发超时关闭
	OKC            // OK 关闭
)

var StatusSlice = []string{
	"全部",
	"已提交",
	"财务已审批",
	"主管已审批",
	"转账已完成",
	"超时关闭",
	"完成关闭",
}

var StatusMap = map[string]int{
	"全部":    -1,
	"已提交":   COK,
	"财务已审批": FOK,
	"主管已审批": MOK,
	"转账已完成": AOK,
	"超时关闭":  TOC,
	"完成关闭":  OKC,
}

var RStatusMap = map[int]string{
	-1:  "全部",
	COK: "已提交",
	FOK: "财务已审批",
	MOK: "主管已审批",
	AOK: "转账已完成",
	TOC: "超时关闭",
	OKC: "完成关闭",
}

// 记录表, 存款和取款
type Recoder struct {
	Id    int64     // ID, 唯一标识
	FId   string    // 财务 ID
	FTime time.Time `orm:"auto_now_add;type(date)"` // 财务确认时间
	// GName     string    // 网关真实姓名
	GWallet string // 网关钱包地址
	// GBankName string    // 网关银行名称
	GBankId string    // 网关银行账号
	MId     string    // 总监 ID
	MTime   time.Time `orm:"auto_now_add;type(date)"` // 总监确认时间
	AId     string    // 会计 ID
	ATime   time.Time `orm:"auto_now_add;type(date)"` // 会计转账完成确认时间
	Status  int       // 记录当前状态
	// Type   int       // 类别 Issue | Redeem | Deposit | Withdrawal
	TxHash string   // tx hash
	R      *Request `orm:"reverse(one)"`
}

// 操作类别
const (
	RequestCommit = iota
	FinCommit
	MasterCommit
)

// 日志表
type Log struct {
	Id          int64
	Name        string    // 名称 Id
	LoginTime   time.Time // 登录时间
	OprateType  int       // 操作类别
	OpratorTime time.Time // 操作时间
	LogoutTime  time.Time // 登出时间
}

func AddReq(gid, gwallet string, t int, u *User, currency string, amount, fees float64) (*Request, error) {
	gbas := Gconf.GBAs
	var gba *GateBankAccount
	for _, g := range gbas {
		if strings.Contains(gid, g.BankId) {
			gba = &g
			break
		}
	}

	req := &Request{
		CTime:        time.Now(),
		UName:        u.UName,
		UWallet:      u.UWallet,
		UBankName:    u.UBankName,
		UBankId:      u.UBankId,
		UContact:     u.UContact,
		UCertificate: u.UCertificate,
		Currency:     currency,
		Amount:       amount,
		Fees:         fees,
		Type:         t,
	}
	req.InvoiceId = getInvoiceId(req)

	if gba == nil && len(gbas) > 0 {
		gba = &gbas[0]
	}

	tm := time.Unix(0, 0) //time.Date(1979, 1, 1, 0, 0, 0, 0, time.UTC)

	rec := &Recoder{
		R:       req,
		FTime:   tm,
		MTime:   tm,
		ATime:   tm,
		GWallet: gwallet,
	}
	// 因为取款和回收会产生多个记录, 所以只有交易完成时, 才变为COK(提交状态)
	// 所以这里只设置存款和发行的COK状态
	if t == Deposit || t == Issue {
		rec.Status = COK
	}
	if gba != nil {
		// req.GName = gba.Name
		// req.GBankName = gba.BankName
		rec.GBankId = gba.BankId
	}
	// log.Printf("%#v", rec)
	req.R = rec
	_, err := Gorm.Insert(rec)
	if err != nil {
		glog.Fatal(err)
	}
	_, err = Gorm.Insert(req)
	if err != nil {
		glog.Fatal(err)
	}
	return req, nil
}

func getInvoiceId(r *Request) string {
	s := fmt.Sprintf("%s%s%s%s%s%s%s%s%f%f%d", r.CTime.String(), r.UName, r.UWallet, r.UBankName, r.UBankId, r.UContact, r.Currency, r.Amount, r.Fees, r.Type)
	hash := sha256.Sum256([]byte(s))
	h, err := data.NewHash256(hash[:])
	if err != nil {
		glog.Error(err)
		return ""
	}
	return h.String()
}

// 获取发行的ICC
func GetIssueIccs() int64 {
	return 0
}
