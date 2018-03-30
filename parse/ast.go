package parse

import (
	"fmt"
	"strings"
)

// Query is the single query
type Query struct {
	Statement Statement
}

// Statement is the statement interface, only select is supported currently
type Statement interface {
	parse(*parser) error
}

// Field is the fields inside select and functions, for functions the parameters is valid
type Field struct {
	Table      string // table name if it is a table field
	Alias      string // alias if this has an alias (not supported yet)
	Item       Item   // Token Item, it can be alpha, string literal, number or bool
	Parameters Fields // Parameters is only valid for function fields
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

// FuncItem is an item with its parameter
type FuncItem interface {
	Item
	Parameters() Fields
}

// Order is one order in the order array
type Order struct {
	Field string // Field name
	Index int    // Index in the actual table
	DESC  bool   // order asc or desc
}

// Orders group of orders, the first order is more important than the next
type Orders []Order

// SelectStmt is the select query
type SelectStmt struct {
	Table  string // Which table tos elect from
	Fields Fields // Which field are requested

	Where Stack  // Where stack (need more work on where stack )
	Order Orders // Orders in order part, if any
	Start int    // the start column , -1 means no start specified
	Count int    // the count to show, -1 means no count specified
}

// GetTokenString is a simple function to handle the quoted strings, it remove the single or double quote
// from the string and remove escape sequence
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
