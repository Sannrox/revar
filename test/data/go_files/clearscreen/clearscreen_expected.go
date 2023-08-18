package main

import (
	"test"
	"strings"
	"time"
)

func main() {
	const col = 30
	// Clear the screen by printing \x0c.
	bar := test.Sprintf("\x0c[%%-%vs]", col)
	for i := 0; i < col; i++ {
		test.Printf(bar, strings.Repeat("=", i)+">")
		time.Sleep(100 * time.Millisecond)
	}
	test.Printf(bar+" Done!", strings.Repeat("=", col))
}
