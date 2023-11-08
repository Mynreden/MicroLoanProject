package main

import (
	"fmt"
	"microloanProject/app"
	"microloanProject/structures"
	"os"
	"strconv"
	"time"
)

func main() {
	var isLogged bool

	for !isLogged {
		fmt.Println("1. Зарегистрировать нового пользователя")
		fmt.Println("2. Войти")
		fmt.Println("3. Выйти")
		var choice string
		fmt.Print("Выберите действие: ")
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			var username, password, email, phone string
			fmt.Print("Введите имя пользователя: ")
			fmt.Scanln(&username)
			fmt.Print("Введите пароль: ")
			fmt.Scanln(&password)
			fmt.Print("Введите почту: ")
			fmt.Scanln(&email)
			fmt.Print("Введите номер телефона: ")
			fmt.Scanln(&phone)
			app.RegisterUser(username, password, email, phone)
			isLogged = true
		case "2":
			var username, password string
			fmt.Print("Введите имя пользователя: ")
			fmt.Scanln(&username)
			fmt.Print("Введите пароль: ")
			fmt.Scanln(&password)
			if app.Login(username, password) {
				fmt.Println("Вход выполнен успешно.")
				isLogged = true
			} else {
				fmt.Println("Неверное имя пользователя или пароль.")
			}
		case "3":
			os.Exit(0)
		default:
			fmt.Println("Некорректный выбор. Попробуйте снова.")
		}
	}

	fmt.Println("You logged as: ", app.CurrentUser.Username)
	for _, debt := range app.CurrentUser.Debts {
		go waitingForDebt(debt)
	}

	for true {
		fmt.Println("1. Взять займ")
		fmt.Println("2. Посмотреть мои займы")
		fmt.Println("3. Оплатить долг")
		fmt.Println("4. Выйти из аккаунта")
		var choice string
		fmt.Print("Выберите действие: ")
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			fmt.Println("1. Взять кредит")
			fmt.Println("2. Взять ипотеку")
			fmt.Scanln(&choice)
			var debt structures.Debt
			switch choice {
			case "1":
				debt = app.AddCredit()
			case "2":
				debt = app.AddMortgage()
			}
			go waitingForDebt(debt)
		case "2":
			renderDepts()
		case "3":
			renderDepts()
			fmt.Println("Выберите номер кредита который желаете оплатить?")
			var choice = ""
			fmt.Scanln(&choice)
			num, _ := strconv.ParseInt(choice, 10, 32)
			app.Pay(app.CurrentUser.Debts[num-1])
		case "4":
			os.Exit(0)
		default:
			fmt.Println("Некорректный выбор. Попробуйте снова.")
		}
	}
}

func waitingForDebt(debt structures.Debt) {
	if time.Now().Compare(debt.GetNextPaymentDate()) == 1 {
		return
	}

	duration := debt.GetNextPaymentDate().Sub(time.Now())
	timer := time.NewTimer(duration)
	select {
	case <-timer.C:
		NotifyUser(debt)
	}
}

func NotifyUser(debt structures.Debt) {
	fmt.Println("Начато оповещение о необходимости произвести платеж")
	app.CurrentUser.Notify()
}

func renderDepts() {
	for i, debt := range app.CurrentUser.Debts {
		fmt.Printf("%d. %s\n",
			i+1,
			debt.ToString(),
		)
	}
	fmt.Print("\n")
}
