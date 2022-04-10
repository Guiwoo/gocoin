package main

import (
	"fmt"
	"time"
)

func countToTen(c chan<- int) {
	for i := range [100]int{} {
		fmt.Printf(">> sending << %d\n", i)
		c <- i
		fmt.Printf(">> sent << %d\n", i)
	}
	close(c)
}

func recive(c <-chan int) {
	for {
		time.Sleep(100 * time.Millisecond)
		a, ok := <-c
		if !ok {
			fmt.Println("we are done")
			break
		}
		fmt.Printf("|| received %d ||", a)
	}
}

func main() {
	// defer db.Close()
	// cli.Start()
	start := time.Now()
	c := make(chan int, 100)

	go countToTen(c)
	recive(c)
	fmt.Println(time.Since(start))

}
