package goql

import (
	"testing"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmbedFn(t *testing.T) {
	st, err := astdata.NewDefinition(`
struct {
ST1
ST2
}
`)
	require.NoError(t, err)
	require.IsType(t, &astdata.StructType{}, st)
	n, err := embedCountFn(0).Execute(Definition{Definition: st})
	assert.NoError(t, err)
	assert.Equal(t, 2.0, n.Get())

	n, err = embedDefFn(0).Execute(Definition{Definition: st}, Number{Number: 1})
	assert.NoError(t, err)
	assert.IsType(t, &astdata.IdentType{}, n.Get())
	assert.Equal(t, "ST1", n.Get().(astdata.Definition).String())

	n, err = embedDefFn(0).Execute(Definition{Definition: st}, Number{Number: 2})
	assert.NoError(t, err)
	assert.IsType(t, &astdata.IdentType{}, n.Get())
	assert.Equal(t, "ST2", n.Get().(astdata.Definition).String())

	i, err := astdata.NewDefinition("int")
	require.NoError(t, err)

	n, err = embedDefFn(0).Execute(Definition{Definition: st})
	assert.Error(t, err)
	assert.Nil(t, n)

	n, err = embedDefFn(0).Execute(Number{}, Number{})
	assert.NoError(t, err)
	assert.Nil(t, n.Get())

	n, err = embedCountFn(0).Execute(Definition{})
	assert.NoError(t, err)
	assert.Nil(t, n.Get())

	n, err = embedCountFn(0).Execute(Definition{Definition: i})
	assert.NoError(t, err)
	assert.Nil(t, n.Get())

	n, err = embedCountFn(0).Execute()
	assert.Error(t, err)
	assert.Nil(t, n)

	n, err = embedDefFn(0).Execute(Definition{Definition: st}, Number{Number: 10})
	assert.NoError(t, err)
	assert.Nil(t, n.Get())

	n, err = embedDefFn(0).Execute(Definition{Definition: i}, Number{Number: 0})
	assert.NoError(t, err)
	assert.Nil(t, n.Get())
}
