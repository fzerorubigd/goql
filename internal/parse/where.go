package parse

import "fmt"

const (
	whereAlpha = 1 << iota
	whereBool
	whereOp
)

var (
	precedence = map[itemType]int{
		itemOr:  -1,
		itemAnd: -1,
	}
)

func isOperand(t itemType) bool {
	return t == itemAnd ||
		t == itemOr ||
		t == itemLike ||
		t == itemEqual ||
		t == itemNotEqual ||
		t == itemGreater ||
		t == itemGreaterEqual ||
		t == itemLesser ||
		t == itemLesserEqual
}

func (p *parser) where() (Stack, error) {
	w := p.scanIgnoreWhiteSpace()
	if w.typ != itemWhere {
		return nil, fmt.Errorf("expect WHERE but got %s", w)
	}

	op := NewStack(10)
	final := NewStack(20)

	var expected = whereAlpha
bigLoop:
	for {
		switch ahead := p.scanIgnoreWhiteSpace(); {
		case ahead.typ == itemOrder || ahead.typ == itemLimit || ahead.typ == itemEOF:
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
			top, err := op.Peek()
			if err == nil && precedence[top.typ] > precedence[ahead.typ] {
				_, _ = op.Pop()
				final.Push(top)
			}
			op.Push(ahead)
			expected = whereBool | whereAlpha
			if ahead.typ == itemLike {
				expected = whereAlpha
			}
		case ahead.typ == itemNumber || ahead.typ == itemAlpha || ahead.typ == itemLiteral1 || ahead.typ == itemLiteral2:
			// operand
			if expected|whereAlpha != expected {
				return nil, fmt.Errorf("not expected operand but got %s", ahead)
			}
			final.Push(ahead)
			expected = whereOp
		}
	}

	for {
		o, err := op.Pop()
		if err != nil {
			break
		}
		final.Push(o)
	}
	// stack len function
	s := NewStack(0)
	for {
		o, err := final.Pop()
		if err != nil {
			break
		}
		s.Push(o)
	}
	return s, nil
}
