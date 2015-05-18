//
//
//

package models

import (
	"encoding/base64"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/scrypt"
	"log"
	"time"
)

var Gorm orm.Ormer

func PassHash(password string) string {
	key, err := scrypt.Key([]byte(password), []byte("wangch"), 16384, 8, 1, 32)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(base64.StdEncoding.EncodeToString(key))
}

func init() {
	orm.RegisterDataBase("default", "sqlite3", "icloud.db")
	orm.RegisterModel(new(Recoder), new(Request))
	orm.RegisterModel(new(Role), new(Log))
	orm.RunSyncdb("default", false, true)
	Gorm = orm.NewOrm()
	Gorm.Using("default")
	for _, r := range Gconf.Roles {
		r.Password = PassHash(r.Password)
		_, _, err := Gorm.ReadOrCreate(&r, "Username")
		if err != nil {
			panic(err)
		}
	}
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

// 请求表, 存款和取款
type Request struct {
	Id     int64     // 请求ID, 唯一标识
	CsId   string    // 客服 ID
	CsTime time.Time `orm:"auto_now_add;type(date)"` // 客服提交时间

	CName     string   // 客户真实姓名
	CWallet   string   // 客户钱包地址
	CBankName string   // 客户银行名称
	CBankId   string   // 客户银行账号
	GName     string   // 网关真实姓名
	GWallet   string   // 网关钱包地址
	GBankName string   // 网关银行名称
	GBankId   string   // 网关银行账号
	Currency  string   // 货币 USD,HKD,CNY,BTC等等
	Amount    float64  // 金额
	Fees      float64  // 费用 总金额 = Amount + Fees
	Type      int      // 类别 Issue | Redeem | Deposit | Withdrawal
	R         *Recoder `orm:"rel(one)"`
}

// 状态
const (
	COK = iota // 提交请求 OK
	FOK        // 财务确认 OK
	MOK        // 总监确认 OK
	AOK        // 会计转账 OK

	TimeoutClosed // 超时关闭 - 财务不确认, 总监不确认的情况下都将引发超时关闭
	OKClosed      // OK 关闭
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
	"超时关闭":  TimeoutClosed,
	"完成关闭":  OKClosed,
}

var RStatusMap = map[int]string{
	-1:            "全部",
	COK:           "已提交",
	FOK:           "财务已审批",
	MOK:           "主管已审批",
	AOK:           "转账已完成",
	TimeoutClosed: "超时关闭",
	OKClosed:      "完成关闭",
}

// 记录表, 存款和取款
type Recoder struct {
	Id     int64     // ID, 唯一标识
	FId    string    // 财务 ID
	FTime  time.Time `orm:"auto_now_add;type(date)"` // 财务确认时间
	MId    string    // 总监 ID
	MTime  time.Time `orm:"auto_now_add;type(date)"` // 总监确认时间
	AId    string    // 会计 ID
	ATime  time.Time `orm:"auto_now_add;type(date)"` // 会计转账完成确认时间
	TxHash string    // tx hash
	Status int       // 记录当前状态
	Type   int       // 类别 Issue | Redeem | Deposit | Withdrawal
	R      *Request  `orm:"reverse(one)"`
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

// for test

var testReqs = []*Request{
	&Request{
		CsId:      "wangchC",
		CsTime:    time.Now(),
		CName:     "习近平",
		CWallet:   "iesadfasdfasdfasdfasdfasdf",
		CBankName: "中国工商银行",
		CBankId:   "0000 0000 0000 0000",
		GName:     "王春晖",
		GWallet:   "idsadfasdfasdfasdfasdfasdf",
		GBankName: "中国人民银行",
		GBankId:   "0000 0000 0000 0001",
		Currency:  "USD",
		Amount:    100.00,
		Fees:      1.00,
	},
	&Request{
		CsId:      "wangchC",
		CsTime:    time.Now(),
		CName:     "李克强",
		CWallet:   "ixsadfasdfasdfasdfasdfasdf",
		CBankName: "中国工商银行",
		CBankId:   "1111 1111 1111 1111",
		GName:     "王春晖",
		GWallet:   "idsadfasdfasdfasdfasdfasdf",
		GBankName: "中国人民银行",
		GBankId:   "0000 0000 0000 0001",
		Currency:  "USD",
		Amount:    200.00,
		Fees:      1.00,
	},
	&Request{
		CsId:      "wangchC",
		CsTime:    time.Now(),
		CName:     "王岐山",
		CWallet:   "iwsadfasdfasdfasdfasdfasdf",
		CBankName: "中国工商银行",
		CBankId:   "2222 2222 2222 2222",
		GName:     "王春晖",
		GWallet:   "idsadfasdfasdfasdfasdfasdf",
		GBankName: "中国人民银行",
		GBankId:   "0000 0000 0000 0001",
		Currency:  "USD",
		Amount:    300.00,
		Fees:      1.00,
	},
	&Request{
		CsId:      "wangchC",
		CsTime:    time.Now(),
		CName:     "张德江",
		CWallet:   "izsadfasdfasdfasdfasdfasdf",
		CBankName: "中国工商银行",
		CBankId:   "3333 3333 3333 3333",
		GName:     "王春晖",
		GWallet:   "idsadfasdfasdfasdfasdfasdf",
		GBankName: "中国人民银行",
		GBankId:   "0000 0000 0000 0001",
		Currency:  "USD",
		Amount:    400.00,
		Fees:      1.00,
	},
	&Request{
		CsId:      "wangchC",
		CsTime:    time.Now(),
		CName:     "俞正声",
		CWallet:   "iesadfasdfasdfasdfasdfasdf",
		CBankName: "中国工商银行",
		CBankId:   "4444 4444 4444 4444",
		GName:     "王春晖",
		GWallet:   "idsadfasdfasdfasdfasdfasdf",
		GBankName: "中国人民银行",
		GBankId:   "0000 0000 0000 0001",
		Currency:  "USD",
		Amount:    500.00,
		Fees:      1.00,
	},
	&Request{
		CsId:      "wangchC",
		CsTime:    time.Now(),
		CName:     "刘云山",
		CWallet:   "ilsadfasdfasdfasdfasdfasdf",
		CBankName: "中国工商银行",
		CBankId:   "5555 5555 5555 5555",
		GName:     "王春晖",
		GWallet:   "idsadfasdfasdfasdfasdfasdf",
		GBankName: "中国人民银行",
		GBankId:   "0000 0000 0000 0001",
		Currency:  "USD",
		Amount:    600.00,
		Fees:      1.00,
	},
	&Request{
		CsId:      "wangchC",
		CsTime:    time.Now(),
		CName:     "张高丽",
		CWallet:   "igsadfasdfasdfasdfasdfasdf",
		CBankName: "中国工商银行",
		CBankId:   "6666 6666 6666 6666",
		GName:     "王春晖",
		GWallet:   "idsadfasdfasdfasdfasdfasdf",
		GBankName: "中国人民银行",
		GBankId:   "0000 0000 0000 0001",
		Currency:  "USD",
		Amount:    700.00,
		Fees:      1.00,
	},
}

var testRecs = []*Recoder{
	&Recoder{
		FId:    "wangchF",
		FTime:  time.Now(),
		MId:    "wangchM",
		MTime:  time.Now(),
		AId:    "wangchA",
		ATime:  time.Now(),
		TxHash: "asflasdflasdfalsdfasdfwsd",
		Status: COK,
	},
	&Recoder{
		FId:    "wangchF",
		FTime:  time.Now(),
		MId:    "wangchM",
		MTime:  time.Now(),
		AId:    "wangchA",
		ATime:  time.Now(),
		TxHash: "asflasdflasdfalsdfasdfwsd",
		Status: COK,
	},
	&Recoder{
		FId:    "wangchF",
		FTime:  time.Now(),
		MId:    "wangchM",
		MTime:  time.Now(),
		AId:    "wangchA",
		ATime:  time.Now(),
		TxHash: "asflasdflasdfalsdfasdfwsd",
		Status: COK,
	},
	&Recoder{
		FId:    "wangchF",
		FTime:  time.Now(),
		MId:    "wangchM",
		MTime:  time.Now(),
		AId:    "wangchA",
		ATime:  time.Now(),
		TxHash: "asflasdflasdfalsdfasdfwsd",
		Status: COK,
	},
	&Recoder{
		FId:    "wangchF",
		FTime:  time.Now(),
		MId:    "wangchM",
		MTime:  time.Now(),
		AId:    "wangchA",
		ATime:  time.Now(),
		TxHash: "asflasdflasdfalsdfasdfwsd",
		Status: COK,
	},
	&Recoder{
		FId:    "wangchF",
		FTime:  time.Now(),
		MId:    "wangchM",
		MTime:  time.Now(),
		AId:    "wangchA",
		ATime:  time.Now(),
		TxHash: "asflasdflasdfalsdfasdfwsd",
		Status: COK,
	},
	&Recoder{
		FId:    "wangchF",
		FTime:  time.Now(),
		MId:    "wangchM",
		MTime:  time.Now(),
		AId:    "wangchA",
		ATime:  time.Now(),
		TxHash: "asflasdflasdfalsdfasdfwsd",
		Status: COK,
	},
}
