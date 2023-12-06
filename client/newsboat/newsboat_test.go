package newsboat

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_urlsSortFunc(t *testing.T) {
	lines := []string{
		"1000",
		"111 \"BBB\"",
		"999 \"AAA\"",
	}

	sortFunc := urlsSortFunc()
	sort.Slice(lines, func(i, j int) bool {
		return sortFunc(lines[i], lines[j])
	})

	assert.Equal(t, []string{
		"999 \"AAA\"",
		"111 \"BBB\"",
		"1000",
	}, lines)
}
