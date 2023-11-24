package main

import (
	"fmt"
	"time"
)

func processando() {
	for i := 0; i < 10; i++ {
		fmt.Println(i)
		time.Sleep(time.Second)
	}

}

// T1
func main() {
	canal := make(chan int)

	go func() {
		canal <- 1 //T2
	}()

	fmt.Println(<-canal)
}
