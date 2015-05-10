//
//
//

package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

var Gorm orm.Ormer

func init() {
	// orm.RegisterDataBase("default", "mysql", "root:root@/icloud?charset=utf8")
	orm.RegisterDataBase("default", "sqlite3", "icloud.db")
	// orm.RegisterModel(new(DepositRequest), new(WithdrawalRequest), new(IssueRequest), new(RedeemRequest))
	// orm.RegisterModel(new(DepositRecoder), new(WithdrawalRecoder), new(IssueRecoder), new(RedeemRecoder))
	orm.RegisterModel(new(Role), new(Log))
	orm.RegisterModel(new(Recoder), new(Request))
	Gorm = orm.NewOrm()
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
	Id           int64
	Username     string
	PasswordHash string
	Type         int
}

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
	Currence  string   // 货币 USD,HKD,CNY,BTC等等
	Amount    float64  // 金额
	Fees      float64  // 费用 总金额 = Amount + Fees
	R         *Recoder `orm:"reverse(one)"`
}

type DepositRequest Request

func (d *DepositRequest) TableName() string {
	return "deposit_req"
}

type WithdrawalRequest Request

func (d *WithdrawalRequest) TableName() string {
	return "withdrawal_req"
}

type IssueRequest Request

func (d *IssueRequest) TableName() string {
	return "issue_req"
}

type RedeemRequest Request

func (d *RedeemRequest) TableName() string {
	return "redeem_req"
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

var StatusMap = map[string]int{
	"已提交":   COK,
	"财务已审批": FOK,
	"主管已审批": MOK,
	"转账已完成": AOK,
	"超时关闭":  TimeoutClosed,
	"完成关闭":  OKClosed,
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
	R      *Request  `orm:"rel(one)"`
}

type DepositRecoder Recoder

func (d *DepositRecoder) TableName() string {
	return "deposit_rec"
}

type WithdrawalRecoder Recoder

func (w *WithdrawalRecoder) TableName() string {
	return "withdrawal_rec"
}

type IssueRecoder Recoder

func (i *IssueRecoder) TableName() string {
	return "issue_rec"
}

type RedeemRecoder Recoder

func (r *RedeemRecoder) TableName() string {
	return "redeem_rec"
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
