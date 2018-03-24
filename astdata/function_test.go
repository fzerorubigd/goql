package astdata

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testFunction = `
package example

type E int 

func Testing(s string) error {
	return nil 
}

func (e *E) Testing(s string) error {
	return nil 
}

func (e E) NoPointer(s string) error {
	return nil 
}

`

func TestFunction(t *testing.T) {
	p, err := ParsePackage("github.com/fzerorubigd/fixture")
	assert.NoError(t, err)

	f, err := p.FindFunction("test")
	assert.NoError(t, err)

	assert.Equal(t, "test", f.Name())
	assert.Equal(t, "// Multi line\n// comment\n// here", f.Docs().String())
	assert.Equal(t, p, f.Package())
	assert.Equal(t, "main.go", f.File().FileName())

	f, err = p.FindMethod("beta", "assert")
	assert.NoError(t, err)

	assert.Equal(t, "assert", f.Name())
	assert.Equal(t, "", f.Docs().String())
	assert.Equal(t, p, f.Package())
	assert.Equal(t, "main.go", f.File().FileName())
	assert.NotNil(t, f.Receiver())
	assert.Equal(t, "", f.Receiver().Name())

	f, err = p.FindMethod("alpha", "testing")
	assert.NoError(t, err)

	assert.Equal(t, "testing", f.Name())
	assert.Equal(t, "", f.Docs().String())
	assert.Equal(t, p, f.Package())
	assert.Equal(t, "main.go", f.File().FileName())
	assert.NotNil(t, f.Receiver())
	assert.Equal(t, "a", f.Receiver().Name())
}

func TestFunctionExtra(t *testing.T) {
	p := &Package{}
	f, err := ParseFile(testFunction, p)
	require.Nil(t, err)
	p.files = append(p.files, f)

	fn, err := p.FindFunction("Testing")
	require.NoError(t, err)

	assert.Equal(t, p, fn.Package())
	assert.Equal(t, f, fn.File())
	assert.Nil(t, fn.Receiver())
	assert.Empty(t, fn.ReceiverType())
	assert.False(t, fn.RecieverPointer())
	assert.Equal(t, "return nil", strings.Trim(fn.Body(), "\n\t "))

	fn, err = p.FindMethod("E", "Testing")
	require.NoError(t, err)

	assert.Equal(t, p, fn.Package())
	assert.Equal(t, f, fn.File())
	require.NotNil(t, fn.Receiver())
	assert.Equal(t, "e", fn.Receiver().Name())
	assert.Equal(t, "E", fn.ReceiverType())
	assert.True(t, fn.RecieverPointer())
	assert.Equal(t, "return nil", strings.Trim(fn.Body(), "\n\t "))

	fn, err = p.FindMethod("E", "NoPointer")
	require.NoError(t, err)

	assert.Equal(t, p, fn.Package())
	assert.Equal(t, f, fn.File())
	require.NotNil(t, fn.Receiver())
	assert.Equal(t, "e", fn.Receiver().Name())
	assert.Equal(t, "E", fn.ReceiverType())
	assert.False(t, fn.RecieverPointer())
	assert.Equal(t, "return nil", strings.Trim(fn.Body(), "\n\t "))

}
