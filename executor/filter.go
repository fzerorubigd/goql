package executor

import (
	"fmt"
	"strconv"
	"strings"

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
		parse.ItemEqual:        equal,
		parse.ItemNotEqual:     notEqual,
		parse.ItemOr:           operOr,
		parse.ItemAnd:          operAnd,
		parse.ItemGreater:      operGreater,
		parse.ItemGreaterEqual: operGreaterEqual,
		parse.ItemLesser:       operLesser,
		parse.ItemLesserEqual:  operLesserEqual,
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
		lv := l(in)
		rv := r(in)
		return lv == castAsLeft(lv, rv)
	}
}

func notEqual(l getter, r getter) getter {
	return func(in []interface{}) interface{} {
		lv := l(in)
		rv := r(in)
		return lv != castAsLeft(lv, rv)
	}
}

func operOr(l getter, r getter) getter {
	return func(in []interface{}) interface{} {
		return toBool(l(in)) || toBool(r(in))
	}
}

func operAnd(l getter, r getter) getter {
	return func(in []interface{}) interface{} {
		return toBool(l(in)) && toBool(r(in))
	}
}

func operGreater(l getter, r getter) getter {
	return func(in []interface{}) interface{} {
		lv := l(in)
		rv := castAsLeft(lv, r(in))
		switch lv.(type) {
		case bool:
			if lv.(bool) != rv.(bool) {
				return lv.(bool)
			}
			return false
		case string:
			return strings.Compare(lv.(string), rv.(string)) > 0
		case float64:
			return lv.(float64) > rv.(float64)
		}
		panic("not supported type")
	}
}

func operGreaterEqual(l getter, r getter) getter {
	return func(in []interface{}) interface{} {
		lv := l(in)
		rv := castAsLeft(lv, r(in))
		switch lv.(type) {
		case bool:
			if lv.(bool) != rv.(bool) {
				return lv.(bool)
			}
			return true
		case string:
			return strings.Compare(lv.(string), rv.(string)) >= 0
		case float64:
			return lv.(float64) >= rv.(float64)
		}
		panic("not supported type")
	}
}

func operLesser(l getter, r getter) getter {
	return func(in []interface{}) (b interface{}) {
		lv := l(in)
		rv := castAsLeft(lv, r(in))
		switch lv.(type) {
		case bool:
			if lv.(bool) != rv.(bool) {
				return !lv.(bool)
			}
			return false
		case string:
			return strings.Compare(lv.(string), rv.(string)) < 0
		case float64:
			return lv.(float64) < rv.(float64)
		}
		panic("not supported type")
	}
}

func operLesserEqual(l getter, r getter) getter {
	return func(in []interface{}) interface{} {
		lv := l(in)
		rv := castAsLeft(lv, r(in))
		switch lv.(type) {
		case bool:
			if lv.(bool) != rv.(bool) {
				return !lv.(bool)
			}
			return true
		case string:
			return strings.Compare(lv.(string), rv.(string)) <= 0
		case float64:
			return lv.(float64) <= rv.(float64)
		}
		panic("not supported type")
	}
}

func castAsLeft(l, r interface{}) interface{} {
	switch l.(type) {
	case bool:
		return toBool(r)
	case float64:
		return toNumber(r)
	case string:
		return toString(r)
	}
	panic(fmt.Sprintf("%T is invalid type", l))
}

func toBool(in interface{}) bool {
	switch t := in.(type) {
	case bool:
		return t
	case string:
		b, _ := strconv.ParseBool(t)
		return b
	case float64:
		return t != 0
	}
	panic(fmt.Sprintf("result from type %T", in))
}

func toNumber(in interface{}) float64 {
	switch t := in.(type) {
	case bool:
		if t {
			return 1
		}
		return 0
	case string:
		f, _ := strconv.ParseFloat(t, 64)
		return f
	case float64:
		return t
	}
	panic(fmt.Sprintf("result from type %T", in))
}

func toString(in interface{}) string {
	switch t := in.(type) {
	case bool:
		return fmt.Sprint(t)
	case string:
		return t
	case float64:
		return fmt.Sprint(t)
	}
	panic(fmt.Sprintf("result from type %T", in))
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
			g, ok := tt.(getter)
			if !ok {
				g = getGetter(tt)
			}
			if _, err := p.Pop(); err == nil {
				return nil, fmt.Errorf("two operand no op?")
			}
			return g, nil
		}
	}
}
