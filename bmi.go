package main

import (
	"fmt"
	"math"
)

func bmi(weight float64, height float64) float64 {
	return weight / math.Pow(height, 2)
}

func main() {
	var weight_kg, height_m float64

	fmt.Scan(&weight_kg, &height_m)

	bmi_kg_m2 := bmi(weight_kg, height_m)

	fmt.Println(bmi_kg_m2)
}
