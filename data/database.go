package data

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"microloanProject/structures"
	"sync"
	"time"
)

const (
	host       = "localhost"
	port       = 5432
	dbUser     = "postgres"
	dbPassword = "sultan2004"
	dbname     = "SDPendterm"
)

type Database struct {
	conn   *sql.DB
	mu     sync.Mutex
	config string
}

var DatabaseInstance *Database
var DatabaseOnce sync.Once

func GetDatabase() *Database {
	DatabaseOnce.Do(func() {
		config := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, dbUser, dbPassword, dbname)
		DatabaseInstance = &Database{config: config}
	})
	DatabaseInstance.OpenConnection()
	return DatabaseInstance
}

func (db *Database) OpenConnection() {
	a, err := sql.Open("postgres", db.config)
	if err != nil {
		panic(err)
	}
	db.conn = a
}

func (db *Database) CloseConnection() {
	err := db.conn.Close()
	if err != nil {
		panic(err)
	}
}

func (db *Database) GetUsers() []*structures.User {
	var arr []*structures.User
	rows, _ := db.conn.Query("SELECT * FROM users")
	var username, password string
	var id int
	for rows.Next() {
		err := rows.Scan(&id, &username, &password)
		if err != nil {
			panic(err)
		}
		user := &structures.User{Id: id, Username: username, Password: password}
		user.Debts = db.getDebtsByUserId(id)
		arr = append(arr, user)
	}
	return arr
}

func (db *Database) AddUserToDB(newUser *structures.User) {
	_, err := db.conn.Exec("INSERT INTO users (id, username, password) VALUES ($1, $2, $3)", newUser.Id, newUser.Username, newUser.Password)
	if err != nil {
		panic(err)
	}
}

func (db *Database) GetUserById(id int) *structures.User {
	var username, password string
	rows, err := db.conn.Query("SELECT * FROM users where id=$1", id)
	if err != nil {
		panic(err)
	}

	err = rows.Scan(&id, &username, &password)
	if err != nil {
		panic(err)
	}

	user := &structures.User{Id: id, Username: username, Password: password}
	user.Debts = db.getDebtsByUser(user)
	return user
}

func (db *Database) getDebtsByUser(user *structures.User) []structures.Debt {
	return db.getDebtsByUserId(user.Id)
}

func (db *Database) getDebtsByUserId(userid int) []structures.Debt {
	var arr []structures.Debt
	rows, _ := db.conn.Query("SELECT * FROM credit_debt WHERE user_id=$1", userid)
	var id, initialAmount, remainder, userId int
	var percent float32
	var startDate, nextPaymentDate time.Time
	for rows.Next() {
		err := rows.Scan(&id, &percent, &initialAmount, &remainder, &startDate, &nextPaymentDate, &userId)
		if err != nil {
			panic(err)
		}
		arr = append(arr, &structures.CreditDebt{Id: id, Percent: percent, InitialAmount: initialAmount, Remainder: remainder, StartDate: startDate, NextPaymentDate: nextPaymentDate})
	}
	rows, _ = db.conn.Query("SELECT * FROM mortgage_debt WHERE user_id=$1", userid)
	var address string
	for rows.Next() {
		err := rows.Scan(&id, &percent, &initialAmount, &remainder, &startDate, &nextPaymentDate, &userId, &address)
		if err != nil {
			panic(err)
		}
		arr = append(arr, &structures.MortgageDebt{Id: id, Percent: percent, InitialAmount: initialAmount, Remainder: remainder, StartDate: startDate, NextPaymentDate: nextPaymentDate, Address: address})
	}
	return arr
}

func (db *Database) DeleteUser(User *structures.User) {
	db.DeleteUserById(User.Id)
}

func (db *Database) DeleteUserById(id int) {
	_, err := db.conn.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		panic(err)
	}
}

func (db *Database) DeleteDebt(debt structures.Debt) {
	switch debt := debt.(type) {
	case *structures.CreditDebt:
		db.deleteCreditDebt(debt.Id)
	case *structures.MortgageDebt:
		db.deleteCreditDebt(debt.Id)
	}
}

func (db *Database) deleteCreditDebt(debtId int) {
	_, err := db.conn.Exec("DELETE FROM credit_debt WHERE id=$1", debtId)
	if err != nil {
		panic(err)
	}
}

func (db *Database) deleteMortgageDebt(debtId int) {
	_, err := db.conn.Exec("DELETE FROM mortgage_debt WHERE id=$1", debtId)
	if err != nil {
		panic(err)
	}
}

func (db *Database) AddDebtToUser(user *structures.User, debt structures.Debt) {
	db.AddDebtToUserById(user.Id, debt)
}

func (db *Database) AddDebtToUserById(id int, debt structures.Debt) {
	switch debt := debt.(type) {
	case *structures.CreditDebt:
		db.addCreditDebtToUserById(id, debt)
	case *structures.MortgageDebt:
		db.addMortgageDebtToUserById(id, debt)
	}
}

func (db *Database) addCreditDebtToUserById(id int, debt *structures.CreditDebt) {
	_, err := db.conn.Exec("INSERT INTO credit_debt (id, percent, initial_amount, remainder, start_date, next_payment_date, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		debt.Id, debt.Percent, debt.InitialAmount, debt.Remainder, debt.StartDate, debt.NextPaymentDate, id)
	if err != nil {
		panic(err)
	}
}

func (db *Database) addMortgageDebtToUserById(id int, debt *structures.MortgageDebt) {
	_, err := db.conn.Exec("INSERT INTO mortgage_debt (id, percent, initial_amount, remainder, start_date, next_payment_date, user_id, address) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		debt.Id, debt.Percent, debt.InitialAmount, debt.Remainder, debt.StartDate, debt.NextPaymentDate, id, debt.Address)
	if err != nil {
		panic(err)
	}
}

func (db *Database) GetMaxUserId() int {
	return db.getMaxId("users")
}

func (db *Database) GetMaxCreditId() int {
	return db.getMaxId("credit_debt")
}

func (db *Database) GetMaxMortgageId() int {
	return db.getMaxId("mortgage_debt")
}

func (db *Database) getMaxId(tableName string) int {
	query := "SELECT MAX(id) FROM " + tableName
	var id int
	err := db.conn.QueryRow(query).Scan(&id)
	if err != nil {
		panic(err)
	}
	return id
}
