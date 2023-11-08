package app

import (
	"fmt"
	"microloanProject/data"
	"microloanProject/structures"
	"strconv"
)

var CurrentUser *structures.User
var creditFactory CreditDebtFactory
var mortgageFactory MortgageDebtFactory

func RegisterUser(username, password, email, phone string) *structures.User {
	newUser := CreateUser(username, password)

	newUser.AddEmailContact(email)
	newUser.AddPhoneContact(phone)

	db := data.GetDatabase()
	defer db.CloseConnection()
	db.AddUserToDB(newUser)
	fmt.Println("Пользователь успешно зарегистрирован.")
	CurrentUser = newUser
	return newUser
}

func Login(username, password string) bool {
	db := data.GetDatabase()
	defer db.CloseConnection()
	for _, user := range db.GetUsers() {
		if user.Username == username && user.Password == password {
			CurrentUser = user
			return true
		}
	}
	return false
}

func CreateUser(username, password string) *structures.User {
	a := GetAutoInc()
	return &structures.User{Id: a.GenerateUserId(), Username: username, Password: password}
}

func AddCredit() structures.Debt {
	fmt.Println("Кредит на какую сумму вы хотите взять?")
	var str string
	fmt.Scanln(&str)
	amount, _ := strconv.ParseFloat(str, 32)
	fmt.Println("Под какой процент кредит вы хотите взять?")
	fmt.Scanln(&str)
	percent, _ := strconv.ParseFloat(str, 32)
	fmt.Println("На какой период вы хотите взять кредит(В месяцах)?")
	fmt.Scanln(&str)
	months, _ := strconv.ParseInt(str, 10, 32)
	credit := creditFactory.CreateDebt(
		int(amount),
		float32(percent),
		int(months),
	)
	CurrentUser.AddDebt(credit)
	db := data.GetDatabase()
	defer db.CloseConnection()
	db.AddDebtToUser(CurrentUser, credit)
	fmt.Println("Вы успешно взяли кредит")
	return credit
}

func AddMortgage() structures.Debt {
	fmt.Println("Ипотеку на какую сумму вы хотите взять?")
	var str string
	fmt.Scanln(&str)
	amount, _ := strconv.ParseFloat(str, 32)
	fmt.Println("Под какой процент ипотеку вы хотите взять?")
	fmt.Scanln(&str)
	percent, _ := strconv.ParseFloat(str, 32)
	fmt.Println("На какой период вы хотите взять ипотеку(В месяцах)?")
	fmt.Scanln(&str)
	months, _ := strconv.ParseInt(str, 10, 32)
	fmt.Println("Назовите адрес дома который вы хотите приобрести:")
	fmt.Scanln(&str)

	mortage := mortgageFactory.CreateDebt(
		int(amount),
		float32(percent),
		int(months),
		str,
	)
	CurrentUser.AddDebt(mortage)
	db := data.GetDatabase()
	defer db.CloseConnection()
	db.AddDebtToUser(CurrentUser, mortage)
	fmt.Println("Вы успешно взяли ипотеку")
	return mortage
}

func Pay(debt structures.Debt) {
	var strategy PaymentStrategy

	fmt.Println("1. Оплатить через KaspiQR")
	fmt.Println("2. Оплатить картой")
	var choice string
	fmt.Scanln(&choice)
	if choice == "1" {
		strategy = &KaspiQRStrategy{}
	} else if choice == "2" {
		strategy = &CardPaymentStrategy{}
	}
	strategy.Pay(debt)
}
