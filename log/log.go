package log

import (
	"fmt"
	"os"
)

func Debug(format string, a ...any) {
	if os.Getenv("debug") == "" {
		//return
	}
	fmt.Printf(format+"\n", a...)
}
