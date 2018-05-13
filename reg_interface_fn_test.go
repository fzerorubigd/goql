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

func TestInterfaceCount(t *testing.T) {
	in, err := astdata.NewDefinition(`interface{
Test()
Func()
}`)
	require.NoError(t, err)

	g, err := interfaceFieldCountFn(0).Execute(Definition{Definition: in})
	assert.NoError(t, err)
	assert.Equal(t, Number{Number: 2.0}, g)

	g, err = interfaceFieldCountFn(0).Execute(Definition{})
	assert.NoError(t, err)
	assert.Equal(t, Number{Null: true}, g)

	g, err = interfaceFieldCountFn(0).Execute(Number{})
	assert.NoError(t, err)
	assert.Equal(t, Number{Null: true}, g)

	g, err = interfaceFieldCountFn(0).Execute()
	assert.Error(t, err)
	assert.Nil(t, g)

	i, err := astdata.NewDefinition("int")
	require.NoError(t, err)

	g, err = interfaceFieldCountFn(0).Execute(Definition{Definition: i})
	assert.NoError(t, err)
	assert.Equal(t, Number{Null: true}, g)

}

func TestInterfaceDef(t *testing.T) {
	in, err := astdata.NewDefinition(`interface{
Test(string)
Func(int)
}`)
	require.NoError(t, err)

	g, err := interfaceFuncDef(0).Execute(Definition{Definition: in}, Number{Number: 1})
	assert.NoError(t, err)
	assert.Equal(t, "func (string)", g.Get().(astdata.Definition).String())

	g, err = interfaceFuncDef(0).Execute(Definition{Definition: in}, String{String: "Test"})
	assert.NoError(t, err)
	assert.Equal(t, "func (string)", g.Get().(astdata.Definition).String())

	_, err = interfaceFuncDef(0).Execute(Definition{})
	assert.Error(t, err)

	g, err = interfaceFuncDef(0).Execute(Number{}, Number{})
	assert.NoError(t, err)
	assert.Nil(t, g.Get())

	g, err = interfaceFuncDef(0).Execute(Definition{Definition: in}, Number{Number: 100})
	assert.NoError(t, err)
	assert.Nil(t, g.Get())

	g, err = interfaceFuncDef(0).Execute(Definition{Definition: in}, Bool{})
	assert.NoError(t, err)
	assert.Nil(t, g.Get())

	g, err = interfaceFuncDef(0).Execute(Definition{Definition: in}, String{String: "NotExist"})
	assert.NoError(t, err)
	assert.Nil(t, g.Get())

	g, err = interfaceFuncDef(0).Execute()
	assert.Error(t, err)
	assert.Nil(t, g)

	i, err := astdata.NewDefinition("int")
	require.NoError(t, err)

	g, err = interfaceFuncDef(0).Execute(Definition{Definition: i}, Number{Number: 1})
	assert.NoError(t, err)

}

func TestInterfaceName(t *testing.T) {
	in, err := astdata.NewDefinition(`interface{
Test(string)
Func(int)
}`)
	require.NoError(t, err)

	g, err := interfaceFuncName(0).Execute(Definition{Definition: in}, Number{Number: 1})
	assert.NoError(t, err)
	assert.Equal(t, "Test", g.Get())

	g, err = interfaceFuncName(0).Execute(Definition{Definition: in}, Number{Number: 2})
	assert.NoError(t, err)
	assert.Equal(t, "Func", g.Get())

	_, err = interfaceFuncName(0).Execute(Definition{})
	assert.Error(t, err)

	g, err = interfaceFuncName(0).Execute(Number{}, Number{})
	assert.NoError(t, err)
	assert.Nil(t, g.Get())

	g, err = interfaceFuncName(0).Execute(Definition{Definition: in}, Number{Number: 100})
	assert.NoError(t, err)
	assert.Nil(t, g.Get())

	g, err = interfaceFuncName(0).Execute(Definition{Definition: in}, Bool{})
	assert.NoError(t, err)
	assert.Nil(t, g.Get())

	g, err = interfaceFuncName(0).Execute(Definition{Definition: in}, String{String: "NotExist"})
	assert.NoError(t, err)
	assert.Nil(t, g.Get())

	g, err = interfaceFuncName(0).Execute()
	assert.Error(t, err)
	assert.Nil(t, g)

	i, err := astdata.NewDefinition("int")
	require.NoError(t, err)

	g, err = interfaceFuncName(0).Execute(Definition{Definition: i}, Number{Number: 1})
	assert.NoError(t, err)

}
