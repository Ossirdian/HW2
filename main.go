package main

import (
	"errors"
	"fmt"
	"sync"
)

var (
	errNotEnoughFunds         = errors.New("Not enough founds for withdraw")
	errTransactionsQueueEmpty = errors.New("TransactionsQueue already empty")
)

type Transaction struct {
	FromID string
	ToID   string
	Amount float64
}

type PaymentSystem struct {
	Users            map[string]*User
	TransactionQueue []Transaction
}

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
		return errNotEnoughFunds
	} else {
		u.Balance = remainder
	}
	return nil
}

func (ps *PaymentSystem) AddUser(us *User) {
	ps.Users[us.ID] = us
}

func (ps *PaymentSystem) AddTransaction(tc Transaction) {
	ps.TransactionQueue = append(ps.TransactionQueue, tc)
}

func (ps *PaymentSystem) ProcessingTransactions(t Transaction) error {
	snd, ok := ps.Users[t.FromID]
	if !ok {
		return fmt.Errorf("User with ID %v not found", t.FromID)
	}
	rcv, ok := ps.Users[t.ToID]
	if !ok {
		return fmt.Errorf("User with ID %v not found", t.ToID)
	}
	if err := snd.Withdraw(t.Amount); err != nil {
		return err
	}
	rcv.Deposit(t.Amount)
	return nil
}

func (ps *PaymentSystem) Worker(ch <-chan Transaction, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range ch {
		if err := ps.ProcessingTransactions(t); err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	paymentSys := &PaymentSystem{
		Users:            make(map[string]*User),
		TransactionQueue: make([]Transaction, 0),
	}

	usr1 := &User{
		ID:      "1",
		Name:    "User1",
		Balance: 1000,
	}
	fmt.Println("Создаю UserID: 1 с балансом 1000")
	usr2 := &User{
		ID:      "2",
		Name:    "User2",
		Balance: 500,
	}
	fmt.Println("Создаю UserID: 2 с балансом 500")

	paymentSys.AddUser(usr1)
	paymentSys.AddUser(usr2)

	trc1 := Transaction{
		FromID: "1",
		ToID:   "2",
		Amount: 200,
	}
	trc2 := Transaction{
		FromID: "2",
		ToID:   "1",
		Amount: 50,
	}
	fmt.Println("Перевожу с UserID: 1 на UserID: 2 сумму в размере 200")
	fmt.Println("Перевожу с UserID: 2 на UserID: 1 сумму в размере 50")

	paymentSys.AddTransaction(trc1)
	paymentSys.AddTransaction(trc2)

	// TODO: Создать канал ch := ...
	ch := make(chan Transaction, len(paymentSys.TransactionQueue))
	// TODO: Инициализируем sync.WaitGroup
	wg := &sync.WaitGroup{}

	for range 3 {
		wg.Add(1)
		go paymentSys.Worker(ch, wg)
	}

	for _, t := range paymentSys.TransactionQueue {
		ch <- t
	}
	close(ch)
	wg.Wait()

	fmt.Println("Result")
	fmt.Printf("User1 balance is %v\n", usr1.Balance)
	fmt.Printf("User2 balance is %v\n", usr2.Balance)
}
