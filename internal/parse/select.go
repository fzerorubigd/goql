package parse

import (
	"fmt"
	"strconv"
)

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

func isOperator(t ItemType) bool {
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

func isOperand(t ItemType) bool {
	return t == ItemNumber ||
		t == ItemAlpha ||
		t == ItemLiteral1 ||
		t == ItemLiteral2 ||
		t == ItemNot ||
		t == ItemNull
}

func isKeyword(t ItemType) bool {
	return t == ItemOrder || t == ItemLimit || t == ItemEOF
}

func (ss *SelectStmt) parseField(p *parser) (Field, error) {
	token := p.scanIgnoreWhiteSpace()
	if token.typ == ItemWildCard {
		return Field{WildCard: true}, nil
	}

	if token.typ == ItemAlpha {
		ahead := p.scanIgnoreWhiteSpace()
		if ahead.typ == ItemParenOpen {
			// function
			fileds, err := ss.parseFields(p, true)
			if err != nil {
				return Field{}, err
			}
			ahead = p.scanIgnoreWhiteSpace()
			if ahead.typ != ItemParenClose {
				return Field{}, fmt.Errorf("expected ) but got %s", ahead)
			}
			return Field{
				Function: &Function{
					Name:       GetTokenString(token),
					Parameters: fileds,
				},
			}, nil
		}
		p.reject()
	}

	if token.typ == ItemAlpha || token.typ == ItemLiteral2 {
		ahead := p.scan() // white space is not allowed here
		if ahead.typ == ItemDot {
			ahead = p.scan()
			if ahead.typ == ItemAlpha || ahead.typ == ItemLiteral2 {
				return Field{
					Table:  GetTokenString(token),
					Column: GetTokenString(ahead),
				}, nil
			}
			return Field{}, fmt.Errorf("expected field name got %s", ahead)
		}
		p.reject() // nope its not a dot
		return Field{
			Column: GetTokenString(token),
		}, nil
	}

	if token.typ == ItemNumber {
		return Field{
			Number: token.Value(),
		}, nil
	}

	if token.typ == ItemLiteral1 {
		return Field{
			String: GetTokenString(token),
		}, nil
	}

	return Field{}, fmt.Errorf("unexpected token, %s", token)
}

func (ss *SelectStmt) parseFields(p *parser, fn bool) (Fields, error) {
	var res Fields
	if fn {
		ahead := p.scanIgnoreWhiteSpace()
		if ahead.typ == ItemParenClose {
			p.reject()
			return Fields{}, nil
		}
		p.reject()
	}
	for {
		field, err := ss.parseField(p)
		if err != nil {
			return nil, err
		}
		res = append(res, field)

		comma := p.scanIgnoreWhiteSpace()
		if comma.typ != ItemComma {
			p.reject()
			break
		}
	}
	return res, nil
}

func (ss *SelectStmt) parse(p *parser) error {
	var err error
	ss.Fields, err = ss.parseFields(p, false)
	if err != nil {
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
	} else {
		p.reject()
	}

	if w := p.scanIgnoreWhiteSpace(); w.typ == ItemOrder {
		p.reject()
		var err error
		ss.Order, err = p.order()
		if err != nil {
			return err
		}
	} else {
		p.reject()
	}

	ss.Start, ss.Count = -1, -1
	if w := p.scanIgnoreWhiteSpace(); w.typ == ItemLimit {
		p.reject()
		var err error
		ss.Start, ss.Count, err = p.limit()
		if err != nil {
			return err
		}
	} else {
		p.reject()
	}

	if w := p.scanIgnoreWhiteSpace(); w.typ != ItemEOF {
		return fmt.Errorf("unexpected token %s", w)
	}

	return nil
}

func (p *parser) operator(ahead item, op, final Stack, expected int) (int, error) {
	// operator
	if expected|whereOp != expected {
		return 0, fmt.Errorf("not expected operator but got %s", ahead)
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
	return whereAlpha, nil
}

func (p *parser) operand(ahead item, op, final Stack, expected int) (int, error) {
	// operand
	if expected|whereAlpha != expected && expected|whereStart != expected {
		return 0, fmt.Errorf("not expected operand but got %s", ahead)
	}
	not := expected|whereNotOp == expected
	if ahead.typ == ItemNot {
		if not {
			return 0, fmt.Errorf("not after not")
		}
		op.Push(ahead)
		expected = whereAlpha | whereNotOp
	} else {
		final.Push(ahead)
		expected = whereOp
	}

	if not {
		top, err := op.Pop()
		assertTrue(err == nil && top.Type() == ItemNot, "why")
		final.Push(top)
	}
	return expected, nil
}

func (p *parser) parenClose(ahead item, op, final Stack, expected int) (int, error) {
	if expected|whereOp != expected {
		return 0, fmt.Errorf("wrong ')' ")
	}

	for {
		o, err := op.Pop()
		assertTrue(err == nil, "why no op in stack?")
		if o.Type() == ItemParenOpen {
			break
		}
		final.Push(o)
	}
	return expected, nil

}

func (p *parser) parenOpen(ahead item, op, final Stack, expected int) (int, error) {
	if expected|whereStart != expected && expected|whereAlpha != expected {
		return 0, fmt.Errorf("wrong '(' ")
	}
	op.Push(ahead)
	return expected, nil
}

func (p *parser) copyStack(op, final Stack) {
	for {
		o, err := op.Pop()
		if err != nil {
			break
		}
		final.Push(o)
	}
}

func (p *parser) where() (Stack, error) {
	w := p.scanIgnoreWhiteSpace()
	assertType(w, ItemWhere)

	op := NewStack(0)
	final := NewStack(0)

	var expected = whereStart
bigLoop:
	for {
		switch ahead := p.scanIgnoreWhiteSpace(); {
		case isKeyword(ahead.typ):
			if expected|whereOp != expected {
				return nil, fmt.Errorf("expected an operand but end of where %d", expected)
			}
			p.reject()
			break bigLoop
		case isOperator(ahead.typ):
			var err error
			expected, err = p.operator(ahead, op, final, expected)
			if err != nil {
				return nil, err
			}
		case isOperand(ahead.typ):
			var err error
			expected, err = p.operand(ahead, op, final, expected)
			if err != nil {
				return nil, err
			}

		case ahead.typ == ItemParenOpen:
			var err error
			expected, err = p.parenOpen(ahead, op, final, expected)
			if err != nil {
				return nil, err
			}
		case ahead.typ == ItemParenClose:
			var err error
			expected, err = p.parenClose(ahead, op, final, expected)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("not expected %s", ahead)
		}
	}
	p.copyStack(op, final)
	return final, nil
}

func (p *parser) order() (Orders, error) {
	var res Orders
	w := p.scanIgnoreWhiteSpace()
	assertType(w, ItemOrder)

	if w := p.scanIgnoreWhiteSpace(); w.typ != ItemBy {
		return nil, fmt.Errorf("invalid token, need by after order , got %s", w)
	}

	for {
		w := p.scanIgnoreWhiteSpace()
		if w.typ != ItemAlpha && w.typ != ItemLiteral2 {
			return nil, fmt.Errorf("need column name got %s", w)
		}
		val := GetTokenString(w)
		res = append(res, Order{
			Field: val,
		})
		ahead := p.scanIgnoreWhiteSpace()
		if ahead.typ == ItemAsc || ahead.typ == ItemDesc {
			if ahead.typ == ItemDesc {
				res[len(res)-1].DESC = true
			}
			ahead = p.scanIgnoreWhiteSpace()
		}

		if ahead.typ != ItemComma {
			p.reject()
			break
		}
	}
	return res, nil
}

func (p *parser) limit() (int, int, error) {
	var start, count int64
	var err error
	w := p.scanIgnoreWhiteSpace()
	assertType(w, ItemLimit)

	w = p.scanIgnoreWhiteSpace()
	if w.typ != ItemNumber {
		return 0, 0, fmt.Errorf("limit need a number but got %s", w)
	}

	start, err = strconv.ParseInt(w.value, 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("error on converting string to int : %s", err)
	}

	w = p.scanIgnoreWhiteSpace()
	if w.typ != ItemComma {
		p.reject()
		// one number means count
		return 0, int(start), nil
	}

	w = p.scanIgnoreWhiteSpace()
	if w.typ != ItemNumber {
		return 0, 0, fmt.Errorf("need the second limit number got %s", w)
	}
	count, err = strconv.ParseInt(w.value, 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("error on converting string to int : %s", err)
	}

	return int(start), int(count), nil
}
