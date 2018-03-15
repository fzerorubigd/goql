package parse

import "fmt"

const (
	whereStart = 1 << iota
	whereAlpha
	whereBool
	whereOp
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
		t == ItemLesserEqual
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
				return nil, fmt.Errorf("expected an operand but end of where")
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
			expected = whereBool | whereAlpha
			if ahead.typ == ItemLike {
				expected = whereAlpha
			}
		case ahead.typ == ItemNumber || ahead.typ == ItemAlpha || ahead.typ == ItemLiteral1 || ahead.typ == ItemLiteral2:
			// operand
			if expected|whereAlpha != expected && expected|whereStart != expected {
				return nil, fmt.Errorf("not expected operand but got %s", ahead)
			}
			final.Push(ahead)
			expected = whereOp
		case ahead.typ == ItemParenOpen:
			if expected|whereStart != expected && expected|whereAlpha != expected {
				return nil, fmt.Errorf("wrong '(' ")
			}
			op.Push(ahead)
		case ahead.typ == ItemParenClose:
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
