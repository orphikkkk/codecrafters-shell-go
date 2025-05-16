package main

import (
	"fmt"
	"os"
)

func dd(message any) {
	fmt.Println(message)
	os.Exit(1)
}
