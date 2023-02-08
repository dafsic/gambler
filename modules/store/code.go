package store

import (
	"strconv"
	"strings"
)

type Encoder struct{}

func (e *Encoder) String(src []int) string {
	var b strings.Builder
	for _, v := range src {
		b.WriteString(strconv.Itoa(v) + "-")
	}

	return b.String()
}

func (e *Encoder) IntArray(src string) ([]int, error) {
	a := strings.Split(src, "-")
	var r []int
	var i int
	var err error
	for _, v := range a {
		if v == "" {
			return r, nil
		}

		if i, err = strconv.Atoi(v); err != nil {
			return nil, err
		}
		r = append(r, i)
	}
	return r, nil
}
