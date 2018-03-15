package executor

import (
	"fmt"
	"strconv"

	"github.com/fzerorubigd/goql/internal/parse"
)

const (
	itemGetter parse.ItemType = -9999
)

type (
	getter     func([]interface{}) interface{}
	opGetter   func(getter, getter) getter
	operGetter func(parse.Item) getter
)

var (
	operGetterMap = map[parse.ItemType]operGetter{
		itemColumn:         fieldGetterGenerator,
		parse.ItemLiteral1: literal1GetterGenerator,
		parse.ItemNumber:   numberGetterGenerator,
	}

	opGetterMap = map[parse.ItemType]opGetter{
		parse.ItemEqual:    equal,
		parse.ItemNotEqual: notEqual,
	}
)

func (g getter) Type() parse.ItemType {
	return itemGetter
}

func (g getter) Pos() int {
	return 0
}

func (g getter) Value() string {
	return ""
}

func (g getter) String() string {
	return ""
}

func fieldGetterGenerator(t parse.Item) getter {
	if t.Type() != itemColumn {
		panic("runtime error")
	}
	var idx = t.Pos()
	return func(in []interface{}) interface{} {
		return in[idx]
	}
}

func literal1GetterGenerator(t parse.Item) getter {
	if t.Type() != parse.ItemLiteral1 {
		panic("runtime error")
	}
	var v = parse.GetTokenString(t)
	return func([]interface{}) interface{} {
		return v
	}
}

func numberGetterGenerator(t parse.Item) getter {
	if t.Type() != parse.ItemNumber {
		panic("runtime error")
	}
	v, _ := strconv.ParseFloat(t.Value(), 10)
	return func([]interface{}) interface{} {
		return v
	}
}

func equal(l getter, r getter) getter {
	return func(in []interface{}) interface{} {
		return l(in) == r(in)
	}
}

func notEqual(l getter, r getter) getter {
	return func(in []interface{}) interface{} {
		return l(in) != r(in)
	}
}

func isOperator(t parse.ItemType) bool {
	return t == parse.ItemAnd ||
		t == parse.ItemOr ||
		t == parse.ItemLike ||
		t == parse.ItemEqual ||
		t == parse.ItemNotEqual ||
		t == parse.ItemGreater ||
		t == parse.ItemGreaterEqual ||
		t == parse.ItemLesser ||
		t == parse.ItemLesserEqual
}

func getGetter(t parse.Item) getter {
	if g, ok := t.(getter); ok {
		return g
	}
	m, ok := operGetterMap[t.Type()]
	if !ok {
		panic(fmt.Sprintf("%T is not belong here", t))
	}
	return m(t)
}

func getOpGetter(op parse.Item, lg, rg getter) getter {
	m, ok := opGetterMap[op.Type()]
	if !ok {
		panic("not implemented?")
	}
	return m(lg, rg)
}

func buildFilter(w parse.Stack) (getter, error) {
	var (
		p = parse.NewStack(0)
	)
	t, err := w.Pop()
	if err != nil {
		return func([]interface{}) interface{} {
			return true
		}, nil
	}
	for {
		if isOperator(t.Type()) {
			r, err := p.Pop()
			if err != nil {
				return nil, fmt.Errorf("invalid operand")
			}
			rg := getGetter(r)
			l, err := p.Pop()
			if err != nil {
				return nil, fmt.Errorf("invalid operand")
			}
			lg := getGetter(l)
			g := getOpGetter(t, lg, rg)
			w.Push(g)
		} else {
			p.Push(t)
		}

		t, err = w.Pop()
		if err != nil {
			tt, err := p.Pop()
			if err != nil {
				return nil, fmt.Errorf("its wrong")
			}
			g := tt.(getter)
			if _, err := p.Pop(); err == nil {
				return nil, fmt.Errorf("two operand no op?")
			}
			return g, nil
		}
	}
}
