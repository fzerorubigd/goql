package astdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefinition(t *testing.T) {
	require.Nil(t, newType(nil, nil, nil))
	def, err := NewDefinition("int")
	assert.NoError(t, err)
	assert.IsType(t, &IdentType{}, def)
	assert.Equal(t, "int", def.String())

	def, err = NewDefinition(`struct{ data int }`)
	assert.NoError(t, err)
	assert.IsType(t, &StructType{}, def)
	assert.Equal(t, "struct {\n\tdata int \n}", def.String())

	_, err = NewDefinition("x := 10")
	assert.Error(t, err)

	_, err = NewDefinition("invalid code")
	assert.Error(t, err)

	vv, err := NewDefinition("func(context.Context)error")
	assert.NoError(t, err)
	assert.IsType(t, &FuncType{}, vv)
	require.Equal(t, 1, len(vv.(*FuncType).Parameters()))
	assert.IsType(t, &SelectorType{}, vv.(*FuncType).Parameters()[0].Definition())
}
