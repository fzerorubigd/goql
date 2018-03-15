package executor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fzerorubigd/goql/internal/parse"
	"github.com/fzerorubigd/goql/structures"
)

const (
	itemGetter   parse.ItemType = -9999
	nullValue    null           = 0
	notNullValue null           = 1
)

type (
	null       int
	getter     func([]structures.Valuer) interface{}
	opGetter   func(getter, getter) getter
	operGetter func(parse.Item) getter
)

var (
	operGetterMap = map[parse.ItemType]operGetter{
		itemColumn:         fieldGetterGenerator,
		parse.ItemAlpha:    alphaGetterGenerator,
		parse.ItemLiteral1: literal1GetterGenerator,
		parse.ItemNumber:   numberGetterGenerator,
		parse.ItemNull:     nullGetterGenerator,
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
		parse.ItemIs:           operIs,
		parse.ItemLike:         operLike,
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

func nullGetterGenerator(t parse.Item) getter {
	if t.Type() != parse.ItemNull {
		panic("runtime error")
	}
	return func(in []structures.Valuer) interface{} {
		return nullValue
	}
}

func alphaGetterGenerator(t parse.Item) getter {
	if t.Type() != parse.ItemAlpha {
		panic("runtime error")
	}
	switch strings.ToLower(t.Value()) {
	case "true":
		return func([]structures.Valuer) interface{} {
			return true
		}
	case "false":
		return func([]structures.Valuer) interface{} {
			return false
		}
	default:
		panic("runtime error")
	}
}

func fieldGetterGenerator(t parse.Item) getter {
	if t.Type() != itemColumn {
		panic("runtime error")
	}
	var idx = t.Pos()
	return func(in []structures.Valuer) interface{} {
		switch t := in[idx].(type) {
		case structures.String:
			if t.Null {
				return nullValue
			}
			return t.String
		case structures.Bool:
			if t.Null {
				return nullValue
			}
			return t.Bool
		case structures.Number:
			if t.Null {
				return nullValue
			}
			return t.Number

		}
		return in[idx]
	}
}

func literal1GetterGenerator(t parse.Item) getter {
	if t.Type() != parse.ItemLiteral1 {
		panic("runtime error")
	}
	var v = parse.GetTokenString(t)
	return func([]structures.Valuer) interface{} {
		return v
	}
}

func numberGetterGenerator(t parse.Item) getter {
	if t.Type() != parse.ItemNumber {
		panic("runtime error")
	}
	v, _ := strconv.ParseFloat(t.Value(), 10)
	return func([]structures.Valuer) interface{} {
		return v
	}
}

func equal(l getter, r getter) getter {
	return func(in []structures.Valuer) interface{} {
		lv := l(in)
		rv := r(in)
		return lv == castAsLeft(lv, rv)
	}
}

func notEqual(l getter, r getter) getter {
	return func(in []structures.Valuer) interface{} {
		lv := l(in)
		rv := r(in)
		return lv != castAsLeft(lv, rv)
	}
}

func operOr(l getter, r getter) getter {
	return func(in []structures.Valuer) interface{} {
		return toBool(l(in)) || toBool(r(in))
	}
}

func operAnd(l getter, r getter) getter {
	return func(in []structures.Valuer) interface{} {
		return toBool(l(in)) && toBool(r(in))
	}
}

func operGreater(l getter, r getter) getter {
	return func(in []structures.Valuer) interface{} {
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
		case null:
			return lv.(null) > rv.(null)
		}
		panic("not supported type")
	}
}

func operGreaterEqual(l getter, r getter) getter {
	return func(in []structures.Valuer) interface{} {
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
		case null:
			return lv.(null) >= rv.(null)
		}
		panic("not supported type")
	}
}

func operLesser(l getter, r getter) getter {
	return func(in []structures.Valuer) interface{} {
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
		case null:
			return lv.(null) < rv.(null)
		}
		panic("not supported type")
	}
}

func operLesserEqual(l getter, r getter) getter {
	return func(in []structures.Valuer) interface{} {
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
		case null:
			return lv.(null) <= rv.(null)
		}
		panic("not supported type")
	}
}

func operIs(l getter, r getter) getter {
	// TODO : there is a problem, if he right operand is a column with null value it works here :)
	return func(in []structures.Valuer) interface{} {
		v := r(in)
		n, ok := v.(null)
		if !ok {
			panic("is only works on null/not null")
		}
		return toNull(l(in)) == n
	}
}

func operLike(l getter, r getter) getter {
	return func(in []structures.Valuer) interface{} {
		re := likeStr(toString(r(in))).regexp()
		v := toString(l(in))
		return re.MatchString(v)
	}
}

func operNot(l getter) getter {
	return func(in []structures.Valuer) interface{} {
		d := l(in)
		switch t := d.(type) {
		case bool:
			return !t
		case null:
			if t == notNullValue {
				return nullValue
			}
			return notNullValue
		case string, float64:
			panic("not is not applicable for string or number")
		}
		return !toBool(l(in))
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
	case null:
		return toNull(r)
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
	case null:
		return false
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
	case null:
		return 0
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
	case null:
		return ""
	}
	panic(fmt.Sprintf("result from type %T", in))
}

func toNull(in interface{}) null {
	switch t := in.(type) {
	case bool:
		return notNullValue
	case string:
		return notNullValue
	case float64:
		return notNullValue
	case null:
		return t
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
		t == parse.ItemLesserEqual ||
		t == parse.ItemIs
}

func getGetter(t parse.Item) getter {
	if g, ok := t.(getter); ok {
		return g
	}

	m, ok := operGetterMap[t.Type()]
	if !ok {
		panic(fmt.Sprintf("%+v is not belong here", t))
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
		return func([]structures.Valuer) interface{} {
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
		} else if t.Type() == parse.ItemNot {
			// pop the operand
			ts, err := p.Pop()
			if err != nil {
				return nil, fmt.Errorf("end of stack")
			}
			// push back the last getter (but after not )
			w.Push(operNot(getGetter(ts)))
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
