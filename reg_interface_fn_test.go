package goql

import (
	"testing"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsInterface(t *testing.T) {
	in, err := astdata.NewDefinition("interface{}")
	require.NoError(t, err)

	g, err := isInterfaceFn(0).Execute(Definition{Definition: in})
	assert.NoError(t, err)
	assert.Equal(t, Bool{Bool: true}, g)

	g, err = isInterfaceFn(0).Execute(nil, nil)
	assert.Error(t, err)

	g, err = isInterfaceFn(0).Execute(Number{})
	assert.NoError(t, err)
	assert.Equal(t, Bool{Null: true}, g)

}
