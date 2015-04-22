//
//
//

package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(DepositRequest), new(WithdrawalRequest), new(IssueRequest), new(RedeemRequest))
	orm.RegisterModel(new(DepositRecoder), new(WithdrawalRecoder), new(IssueRecoder), new(RedeemRecoder))
	orm.RegisterModel(new(Role), new(Log))
}

// 人员类别
const (
	Cs     int = iota // 客服
	Fin               // 财务
	Master            // 总监
)

// 人员表
type Role struct {
	Id           string
	PasswordHash string
	Type         int
}

// 请求表, 存款和取款
type Request struct {
	Id     int64     // 请求ID, 唯一标识
	CsId   string    // 客服 ID
	CsTime time.Time // 客服提交时间

	Name     string  // 真实姓名
	BankName string  // 银行名称
	BankId   string  // 银行账号
	Currence string  // 货币 USD,HKD,CNY,BTC等等
	Amount   float64 // 金额
	Fees     float64 // 费用  总金额 = Amount + Fees
	Wallet   string  // 钱包地址
}

type DepositRequest Request
type WithdrawalRequest Request
type IssueRequest Request
type RedeemRequest Request

// 记录状态
const (
	Commited = iota // 提交请求
	FinOK           // 财务确认 OK
	MasterOK        // 总监确认 OK	Id       int64     // 请求ID, 唯一标识

	TimeoutClosed // 超时关闭 - 财务不确认, 总监不确认的情况下都将引发超时关闭
	OKClosed      // OK 关闭
)

// 记录表, 存款和取款
type Recoder struct {
	Id         int64     // 请求ID, 唯一标识
	FinId      string    // 财务 ID
	FinTime    time.Time // 财务确认时间
	MasterId   string    // 总监 ID
	MasterTime time.Time // 总监确认时间
	TxHash     string    // tx hash
	Status     int       // 记录当前状态
	R          *Request  // `orm:"rel(one)"`
}

type DepositRecoder Recoder
type WithdrawalRecoder Recoder
type IssueRecoder Recoder
type RedeemRecoder Recoder

// 操作类别
const (
	RequestCommit = iota
	FinCommit
	MasterCommit
)

// 日志表
type Log struct {
	Name        string    // 名称 Id
	LoginTime   time.Time // 登录时间
	OprateType  int       // 操作类别
	OpratorTime time.Time // 操作时间
	LogoutTime  time.Time // 登出时间
}
