package main

import (
	"fmt"
	"log"
)

const (
	VERSION = "0.01"
)

func init() {
	log.SetFlags(log.Llongfile)
}

func main() {
	fmt.Println("vim-go")
}
