package main

import (
	"fmt"
)

type sta int

func (s *sta) Set(a sta) {
	if a != EVEN && a != ODD && a != INVALID {
		panic("sta error")
	}
	*s = sta(a)
}

func (s *sta) Get() sta {
	return *s
}

const (
	EVEN sta = iota
	ODD
	INVALID
)

func main() {
	var state sta
	state.Set(5)
	fmt.Println(state)
	fmt.Println(state.Get())
	var ss2 sta
	ss2 = state.Get()
	fmt.Println(ss2)
}
