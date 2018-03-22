package astdata

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testStruct = `
package example

type ST1 struct {
}

type ST2 struct {
   A1 int ` + "`tag:\"tag_test\"`" + `
}

type ST3 struct {
   ST2 ` + "`tag:\"test_2\"`" + `
   A2 string
}

`

func TestStructType(t *testing.T) {
	p := &Package{}

	f, err := ParseFile(testStruct, p)
	assert.NoError(t, err)
	p.files = append(p.files, f)

	ST1, err := p.FindType("ST1")
	assert.NoError(t, err)
	assert.Equal(t, p, ST1.Package())
	assert.Equal(t, f, ST1.File())

	assert.IsType(t, &StructType{}, ST1.Definition())
	def := ST1.Definition().(*StructType)
	assert.Equal(t, 0, len(def.Fields()))
	assert.Equal(t, 0, len(def.Embeds()))
	assert.Equal(t, "struct{}", def.String())

	ST2, err := p.FindType("ST2")
	assert.NoError(t, err)
	assert.Equal(t, p, ST2.Package())
	assert.Equal(t, f, ST2.File())

	assert.IsType(t, &StructType{}, ST2.Definition())
	def = ST2.Definition().(*StructType)
	assert.Equal(t, 1, len(def.Fields()))
	assert.Equal(t, 0, len(def.Embeds()))
	f1 := def.Fields()[0]
	assert.Equal(t, "A1", f1.Name())
	assert.Equal(t, reflect.StructTag("tag:\"tag_test\""), f1.Tags())

	assert.Equal(t, "struct {\n\tA1 int `tag:\"tag_test\"`\n}", def.String())

	ST3, err := p.FindType("ST3")
	assert.NoError(t, err)
	assert.Equal(t, p, ST3.Package())
	assert.Equal(t, f, ST3.File())

	assert.IsType(t, &StructType{}, ST3.Definition())
	def = ST3.Definition().(*StructType)
	assert.Equal(t, 1, len(def.Fields()))
	assert.Equal(t, 1, len(def.Embeds()))
	f1 = def.Fields()[0]
	assert.Equal(t, "A2", f1.Name())
	assert.Equal(t, "string", f1.Definition().String())
	assert.Equal(t, "", f1.Docs().String())

	e1 := def.Embeds()[0]
	assert.Equal(t, reflect.StructTag("tag:\"test_2\""), e1.Tags())
	assert.Equal(t, "struct {\n\tST2\n\tA2 string \n}", def.String())
	assert.Equal(t, "", e1.Docs().String())

	assert.Equal(t, p, def.Package())
	assert.Equal(t, f, def.File())
}
