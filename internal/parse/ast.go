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

// Function is an sql function
type Function struct {
	Name       string
	Parameters Fields
}

// Field is the fields inside select
type Field struct {
	WildCard bool      // aka '*'
	String   string    // is this an string ('string')  empty means no
	Number   string    // is this an number? (19991) empty means no
	Function *Function // is this a function? nil means no
	Table    string    // the part before dot
	Column   string    // the column
	Alias    string    // alias of the column
}

// Fields is the collection of fields with order
type Fields []Field

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
