package log

import (
	"fmt"
)

var DebugEnabled = false

func Debug(format string, a ...any) {
	if !DebugEnabled {
		return
	}
	fmt.Printf(format+"\n", a...)
}

func Info(format string, a ...any) {
	fmt.Printf(format+"\n", a...)
}

func Warn(format string, a ...any) {
	fmt.Printf("WARN! "+format+"\n", a...)
}
