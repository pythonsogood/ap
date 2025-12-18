package Shapes

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Print() string
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Print() string {
	return fmt.Sprintf("[Rectangle] Width: %f, Height: %f", r.Width, r.Height)
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * math.Pow(c.Radius, 2)
}

func (c Circle) Print() string {
	return fmt.Sprintf("[Circle] Radius: %f", c.Radius)
}

type Square struct {
	Side float64
}

func (s Square) Area() float64 {
	return math.Pow(s.Side, 2)
}

func (s Square) Print() string {
	return fmt.Sprintf("[Square] Side: %f", s.Side)
}

type Triangle struct {
	Side1 float64
	Side2 float64
	Side3 float64
}

func (t Triangle) Area() float64 {
	p := (t.Side1 + t.Side2 + t.Side3) / 2

	return math.Sqrt(p * (p - t.Side1) * (p - t.Side2) * (p - t.Side3))
}

func (t Triangle) Print() string {
	return fmt.Sprintf("[Triangle] Side 1: %f, Side 2: %f, Side 3: %f", t.Side1, t.Side2, t.Side3)
}

func remove_slice[T any](slice []T, i int) ([]T, error) {
	// https://stackoverflow.com/a/37335777/19338842

	if i < 0 || i >= len(slice) {
		return slice, fmt.Errorf("index %d out of range", i)
	}

	return append(slice[:i], slice[i+1:]...), nil
}

func CLI() {
	shapes := []Shape{}

	for {
		fmt.Println("\n--- Shapes ---")
		fmt.Println("[1] Add Rectangle")
		fmt.Println("[2] Add Circle")
		fmt.Println("[3] Add Square")
		fmt.Println("[4] Add Triangle")
		fmt.Println("[5] Remove Shape")
		fmt.Println("[0] Return")

		fmt.Println("\nShapes:")

		for i, shape := range shapes {
			fmt.Printf("%d %s\n", i, shape.Print())
		}

		fmt.Print("\n>>> ")

		var choice int

		fmt.Scanln(&choice)

		switch choice {
		case 1:
			var width, height float64

			fmt.Println("Enter Width:")
			fmt.Scanln(&width)

			fmt.Println("Enter Height:")
			fmt.Scanln(&height)

			shapes = append(shapes, Rectangle{Width: width, Height: height})

		case 2:
			var radius float64

			fmt.Println("Enter Radius:")
			fmt.Scanln(&radius)

			shapes = append(shapes, Circle{Radius: radius})

		case 3:
			var side float64

			fmt.Println("Enter Side:")
			fmt.Scanln(&side)

			shapes = append(shapes, Square{Side: side})

		case 4:
			var side1, side2, side3 float64

			fmt.Println("Enter Side 1:")
			fmt.Scanln(&side1)

			fmt.Println("Enter Side 2:")
			fmt.Scanln(&side2)

			fmt.Println("Enter Side 3:")
			fmt.Scanln(&side3)

			shapes = append(shapes, Triangle{Side1: side1, Side2: side2, Side3: side3})

		case 5:
			var index int

			fmt.Println("Enter Index:")
			fmt.Scanln(&index)

			shapes_, err := remove_slice(shapes, index)

			if err != nil {
				fmt.Println(err)
			} else {
				shapes = shapes_
			}

		case 0:
			return

		default:
			fmt.Println("Invalid choice")
		}
	}
}
