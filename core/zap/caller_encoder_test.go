package zap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShort(t *testing.T) {
	assert.Equal(t, "", short(""))
	assert.Equal(t, "a/bbb/ccc.ddd.eee", short("aaa/bbb/ccc.ddd.eee"))
}
