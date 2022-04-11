package main

import (
	"fmt"
)

func main() {
	// defer db.Close()
	// cli.Start()
	a := 0
	b := 0
	result := 0
	for i := 1; i < 54; i++ {
		a = (8900 + (9 * i))
		b = (1250 + (3 * i * i))
		result = a - b
		if result < 0 {
			fmt.Println("ho", i)
		}
	}
}
