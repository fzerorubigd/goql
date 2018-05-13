package goql

import (
	"testing"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsStruct(t *testing.T) {
	st, err := astdata.NewDefinition("struct{}")
	require.NoError(t, err)

	g, err := isStructFn(0).Execute(Definition{Definition: st})
	assert.NoError(t, err)
	assert.Equal(t, Bool{Bool: true}, g)

	_, err = isStructFn(0).Execute(nil, nil)
	assert.Error(t, err)

	g, err = isStructFn(0).Execute(Number{})
	assert.NoError(t, err)
	assert.Equal(t, Bool{Null: true}, g)
}

func TestStructFiledCount(t *testing.T) {
	st0, err := astdata.NewDefinition("struct{pkg.Struct}")
	require.NoError(t, err)

	st1, err := astdata.NewDefinition(`struct{
pkg.Struct
pkg2.Struct2
test int
}`)
	require.NoError(t, err)

	ident, err := astdata.NewDefinition("int")
	require.NoError(t, err)

	g, err := structFieldCountFn(0).Execute(Definition{Definition: st0})
	assert.NoError(t, err)
	assert.Equal(t, 0.0, g.Get())

	g, err = embedCountFn(0).Execute(Definition{Definition: st0})
	assert.NoError(t, err)
	assert.Equal(t, 1.0, g.Get())

	g, err = structFieldCountFn(0).Execute(Definition{Definition: st1})
	assert.NoError(t, err)
	assert.Equal(t, 1.0, g.Get())

	g, err = embedCountFn(0).Execute(Definition{Definition: st1})
	assert.NoError(t, err)
	assert.Equal(t, 2.0, g.Get())

	_, err = structFieldCountFn(0).Execute(nil, nil)
	assert.Error(t, err)

	g, err = structFieldCountFn(0).Execute(Definition{Definition: ident})
	assert.NoError(t, err)
	assert.Equal(t, Number{Null: true}, g)

	g, err = structFieldCountFn(0).Execute(Number{})
	assert.NoError(t, err)
	assert.Equal(t, Number{Null: true}, g)
}

func TestStructFiledDef(t *testing.T) {
	st1, err := astdata.NewDefinition(`struct{
pkg.Struct
pkg2.Struct2
test int
}`)
	require.NoError(t, err)

	ident, err := astdata.NewDefinition("int")
	require.NoError(t, err)

	_, err = structFieldDefFn(0).Execute(nil)
	assert.Error(t, err)

	def, err := structFieldDefFn(0).Execute(Definition{Definition: ident}, Number{Number: 1})
	assert.NoError(t, err)
	assert.Nil(t, def.Get())

	// filed
	def, err = structFieldDefFn(0).Execute(Definition{Definition: st1}, Number{Number: 1})
	assert.NoError(t, err)
	require.IsType(t, &astdata.IdentType{}, def.Get())
	assert.Equal(t, "int", def.Get().(astdata.Definition).String())

	// embed
	def, err = embedDefFn(0).Execute(Definition{Definition: st1}, Number{Number: 1})
	assert.NoError(t, err)
	require.IsType(t, &astdata.SelectorType{}, def.Get())
	assert.Equal(t, "pkg.Struct", def.Get().(astdata.Definition).String())

	def, err = structFieldDefFn(0).Execute(Definition{Definition: st1}, String{String: "test"})
	assert.NoError(t, err)
	require.IsType(t, &astdata.IdentType{}, def.Get())
	assert.Equal(t, "int", def.Get().(astdata.Definition).String())

	def, err = structFieldDefFn(0).Execute(Definition{Definition: st1}, String{String: "notexist"})
	assert.NoError(t, err)
	assert.Nil(t, def.Get())

	// embeds are not available with name (yet)
	def, err = structFieldDefFn(1).Execute(Definition{Definition: st1}, String{String: "pkg.Struct"})
	assert.NoError(t, err)
	assert.Nil(t, def.Get())

	def, err = structFieldDefFn(0).Execute(Definition{Definition: st1}, Number{})
	assert.NoError(t, err)
	assert.Nil(t, def.Get())

	def, err = structFieldDefFn(1).Execute(Definition{Definition: st1}, Number{})
	assert.NoError(t, err)
	assert.Nil(t, def.Get())

	def, err = structFieldDefFn(0).Execute(Definition{Definition: st1}, Bool{})
	assert.NoError(t, err)
	assert.Nil(t, def.Get())

}

func TestStructFieldName(t *testing.T) {
	st1, err := astdata.NewDefinition(`struct{
pkg.Struct
pkg2.Struct2
test int
}`)
	require.NoError(t, err)

	ident, err := astdata.NewDefinition("int")
	require.NoError(t, err)

	_, err = structFieldNameFn(0).Execute(nil)
	assert.Error(t, err)

	def, err := structFieldNameFn(0).Execute(Definition{Definition: ident}, Number{Number: 1})
	assert.NoError(t, err)
	assert.Nil(t, def.Get())

	def, err = structFieldNameFn(0).Execute(Definition{Definition: st1}, Number{Number: 1})
	assert.NoError(t, err)
	assert.Equal(t, "test", def.Get())

	def, err = structFieldNameFn(0).Execute(Definition{Definition: st1}, Number{})
	assert.NoError(t, err)
	assert.Equal(t, nil, def.Get())

}

func TestStructFiledTag(t *testing.T) {
	st1, err := astdata.NewDefinition(`struct{
pkg.Struct ` + "`tag:\"test\"`" + `
pkg2.Struct2
test int ` + "`tag2:\"test2\"`" + `
}`)
	require.NoError(t, err)

	ident, err := astdata.NewDefinition("int")
	require.NoError(t, err)

	_, err = structFieldTagFn(0).Execute(nil)
	assert.Error(t, err)

	def, err := structFieldTagFn(0).Execute(Definition{Definition: ident}, Number{Number: 1})
	assert.NoError(t, err)
	assert.Nil(t, def.Get())

	// filed
	def, err = structFieldTagFn(0).Execute(Definition{Definition: st1}, Number{Number: 1})
	assert.NoError(t, err)
	assert.Equal(t, "tag2:\"test2\"", def.Get())

	def, err = structFieldTagFn(0).Execute(Definition{Definition: st1}, Number{Number: 1}, String{String: "tag2"})
	assert.NoError(t, err)
	assert.Equal(t, "test2", def.Get())

	// embed
	def, err = structFieldTagFn(1).Execute(Definition{Definition: st1}, Number{Number: 1})
	assert.NoError(t, err)
	assert.Equal(t, "tag:\"test\"", def.Get())

	def, err = structFieldTagFn(1).Execute(Definition{Definition: st1}, Number{Number: 1}, String{String: "tag"})
	assert.NoError(t, err)
	assert.Equal(t, "test", def.Get())

	def, err = structFieldTagFn(0).Execute(Definition{Definition: st1}, String{String: "test"})
	assert.NoError(t, err)
	assert.Equal(t, "tag2:\"test2\"", def.Get())

	def, err = structFieldTagFn(0).Execute(Definition{Definition: st1}, String{String: "notexist"})
	assert.NoError(t, err)
	assert.Nil(t, def.Get())

	// embeds are not available with name (yet)
	def, err = structFieldTagFn(1).Execute(Definition{Definition: st1}, String{String: "pkg.Struct"})
	assert.NoError(t, err)
	assert.Nil(t, def.Get())

	def, err = structFieldTagFn(0).Execute(Definition{Definition: st1}, Number{})
	assert.NoError(t, err)
	assert.Nil(t, def.Get())

	def, err = structFieldTagFn(1).Execute(Definition{Definition: st1}, Number{})
	assert.NoError(t, err)
	assert.Nil(t, def.Get())

	def, err = structFieldTagFn(0).Execute(Definition{Definition: st1}, Bool{})
	assert.NoError(t, err)
	assert.Nil(t, def.Get())

}
