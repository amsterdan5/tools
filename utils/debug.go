package utils

import (
	"log"
	"runtime/debug"
	"strings"
)

func CatchPanic(panicCallback func(err interface{})) {
	if err := recover(); err != nil {
		log.Printf("PROCESS PANIC: err %s", err)
		st := debug.Stack()
		if len(st) > 0 {
			log.Printf("dump stack (%s):", err)
			lines := strings.Split(string(st), "\n")
			for _, line := range lines {
				log.Print("  ", line)
			}
		} else {
			log.Printf("stack is empty (%s)", err)
		}
		if panicCallback != nil {
			panicCallback(err)
		}
	}
}
