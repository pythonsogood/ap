package Bank

import (
	"errors"
	"fmt"
)

type BankAccount struct {
	Balance uint64
}

func (b *BankAccount) Deposit(amount uint64) {
	b.Balance += amount
}

func (b *BankAccount) Withdraw(amount uint64) error {
	if b.Balance < amount {
		return errors.New("Insufficient funds")
	}

	b.Balance -= amount

	return nil
}

func (b *BankAccount) GetBalance() uint64 {
	return b.Balance
}

func NewBankAccount() *BankAccount {
	return &BankAccount{Balance: 0}
}

func CLI() {
	bank := NewBankAccount()

	for {
		fmt.Println("\n--- Bank ---")
		fmt.Println("[1] Deposit")
		fmt.Println("[2] Withdraw")
		fmt.Println("[3] Get Balance")
		fmt.Println("[0] Return")
		fmt.Print("\n>>> ")

		var choice int

		fmt.Scanln(&choice)

		switch choice {
		case 1:
			var amount uint64

			fmt.Println("Enter amount to deposit:")
			fmt.Scanln(&amount)

			bank.Deposit(amount)

		case 2:
			var amount uint64

			fmt.Println("Enter amount to withdraw:")
			fmt.Scanln(&amount)

			err := bank.Withdraw(amount)

			if err != nil {
				fmt.Println(err)
			}

		case 3:
			fmt.Println("Balance:", bank.GetBalance())

		case 0:
			return

		default:
			fmt.Println("Invalid choice")
		}
	}
}
