package parse

import "fmt"

const (
	whereStart = 1 << iota
	whereAlpha
	whereOp
	whereNotOp
)

var (
	precedence = map[ItemType]int{
		ItemOr:  -1,
		ItemAnd: -1,
	}
)

func isOperand(t ItemType) bool {
	return t == ItemAnd ||
		t == ItemOr ||
		t == ItemLike ||
		t == ItemEqual ||
		t == ItemNotEqual ||
		t == ItemGreater ||
		t == ItemGreaterEqual ||
		t == ItemLesser ||
		t == ItemLesserEqual ||
		t == ItemIs
}

func (ss *SelectStmt) parseField(p *parser) (Field, error) {
	token := p.scanIgnoreWhiteSpace()
	if token.typ == ItemWildCard {
		return Field{WildCard: true}, nil
	}

	if token.typ == ItemAlpha || token.typ == ItemLiteral2 {
		ahead := p.scan() // white space is not allowed here
		if ahead.typ != ItemDot {
			p.reject()
			return Field{Column: GetTokenString(token)}, nil
		}
		ahead = p.scan()
		if ahead.typ == ItemAlpha || ahead.typ == ItemLiteral2 {
			return Field{
				Table:  GetTokenString(token),
				Column: GetTokenString(ahead),
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
		if comma.typ != ItemComma {
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
	if t.typ != ItemFrom {
		return fmt.Errorf("unexpected %s , expected FROM or COMMA (,)", t)
	}

	t = p.scanIgnoreWhiteSpace()
	if t.typ != ItemAlpha && t.typ != ItemLiteral2 {
		return fmt.Errorf("unexpected input %s , need table name", t)
	}
	ss.Table = GetTokenString(t)

	if w := p.scanIgnoreWhiteSpace(); w.typ == ItemWhere {
		p.reject()
		var err error
		ss.Where, err = p.where()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *parser) where() (Stack, error) {
	w := p.scanIgnoreWhiteSpace()
	if w.typ != ItemWhere {
		return nil, fmt.Errorf("expect WHERE but got %s", w)
	}

	op := NewStack(0)
	final := NewStack(0)

	var expected = whereStart
bigLoop:
	for {
		switch ahead := p.scanIgnoreWhiteSpace(); {
		case ahead.typ == ItemOrder || ahead.typ == ItemLimit || ahead.typ == ItemEOF:
			if expected|whereOp != expected {
				return nil, fmt.Errorf("expected an operand but end of where %d", expected)
			}
			p.reject()
			break bigLoop
		case isOperand(ahead.typ):
			// operator
			if expected|whereOp != expected {
				return nil, fmt.Errorf("not expected operator but got %s", ahead)
			}
			// so we got the op, push it into stack
			for {
				top, err := op.Peek()
				if err == nil && top.Type() != ItemParenOpen && precedence[top.Type()] > precedence[ahead.Type()] {
					_, _ = op.Pop()
					final.Push(top)
				} else {
					break
				}
			}
			op.Push(ahead)
			expected = whereAlpha
		case ahead.typ == ItemNumber || ahead.typ == ItemAlpha || ahead.typ == ItemLiteral1 || ahead.typ == ItemLiteral2 || ahead.typ == ItemNot || ahead.typ == ItemNull:
			// operand
			if expected|whereAlpha != expected && expected|whereStart != expected {
				return nil, fmt.Errorf("not expected operand but got %s", ahead)
			}
			not := expected|whereNotOp == expected
			if ahead.typ == ItemNot {
				if not {
					return nil, fmt.Errorf("not after not")
				}
				op.Push(ahead)
				expected = whereAlpha | whereNotOp
			} else {
				final.Push(ahead)
				expected = whereOp
			}

			if not {
				top, err := op.Pop()
				if err != nil || top.Type() != ItemNot {
					panic("why")
				}
				final.Push(top)
			}

		case ahead.typ == ItemParenOpen:
			if expected|whereStart != expected && expected|whereAlpha != expected {
				return nil, fmt.Errorf("wrong '(' ")
			}
			op.Push(ahead)
		case ahead.typ == ItemParenClose:
			if expected|whereOp != expected {
				return nil, fmt.Errorf("wrong ')' ")
			}

			for {
				o, err := op.Pop()
				if err != nil {
					return nil, fmt.Errorf("invalid ')")
				}
				if o.Type() == ItemParenOpen {
					break
				}
				final.Push(o)
			}
		default:
			return nil, fmt.Errorf("not expected %s", ahead)
		}
	}

	for {
		o, err := op.Pop()
		if err != nil {
			break
		}
		final.Push(o)
	}
	return final, nil
}
