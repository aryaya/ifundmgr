//
//
//

package models

import (
	"log"
)

const (
	RoleTable       = "role"
	DepositTable    = "deposit"
	WithdrawalTable = "withdrawal"
	IssuseTable     = "issuse"
	RedeemTable     = "redeem"
)

type DataBase struct {
}

func (db *DataBase) CreateTable(tableName string) error {
	switch tableName {
	case RoleTable:
	case DepositTable:
	case WithdrawalTable:
	case IssuseTable:
	case RedeemTable:
	default:
		panic("can't go here")
	}
	return nil
}

func (db *DataBase) Commite(tableName, csid string, r *Request) error {
	return nil
}

func (db *DataBase) FinVerificate(tableName string, id int64) error {
	return nil
}

func (db *DataBase) MasterVerficate(tableName string, id int64) error {
	return nil
}

func (db *DataBase) CsQuery(id, where string) ([]*Request, error) {
	return nil, nil
}

func (db *DataBase) FinQuery(id, where string) ([]*Recoder, error) {
	return nil, nil
}

func (db *DataBase) MasterQuery(id, where string) ([]*Recoder, error) {
	return nil, nil
}
