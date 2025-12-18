package Company

import (
	"errors"
	"fmt"
)

type Employee interface {
	GetDetails() string
}

type FullTimeEmployee struct {
	Id   uint64
	Name string
}

func (f FullTimeEmployee) GetDetails() string {
	return fmt.Sprintf("Full Time Employee Id: %d, Name: %s", f.Id, f.Name)
}

type PartTimeEmployee struct {
	Id   uint64
	Name string
}

func (p PartTimeEmployee) GetDetails() string {
	return fmt.Sprintf("Part Time Employee Id: %d, Name: %s", p.Id, p.Name)
}

type Company struct {
	Employees map[uint64]Employee
}

func (c *Company) AddEmployee(id uint64, name string, full_time bool) (Employee, error) {
	employee, ok := c.Employees[id]

	if ok {
		return nil, errors.New("Employee already exists")
	}

	if full_time {
		employee = FullTimeEmployee{Id: id, Name: name}
	} else {
		employee = PartTimeEmployee{Id: id, Name: name}
	}

	c.Employees[id] = employee

	return employee, nil
}

func (c *Company) ListEmployees() []Employee {
	var employees []Employee

	for _, employee := range c.Employees {
		employees = append(employees, employee)
	}

	return employees
}

func NewCompany() *Company {
	return &Company{Employees: make(map[uint64]Employee)}
}

func CLI() {
	company := NewCompany()

	for {
		fmt.Println("\n--- Company ---")
		fmt.Println("[1] Add Employee")
		fmt.Println("[2] List Employees")
		fmt.Println("[0] Return")
		fmt.Print("\n>>> ")

		var choice int

		fmt.Scanln(&choice)

		switch choice {
		case 1:
			var id uint64
			var name string
			var full_time bool

			fmt.Println("Enter Employee ID:")
			fmt.Scanln(&id)

			fmt.Println("Enter Employee Name:")
			fmt.Scanln(&name)

			fmt.Println("Is Employee Full Time? (true/false)")
			fmt.Scanln(&full_time)

			_, err := company.AddEmployee(id, name, full_time)

			if err != nil {
				fmt.Println(err)
			}

		case 2:
			employees := company.ListEmployees()

			for _, employee := range employees {
				fmt.Println(employee.GetDetails())
			}

		case 0:
			return

		default:
			fmt.Println("Invalid choice")
		}
	}
}
