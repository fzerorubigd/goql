package goql

import (
	"testing"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsArray(t *testing.T) {
	st, err := astdata.NewDefinition("[]string")
	require.NoError(t, err)

	g, err := isArrayFn(0).Execute(Definition{Definition: st})
	assert.NoError(t, err)
	assert.Equal(t, Bool{Bool: true}, g)

	g, err = isArrayFn(0).Execute(nil, nil)
	assert.Error(t, err)

	g, err = isArrayFn(0).Execute(Number{})
	assert.NoError(t, err)
	assert.Equal(t, Bool{Null: true}, g)
}

func TestArrayItem(t *testing.T) {
	at, err := astdata.NewDefinition("[]string")
	require.NoError(t, err)

	et, err := astdata.NewDefinition("[...]string{}")
	require.NoError(t, err)

	ident, err := astdata.NewDefinition("int")
	require.NoError(t, err)

	_, err = arrayItemFn(0).Execute(nil, nil)
	assert.Error(t, err)

	g, err := arrayItemFn(0).Execute(Definition{})
	assert.NoError(t, err)
	assert.Nil(t, g.Get())

	g, err = arrayItemFn(0).Execute(Definition{Definition: ident})
	assert.NoError(t, err)
	assert.Nil(t, g.Get())

	g, err = arrayItemFn(0).Execute(Definition{Definition: at})
	assert.NoError(t, err)
	require.IsType(t, &astdata.IdentType{}, g.Get())
	assert.Equal(t, "string", g.Get().(astdata.Definition).String())

	g, err = arrayItemFn(0).Execute(Definition{Definition: et})
	assert.NoError(t, err)
	require.IsType(t, &astdata.IdentType{}, g.Get())
	assert.Equal(t, "string", g.Get().(astdata.Definition).String())
}

func TestArrayType(t *testing.T) {
	at, err := astdata.NewDefinition("[]string")
	require.NoError(t, err)

	at2, err := astdata.NewDefinition("[10]string")
	require.NoError(t, err)

	et, err := astdata.NewDefinition("[...]string{}")
	require.NoError(t, err)

	ident, err := astdata.NewDefinition("int")
	require.NoError(t, err)

	_, err = arrayTypeFn(0).Execute(nil, nil)
	assert.Error(t, err)

	g, err := arrayTypeFn(0).Execute(Definition{})
	assert.NoError(t, err)
	assert.Nil(t, g.Get())

	g, err = arrayTypeFn(0).Execute(Definition{Definition: ident})
	assert.NoError(t, err)
	assert.Nil(t, g.Get())

	g, err = arrayTypeFn(0).Execute(Definition{Definition: at})
	assert.NoError(t, err)
	assert.Equal(t, "slice", g.Get())

	g, err = arrayTypeFn(0).Execute(Definition{Definition: at2})
	assert.NoError(t, err)
	assert.Equal(t, "array", g.Get())

	g, err = arrayTypeFn(0).Execute(Definition{Definition: et})
	assert.NoError(t, err)
	assert.Equal(t, "ellipsis", g.Get())

}
