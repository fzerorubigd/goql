package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	s := NewStack(10)
	o, err := s.Peek()
	assert.Error(t, err, "stack is empty")
	assert.Equal(t, o, item{})

	o, err = s.Pop()
	assert.Error(t, err, "stack is empty")
	assert.Equal(t, o, item{})

	s.Push(item{
		typ:   itemAlpha,
		pos:   10,
		value: "string",
	})

	o, err = s.Peek()
	assert.NoError(t, err)
	assert.Equal(t, o.typ, itemAlpha)
	assert.Equal(t, o.String(), "pos 10, token string")

	s.Push(item{
		typ:   itemAnd,
		pos:   11,
		value: "and",
	})

	o, err = s.Peek()
	assert.NoError(t, err)
	assert.Equal(t, o.typ, itemAnd)
	assert.Equal(t, o.String(), "pos 11, token and")

	o, err = s.Pop()
	assert.NoError(t, err)
	assert.Equal(t, o.typ, itemAnd)
	assert.Equal(t, o.String(), "pos 11, token and")

	o, err = s.Pop()
	assert.NoError(t, err)
	assert.Equal(t, o.typ, itemAlpha)
	assert.Equal(t, o.String(), "pos 10, token string")

}
