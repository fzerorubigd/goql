package astdata

import (
	"go/ast"
)

// ChannelType is the channel type in go source code
type ChannelType struct {
	pkg *Package

	Direction ast.ChanDir
	Type      Definition
}

// String represent string version of the data
func (c *ChannelType) String() string {
	switch c.Direction {
	case ast.SEND:
		return "chan<- " + c.Type.String()
	case ast.RECV:
		return "<-chan " + c.Type.String()
	default:
		return "chan " + c.Type.String()
	}
}

// Package is the package of channel
func (c *ChannelType) Package() *Package {
	return c.pkg
}

func getChannel(p *Package, f *File, t *ast.ChanType) Definition {
	return &ChannelType{
		pkg:       p,
		Direction: t.Dir,
		Type:      newType(p, f, t.Value),
	}
}
