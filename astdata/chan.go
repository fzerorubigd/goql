package astdata

import (
	"go/ast"
)

// ChanDir The direction of a channel type is indicated by a bit
// mask including one or both of the following constants.
type ChanDir int

const (
	// SEND is send only channel
	SEND ChanDir = 1 << iota
	// RECV is the receive only channel
	RECV
)

// ChannelType is the channel type in go source code
type ChannelType struct {
	embededData

	direction ChanDir
	def       Definition
}

// String represent string version of the data
func (c *ChannelType) String() string {
	// its bitwise :)) i should reconsider it
	switch c.direction {
	case SEND:
		return "chan<- " + c.def.String()
	case RECV:
		return "<-chan " + c.def.String()
	default:
		return "chan " + c.def.String()
	}
}

// Direction return the channel direction
func (c *ChannelType) Direction() ChanDir {
	return c.direction
}

// ValueDefinition return the definition of the type in channel
func (c *ChannelType) ValueDefinition() Definition {
	return c.def
}

// Compare try to compare this to def
func (c *ChannelType) Compare(def Definition) bool {
	return c.String() == def.String()
}

func getChannel(p *Package, f *File, t *ast.ChanType) Definition {
	return &ChannelType{
		embededData: embededData{
			pkg:  p,
			fl:   f,
			node: t,
		},
		direction: ChanDir(t.Dir),
		def:       newType(p, f, t.Value),
	}
}
