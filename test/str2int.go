package main

import (
	"fmt"
)

func main() {

	var addr = "417e36ce97d6f8a5e47b1c95ae01fe815bc8f2c8cd"
	var para = "0000000000000000000000007e36ce97d6f8a5e47b1c95ae01fe815bc8f2c8cd00000000000000000000000000000000000000000000000000000000000f4240"
	p := encodeParameter(addr, 1000000)
	fmt.Println(p)
	if p == para {
		fmt.Println("true")
	}
}

func encodeParameter(addr string, amount int64) string {
	p := fmt.Sprintf("%064s%064x", addr[2:], amount)
	return p
}
