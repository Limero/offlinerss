package helpers

import (
	"fmt"
	"os/exec"
	"strings"
)

func SqlDiff(db1 string, db2 string) ([]string, error) {
	// Output SQL text that would transform DB1 into DB2
	fmt.Printf("Comparing database %s with %s\n", db1, db2)
	out, err := exec.Command("sqldiff", db1, db2).Output()
	if err != nil {
		return nil, fmt.Errorf("sqldiff failed: %w", err)
	}

	return strings.Split(string(out[:]), "\n"), nil
}
