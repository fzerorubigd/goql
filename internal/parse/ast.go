package parse

import (
	"fmt"
	"strings"
)

// Query is the single query ast
type Query struct {
	Statement Statement
}

// Statement is the statement interface
type Statement interface {
	parse(*parser) error
}

// Field is the fields inside select and functions, for functions the parameters is valid
type Field struct {
	Table      string
	Alias      string
	Item       Item
	Parameters Fields
}

// Fields is the collection of fields with order
type Fields []Field

type itemFn struct {
	item
	parameter Fields
}

func (i itemFn) Parameters() Fields {
	return i.parameter
}

// FuncItem is an item with its parameter for where
type FuncItem interface {
	Item
	Parameters() Fields
}

// Order is one order in the order array
type Order struct {
	Field string
	Index int
	DESC  bool
}

// Orders group of orders
type Orders []Order

// SelectStmt is the select query
type SelectStmt struct {
	Table  string
	Fields Fields

	Where Stack
	Order Orders
	Start int
	Count int
}

// GetTokenString is a simple function to handle the quoted strings
func GetTokenString(t Item) string {
	v := t.Value()
	var l string
	var c byte
	switch t.Type() {
	case ItemLiteral1:
		l = `'`
		c = '\''
	case ItemLiteral2:
		l = `"`
		c = '"'
	default:
		return v
	}
	v = strings.Replace(t.Value(), `\`+l, l, -1)

	if len(v) < 2 || v[0] != c || v[len(v)-1] != c {
		panic("un-terminated literal")
	}

	return v[1 : len(v)-1]
}

func newStatement(p *parser) (Statement, error) {
	start := p.scan()
	switch start.typ {
	case ItemSelect:
		sel := &SelectStmt{}
		err := sel.parse(p)
		if err != nil {
			return nil, err
		}
		return sel, nil
	default:
		return nil, fmt.Errorf("token %s is not a valid token", start.value)
	}
}
