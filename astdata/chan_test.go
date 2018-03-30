package astdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testChan = `
package example

func XX (in chan int, out <-chan string ) chan<- float64 {

}
`

func TestChanType(t *testing.T) {
	p := &Package{}

	f, err := ParseFile(testChan, p)
	assert.NoError(t, err)

	p.files = append(p.files, f)

	fn, err := p.FindFunction("XX")
	assert.NoError(t, err)
	assert.IsType(t, &FuncType{}, fn.Definition())

	assert.Equal(t, 2, len(fn.def.parameters))
	assert.Equal(t, 1, len(fn.def.results))

	ch1 := fn.def.parameters[0].def
	assert.IsType(t, &ChannelType{}, ch1)
	ch1c := ch1.(*ChannelType)
	assert.Equal(t, p, ch1c.Package())
	assert.Equal(t, f, ch1c.File())
	assert.Equal(t, SEND|RECV, ch1c.Direction())
	assert.IsType(t, &IdentType{}, ch1c.ValueDefinition())

	assert.Equal(t, "chan int", ch1c.String())

	ch1 = fn.def.parameters[1].def
	assert.IsType(t, &ChannelType{}, ch1)
	ch1c = ch1.(*ChannelType)
	assert.Equal(t, p, ch1c.Package())
	assert.Equal(t, f, ch1c.File())
	assert.Equal(t, RECV, ch1c.Direction())
	assert.IsType(t, &IdentType{}, ch1c.ValueDefinition())

	assert.Equal(t, "<-chan string", ch1c.String())

	ch1 = fn.def.results[0].def
	assert.IsType(t, &ChannelType{}, ch1)
	ch1c = ch1.(*ChannelType)
	assert.Equal(t, p, ch1c.Package())
	assert.Equal(t, f, ch1c.File())
	assert.Equal(t, SEND, ch1c.Direction())
	assert.IsType(t, &IdentType{}, ch1c.ValueDefinition())

	assert.Equal(t, "chan<- float64", ch1c.String())

}
