package log

import (
	"fmt"
	"os"
)

func Debug(s string) {
	if os.Getenv("debug") == "" {
		//return
	}
	fmt.Println(s)
}
