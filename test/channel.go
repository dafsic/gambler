package main

import (
	"fmt"
	"time"
)

var ac chan int

func main() {
	ac = make(chan int, 8)
	go func() {
		for {
			a := <-ac
			fmt.Println(a)
			time.Sleep(1 * time.Second)
		}
	}()
	for i := 0; i < 20; i++ {
		w(i)
	}
	time.Sleep(10 * time.Second)

}

func w(a int) {
	fmt.Printf("w %d\n", a)
	select {
	case ac <- a:
	default:
		fmt.Println("wait")
	}
}
