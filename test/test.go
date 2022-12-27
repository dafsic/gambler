package main

import "fmt"

func main() {
	hash := "1eabcAdc"

	l := len(hash)
	for i := 1; i <= l; i++ {
		if hash[l-i] > '9' {
			continue
		}
		fmt.Println(aaa[hash[l-i]])
		return
	}
	fmt.Println("no number")
}

var aaa = map[byte]bool{'0': true, '1': false, '2': true, '3': false, '4': true, '5': false, '6': true, '7': false, '8': true, '9': false}
