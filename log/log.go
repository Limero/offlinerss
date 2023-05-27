package log

import (
	"fmt"
	"os"
)

func Debug(format string, a ...any) {
	if os.Getenv("DEBUG") == "" {
		return
	}
	fmt.Printf(format+"\n", a...)
}

func Info(format string, a ...any) {
	fmt.Printf(format+"\n", a...)
}
