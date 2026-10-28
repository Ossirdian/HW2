package main

import (
	"errors"
	"fmt"
)

var errNotEnoughFounds = errors.New("Недостаточо средств на балансе для списания")

type User struct {
	ID      string
	Name    string
	Balance float64
}

func (u *User) Deposit(deposit float64) {
	u.Balance = u.Balance + deposit
}

func (u *User) Withdraw(withdraw float64) error {
	if remainder := u.Balance - withdraw; remainder < 0 {
		return errNotEnoughFounds
	} else {
		u.Balance = remainder
	}
	return nil
}

func main() {
	user1 := &User{
		ID:      "0",
		Name:    "user1",
		Balance: 2000.62,
	}
	user2 := &User{
		ID:      "1",
		Name:    "user2",
		Balance: 321.45,
	}
	user1.Withdraw(1000)
	println("Баланс пользователя 1: ", user1.Balance)
	user1.Deposit(2000)
	println("Баланс пользователя 1: ", user1.Balance)
	if err := user1.Withdraw(20000); err != nil {
		fmt.Println(err)
	}
	println("Баланс пользователя 1: ", user1.Balance)
	user2.Deposit(20000)
	println("Баланс пользователя 2: ", user2.Balance)
	user2.Withdraw(2000.41)
	println("Баланс пользователя 2: ", user2.Balance)
}
