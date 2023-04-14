package llog

import (
	"fmt"
	"os"
	"runtime/debug"
)

func Throw(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, msg, args...)
	os.Exit(1)
}

func Raise(exception string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, "exception: %s\n", fmt.Sprintf(exception, args...))
	debug.PrintStack()
	os.Exit(1)
}

func BugOn(expr bool, msg string) {
	if !expr {
		Raise(msg)
	}
}
