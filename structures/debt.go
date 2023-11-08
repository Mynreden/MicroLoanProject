package structures

import (
	"fmt"
	"time"
)

type Debt interface {
	GetNextPaymentDate() time.Time
	ToString() string
}

type CreditDebt struct {
	Id              int       `json:"id"`
	Percent         float32   `json:"percent"`
	InitialAmount   int       `json:"initialAmount"`
	Remainder       int       `json:"remainder"`
	StartDate       time.Time `json:"startDate"`
	NextPaymentDate time.Time `json:"nextPaymentDate"`
}

func (debt *CreditDebt) ToString() string {
	return fmt.Sprintf("Credit: Percent: %.2f;   Remainder: %d;   Next Payment: %s",
		debt.Percent,
		debt.Remainder,
		debt.NextPaymentDate.Format("2006-01-02 15:04:05"),
	)
}

func (d *CreditDebt) GetNextPaymentDate() time.Time {
	return d.NextPaymentDate
}

type MortgageDebt struct {
	Id              int       `json:"id"`
	Address         string    `json:"address"`
	Percent         float32   `json:"percent"`
	InitialAmount   int       `json:"initialAmount"`
	Remainder       int       `json:"remainder"`
	StartDate       time.Time `json:"dateOfCollection"`
	NextPaymentDate time.Time `json:"nextPaymentDate"`
}

func (debt *MortgageDebt) ToString() string {
	return fmt.Sprintf("Mortgage: Percent: %.2f;   Remainder: %d;   Next Payment: %s",
		debt.Percent,
		debt.Remainder,
		debt.NextPaymentDate.Format("2006-01-02 15:04:05"),
	)
}

func (d *MortgageDebt) GetNextPaymentDate() time.Time {
	return d.NextPaymentDate
}
