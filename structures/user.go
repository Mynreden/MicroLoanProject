package structures

import (
	"fmt"
)

type User struct {
	Id       int       `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Debts    []Debt    `json:"debts"`
	Contacts []Contact `json:"contacts"`
}

func (u *User) AddDebt(dept Debt) {
	u.Debts = append(u.Debts, dept)
}

func (u *User) PayDebt(debt Debt) {
	for i, d := range u.Debts {
		if d == debt {
			u.Debts = append(u.Debts[:i], u.Debts[i+1:]...)
		}
	}
}

func (u *User) AddEmailContact(contact string) {
	u.Contacts = append(u.Contacts, Email{contact})
}

func (u *User) AddPhoneContact(contact string) {
	u.Contacts = append(u.Contacts, PhoneNumber{contact})
}

func (u *User) Notify() {
	for _, contact := range u.Contacts {
		contact.Update()
	}
}

type Contact interface {
	Update()
}

type PhoneNumber struct {
	Number string `json:"number"`
}

func (ph PhoneNumber) Update() {
	fmt.Println("Phone notifier to pay credit")
}

type Email struct {
	Email string `json:"email"`
}

func (e Email) Update() {
	fmt.Println("Email notifier to pay credit")
}
