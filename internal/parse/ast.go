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
}

func getTokenString(t item) string {
	v := t.value
	if t.typ == itemLiteral1 {
		v = strings.Trim(strings.Replace(t.value, `\'`, `'`, -1), "'")
	}
	if t.typ == itemLiteral2 {
		v = strings.Trim(strings.Replace(t.value, `\"`, `"`, -1), "\"")
	}
	return v
}

func (ss *SelectStmt) parseField(p *parser) (Field, error) {
	token := p.scanIgnoreWhiteSpace()
	if token.typ == itemWildCard {
		return Field{WildCard: true}, nil
	}

	if token.typ == itemAlpha || token.typ == itemLiteral2 {
		ahead := p.scan() // white space is not allowed here
		if ahead.typ != itemDot {
			p.reject()
			return Field{Column: getTokenString(token)}, nil
		}
		ahead = p.scan()
		if ahead.typ == itemAlpha || ahead.typ == itemLiteral2 {
			return Field{
				Table:  getTokenString(token),
				Column: getTokenString(ahead),
			}, nil
		}
	}

	return Field{}, fmt.Errorf("unexpected token, %s", token)
}

func (ss *SelectStmt) parseFields(p *parser) error {
	for {
		field, err := ss.parseField(p)
		if err != nil {
			return err
		}
		ss.Fields = append(ss.Fields, field)

		comma := p.scanIgnoreWhiteSpace()
		if comma.typ != itemComma {
			p.reject()
			break
		}
	}
	return nil
}

func (ss *SelectStmt) parse(p *parser) error {
	if err := ss.parseFields(p); err != nil {
		return err
	}

	t := p.scanIgnoreWhiteSpace()
	// must be from
	if t.typ != itemFrom {
		return fmt.Errorf("unexpected %s , expected FROM or COMMA (,)", t)
	}

	t = p.scanIgnoreWhiteSpace()
	if t.typ != itemAlpha && t.typ != itemLiteral2 {
		return fmt.Errorf("unexpected input %s , need table name", t)
	}
	ss.Table = getTokenString(t)

	return nil
}

func newStatement(p *parser) (Statement, error) {
	start := p.scan()
	switch start.typ {
	case itemSelect:
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
