package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCond(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		assert.Equal(t, "hello", Cond(true, "hello", "world"))
		assert.Equal(t, "world", Cond(false, "hello", "world"))
	})

	t.Run("ints", func(t *testing.T) {
		assert.Equal(t, 13, Cond(true, 13, 37))
		assert.Equal(t, 37, Cond(false, 13, 37))
	})
}
