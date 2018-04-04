package goql

import (
	"testing"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsMap(t *testing.T) {
	st, err := astdata.NewDefinition("map[string]int")
	require.NoError(t, err)

	g, err := isMapFn(0).Execute(Definition{Definition: st})
	assert.NoError(t, err)
	assert.Equal(t, Bool{Bool: true}, g)

	_, err = isMapFn(0).Execute(nil, nil)
	assert.Error(t, err)

	g, err = isMapFn(0).Execute(Number{})
	assert.NoError(t, err)
	assert.Equal(t, Bool{Null: true}, g)

}

func TestMapKey(t *testing.T) {
	st, err := astdata.NewDefinition("map[string]int")
	require.NoError(t, err)

	g, err := mapKeyFn(0).Execute(Definition{Definition: st})
	assert.NoError(t, err)
	assert.IsType(t, &astdata.IdentType{}, g.Get())
	assert.Equal(t, "string", g.Get().(astdata.Definition).String())

	g, err = mapKeyFn(0).Execute(g)
	assert.NoError(t, err)
	assert.Nil(t, g.Get())

	_, err = mapKeyFn(0).Execute(nil, nil)
	assert.Error(t, err)

	g, err = mapKeyFn(0).Execute(Number{})
	assert.NoError(t, err)
	assert.Equal(t, Definition{}, g)

}

func TestMapValue(t *testing.T) {
	st, err := astdata.NewDefinition("map[string]int")
	require.NoError(t, err)

	g, err := mapValFn(0).Execute(Definition{Definition: st})
	assert.NoError(t, err)
	assert.IsType(t, &astdata.IdentType{}, g.Get())
	assert.Equal(t, "int", g.Get().(astdata.Definition).String())

	g, err = mapValFn(0).Execute(g)
	assert.NoError(t, err)
	assert.Nil(t, g.Get())

	_, err = mapValFn(0).Execute(nil, nil)
	assert.Error(t, err)

	g, err = mapValFn(0).Execute(Number{})
	assert.NoError(t, err)
	assert.Equal(t, Definition{}, g)

}
