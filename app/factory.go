package app

import (
	"microloanProject/data"
	"microloanProject/structures"
	"sync"
	"time"
)

type DebtFactory interface {
	CreateDebt(amount float32, percent float32) structures.Debt
}

type CreditDebtFactory struct{}

func (f *CreditDebtFactory) CreateDebt(amount int, percent float32, months int) structures.Debt {
	// some logic to calculate
	a := GetAutoInc()
	return &structures.CreditDebt{a.GenerateCreditId(),
		percent,
		amount,
		amount,
		time.Now(),
		time.Now().Add(time.Second * 5),
	}
}

type MortgageDebtFactory struct{}

func (f *MortgageDebtFactory) CreateDebt(amount int, percent float32, months int, address string) structures.Debt {
	// some logic to calculate
	a := GetAutoInc()
	return &structures.MortgageDebt{a.GenerateMortgageId(),
		address,
		percent,
		amount,
		amount,
		time.Now(),
		time.Now().Add(time.Minute * 5),
	}
}

type AutoInc struct {
	userId     int
	creditId   int
	mortgageId int
}

var autoIncOnce sync.Once
var autoIncInstance *AutoInc

func GetAutoInc() *AutoInc {
	autoIncOnce.Do(func() {
		db := data.GetDatabase()
		defer db.CloseConnection()
		autoIncInstance = &AutoInc{db.GetMaxUserId(), db.GetMaxCreditId(), db.GetMaxMortgageId()}
	})
	return autoIncInstance
}

func (a *AutoInc) GenerateUserId() int {
	a.userId++
	return a.userId
}

func (a *AutoInc) GenerateCreditId() int {
	a.creditId++
	return a.creditId
}

func (a *AutoInc) GenerateMortgageId() int {
	a.mortgageId++
	return a.mortgageId
}
