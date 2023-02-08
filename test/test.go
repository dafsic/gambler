package main

import "fmt"

type T struct {
	a int
	b string
}

func main() {
	var t T
	var a int
	var s string

	a = 5
	s = "12345"

	*(&(t.a)) = a
	*(&(t.b)) = s
	fmt.Printf("%+v\n", t)
}
