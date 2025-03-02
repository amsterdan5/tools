package main

import (
	"fmt"

	. "github.com/amsterdan/tools/builder/lib"
)

// usage: -a "left op" -b right op""
func Add() {
	a := OptInt("a")
	b := OptInt("b")
	fmt.Print(a + b)
}

// usage: -a "left op" -b right op""
func Sub() {
	a := OptInt("a")
	b := OptInt("b")
	fmt.Print(a - b)
}
func SayHello() {
	fmt.Printf("hello world\n")
}
