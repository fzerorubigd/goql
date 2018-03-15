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

// Field is the fields inside select
type Field struct {
	WildCard bool // aka '*'
	//Alias    string // add support for this? :))
	Table  string // the part before dot
	Column string // the column
}

// Fields is the collection of fields with order
type Fields []Field

// SelectStmt is the select query
type SelectStmt struct {
	Table  string
	Fields Fields

	Where Stack
}

// GetTokenString is a simple function to handle the quoted strings
func GetTokenString(t Item) string {
	v := t.Value()
	if t.Type() == ItemLiteral1 {
		v = strings.Trim(strings.Replace(t.Value(), `\'`, `'`, -1), "'")
	}
	if t.Type() == ItemLiteral2 {
		v = strings.Trim(strings.Replace(t.Value(), `\"`, `"`, -1), "\"")
	}
	return v
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
