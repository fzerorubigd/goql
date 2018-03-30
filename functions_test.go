package goql

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fn struct {
}

//
func (f *fn) Execute(v ...Valuer) (Valuer, error) {
	if len(v) != 2 {
		return nil, fmt.Errorf("test function need two arg")
	}

	switch t := v[0].Value().(type) {
	case string:
		if r, ok := v[1].Value().(string); ok {
			return String{String: t + r}, nil
		}
	}

	return nil, fmt.Errorf("invalid types %T, %T", v[0].Value(), v[1].Value())
}

func TestRegister(t *testing.T) {
	concat := &fn{}
	assert.NotPanics(t, func() { RegisterFunction("fnconcat", concat) })
	assert.Panics(t, func() { RegisterFunction("fnconcat", concat) })
	assert.True(t, hasFunction("fnconcat"))
	assert.False(t, hasFunction("invalid-func"))

	res, err := executeFunction("fnconcat", String{}, Bool{})
	assert.Error(t, err)
	assert.Nil(t, res)

	res, err = executeFunction("fnconcat", String{String: "Hello"}, String{String: "World"})
	assert.NoError(t, err)
	assert.IsType(t, String{}, res)
	assert.Equal(t, "HelloWorld", res.Value().(string))

	res, err = executeFunction("notexists")
	assert.Error(t, err)
	assert.Nil(t, res)
}
