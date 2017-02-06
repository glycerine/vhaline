package main

import (
	"fmt"
)

// p is a shortcut for a call to fmt.Printf that implicitly starts
// and ends its message with a newline.
func p(format string, stuff ...interface{}) {
	fmt.Printf("\n "+format+"\n", stuff...)
}

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}
