package pugjs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNil(t *testing.T) {
	n := Nil{}

	assert.Equal(t, false, n.True())
	assert.Equal(t, "", n.String())
	assert.Equal(t, Nil{}, n.Member(""))
	assert.Equal(t, Nil{}, n.Member("aaa"))
	assert.Equal(t, Nil{}, n.copy())
}

func TestBool(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		b := Bool(true)
		assert.Equal(t, true, b.True())
		assert.Equal(t, "true", b.String())
		assert.Equal(t, Nil{}, b.Member(""))
		assert.Equal(t, Nil{}, b.Member("aaa"))
		assert.Equal(t, Bool(true), b.copy())
	})

	t.Run("false", func(t *testing.T) {
		b := Bool(false)
		assert.Equal(t, false, b.True())
		assert.Equal(t, "false", b.String())
		assert.Equal(t, Nil{}, b.Member(""))
		assert.Equal(t, Nil{}, b.Member("aaa"))
		assert.Equal(t, Bool(false), b.copy())
	})
}

func TestNumber(t *testing.T) {
	n := Number(1.2)

	assert.Equal(t, "1.2", n.String())
	assert.Equal(t, "1", Number(1).String())
	assert.Equal(t, "0", Number(0).String())
	assert.Equal(t, "-1", Number(-1).String())

	assert.Equal(t, Nil{}, n.Member(""))
	assert.Equal(t, Nil{}, n.Member("aaa"))

	assert.Equal(t, n, n.copy())
}

func TestArray_Splice(t *testing.T) {
	arr := convert([]int{1, 2, 3, 4, 5}).(*Array)

	assert.Len(t, arr.items, 5)
	leftover := arr.Splice(Number(2)).(*Array)
	assert.Len(t, arr.items, 2)
	assert.Len(t, leftover.items, 3)

	assert.Contains(t, arr.items, Number(1))
	assert.Contains(t, arr.items, Number(2))
	assert.Contains(t, leftover.items, Number(3))
	assert.Contains(t, leftover.items, Number(4))
	assert.Contains(t, leftover.items, Number(5))
}

func TestArray_Slice(t *testing.T) {
	arr := convert([]int{1, 2, 3, 4, 5}).(*Array)

	assert.Len(t, arr.items, 5)
	leftover := arr.Slice(Number(2)).(*Array)
	assert.Len(t, arr.items, 5)
	assert.Len(t, leftover.items, 3)

	assert.Contains(t, arr.items, Number(1))
	assert.Contains(t, arr.items, Number(2))
	assert.Contains(t, arr.items, Number(3))
	assert.Contains(t, arr.items, Number(4))
	assert.Contains(t, arr.items, Number(5))
	assert.Contains(t, leftover.items, Number(3))
	assert.Contains(t, leftover.items, Number(4))
	assert.Contains(t, leftover.items, Number(5))
}
