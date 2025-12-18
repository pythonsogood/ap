package main

import (
	"fmt"

	"github.com/SultanYakupov/Assignment1/Bank"
	"github.com/SultanYakupov/Assignment1/Company"
	"github.com/SultanYakupov/Assignment1/Library"
	"github.com/SultanYakupov/Assignment1/Shapes"
)

func main() {
main_menu:
	for {
		var choice int

		fmt.Println("\n--- Main menu ---")
		fmt.Println("[1] Library Management System")
		fmt.Println("[2] Shapes & Interfaces")
		fmt.Println("[3] Employee Management System")
		fmt.Println("[4] Bank Account Simulation")
		fmt.Println("[0] Exit")
		fmt.Print("\n>>> ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			Library.CLI()

		case 2:
			Shapes.CLI()

		case 3:
			Company.CLI()

		case 4:
			Bank.CLI()

		case 0:
			break main_menu

		default:
			fmt.Println("Invalid choice")
		}
	}
}
