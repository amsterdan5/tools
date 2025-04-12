package main

import (
	tools_lib "github.com/amsterdan5/tools/builder/lib"
)

func wrapperAdd() {
	Add()
}
func wrapperSub() {
	Sub()
}
func wrapperSayHello() {
	SayHello()
}
func main() {
	tools_lib.Register("Add", `-a "left op" -b right op""`, wrapperAdd)
	tools_lib.Register("Sub", `-a "left op" -b right op""`, wrapperSub)
	tools_lib.Register("SayHello", ``, wrapperSayHello)
	tools_lib.Run()
}
