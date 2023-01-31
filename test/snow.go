package main

//package snowID

import (
	"fmt"
	"os"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func init() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(12)
	}
}

func Generate() string {
	return node.Generate().String()
}

func main() {
	fmt.Printf("%v\n", Generate())
}
