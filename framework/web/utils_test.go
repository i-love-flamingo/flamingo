package web

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLTitle(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "test_a-123-name-test", URLTitle("test/a 123 name % / _ - _ test"))
}

// Example usage of the URLTitle helper function
func ExampleURLTitle() {
	fmt.Println(URLTitle("test/a 123 name % / _ - _ test"))
	// Output: test_a-123-name-test
}
