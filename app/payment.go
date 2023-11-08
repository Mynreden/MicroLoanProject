package app

import (
	"fmt"
	"microloanProject/data"
	"microloanProject/structures"
	"time"
)

type PaymentStrategy interface {
	Pay(debt structures.Debt)
}

type KaspiQRStrategy struct {
}

func (s *KaspiQRStrategy) Pay(debt structures.Debt) {
	deleteDebt(debt)
	time.Sleep(time.Second * 2)
	fmt.Println("Оплата прошла успешно")
	// some KaspiQR logic
}

type CardPaymentStrategy struct {
}

func (s *CardPaymentStrategy) Pay(debt structures.Debt) {
	deleteDebt(debt)
	fmt.Println("Оплата картой прошла успешно")

	// some card payment logic
}

func deleteDebt(debt structures.Debt) {
	db := data.GetDatabase()
	defer db.CloseConnection()
	CurrentUser.PayDebt(debt)
	db.DeleteDebt(debt)
}
