package main

import (
	"fmt"
	"math"
	"math/cmplx"
)

func quadratic_roots(a float64, b float64, c float64) (complex128, complex128) {
	D := complex(math.Pow(b, 2)-4*a*c, 0)
	D_sqrt := cmplx.Sqrt(D)

	numerator_b := complex(-b, 0)
	denominator := complex(2*a, 0)

	x1 := (numerator_b + D_sqrt) / denominator
	x2 := (numerator_b - D_sqrt) / denominator

	return x1, x2
}

func main() {
	var a, b, c float64

	fmt.Println("Enter a, b, c for equation a*x^2 + b*x + c = 0")

	_, scan_error := fmt.Scan(&a, &b, &c)

	for scan_error != nil {
		fmt.Println(scan_error, "\nTry again")
		_, scan_error = fmt.Scan(&a, &b, &c)
	}

	x1, x2 := quadratic_roots(a, b, c)

	fmt.Println("Roots:", x1, x2)
}
