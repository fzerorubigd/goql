package goql

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/fzerorubigd/goql/parse"
)

const (
	itemGetter   parse.ItemType = -9999
	nullValue    null           = 0
	notNullValue null           = 1
)

type (
	null       int
	getter     func([]Getter) interface{}
	opGetter   func(getter, getter) getter
	operGetter func(parse.Item) getter
)

var (
	operGetterMap = map[parse.ItemType]operGetter{
		itemColumn:         fieldGetterGenerator,
		parse.ItemLiteral1: literal1GetterGenerator,
		parse.ItemNumber:   numberGetterGenerator,
		parse.ItemNull:     nullGetterGenerator,
		parse.ItemTrue:     boolGetterGenerator,
		parse.ItemFalse:    boolGetterGenerator,
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

func (g getter) Data() int {
	return 0
}

func (g getter) String() string {
	return ""
}

func assertType(t parse.Item, tp ...parse.ItemType) {
	for i := range tp {
		if tp[i] == t.Type() {
			return
		}
	}
	panic("runtime error")
}

func nullGetterGenerator(t parse.Item) getter {
	assertType(t, parse.ItemNull)
	return func(in []Getter) interface{} {
		return nullValue
	}
}

func boolGetterGenerator(t parse.Item) getter {
	assertType(t, parse.ItemTrue, parse.ItemFalse)
	res := t.Type() == parse.ItemTrue
	return func(in []Getter) interface{} {
		return res
	}
}

func fieldGetterGenerator(t parse.Item) getter {
	assertType(t, itemColumn)
	var idx = t.Data()
	return func(in []Getter) interface{} {
		switch t := in[idx].(type) {
		case String:
			if t.Null {
				return nullValue
			}
			return t.String
		case Bool:
			if t.Null {
				return nullValue
			}
			return t.Bool
		case Number:
			if t.Null {
				return nullValue
			}
			return t.Number
		case Definition:
			return t.Definition
		}
		return in[idx]
	}
}

func literal1GetterGenerator(t parse.Item) getter {
	assertType(t, parse.ItemLiteral1)
	var v = parse.GetTokenString(t)
	return func([]Getter) interface{} {
		return v
	}
}

func numberGetterGenerator(t parse.Item) getter {
	assertType(t, parse.ItemNumber)
	v, _ := strconv.ParseFloat(t.Value(), 10)
	return func([]Getter) interface{} {
		return v
	}
}

func equal(l getter, r getter) getter {
	return func(in []Getter) interface{} {
		lv := l(in)
		rv := r(in)
		switch t := lv.(type) {
		case astdata.Definition:
			def := toDefinition(rv)
			if t == nil || def == nil {
				return t == def
			}
			return t.Compare(toDefinition(rv))
		default:
			return lv == castAsLeft(lv, rv)
		}
	}
}

func notEqual(l getter, r getter) getter {
	return func(in []Getter) interface{} {
		lv := l(in)
		rv := r(in)
		switch t := lv.(type) {
		case astdata.Definition:
			def := toDefinition(rv)
			if t == nil || def == nil {
				return t != def
			}
			return !t.Compare(toDefinition(rv))
		default:
			return lv != castAsLeft(lv, rv)
		}
	}
}

func operOr(l getter, r getter) getter {
	return func(in []Getter) interface{} {
		return toBool(l(in)) || toBool(r(in))
	}
}

func operAnd(l getter, r getter) getter {
	return func(in []Getter) interface{} {
		return toBool(l(in)) && toBool(r(in))
	}
}

func operGreater(l getter, r getter) getter {
	return func(in []Getter) interface{} {
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
		case astdata.Definition:
			return false // definition dose not support this operators
		case null:
			return lv.(null) > rv.(null)
		}
		panic("not supported type")
	}
}

func operGreaterEqual(l getter, r getter) getter {
	return func(in []Getter) interface{} {
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
		case astdata.Definition:
			return false // definition dose not support this operators
		case null:
			return lv.(null) >= rv.(null)
		}
		panic("not supported type")
	}
}

func operLesser(l getter, r getter) getter {
	return func(in []Getter) interface{} {
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
		case astdata.Definition:
			return false // definition dose not support this operators
		case null:
			return lv.(null) < rv.(null)
		}
		panic("not supported type")
	}
}

func operLesserEqual(l getter, r getter) getter {
	return func(in []Getter) interface{} {
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
		case astdata.Definition:
			return false // definition dose not support this operators
		case null:
			return lv.(null) <= rv.(null)
		}
		panic("not supported type")
	}
}

func operIs(l getter, r getter) getter {
	// TODO : there is a problem, if he right operand is a column with null value it works here :)
	return func(in []Getter) interface{} {
		v := r(in)
		n, ok := v.(null)
		if !ok {
			panic("is only works on null/not null")
		}
		return toNull(l(in)) == n
	}
}

func operLike(l getter, r getter) getter {
	return func(in []Getter) interface{} {
		re := likeStr(toString(r(in))).regexp()
		v := toString(l(in))
		return re.MatchString(v)
	}
}

func operNot(l getter) getter {
	return func(in []Getter) interface{} {
		d := l(in)
		switch t := d.(type) {
		case bool:
			return !t
		case null:
			if t == notNullValue {
				return nullValue
			}
			return notNullValue
		default:
			panic("not is not applicable for string or number or definition")
		}
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
	case astdata.Definition:
		return toDefinition(r)
	case null:
		return toNull(r)
	}
	return r
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
	case astdata.Definition:
		return t != nil
	case null:
		return t != nullValue
	case nil:
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
	case astdata.Definition:
		return 0
	case null:
		return float64(t)
	case nil:
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
	case astdata.Definition:
		return t.String()
	case nil:
		return ""
	}
	panic(fmt.Sprintf("result from type %T", in))
}

func toDefinition(in interface{}) astdata.Definition {
	switch t := in.(type) {
	case bool, float64, nil:
		return nil
	case string:
		def, err := astdata.NewDefinition(t)
		if err != nil {
			return nil
		}
		return def
	case astdata.Definition:
		return t
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
	case astdata.Definition:
		return notNullValue
	case nil:
		return nullValue
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

func buildFilter(w parse.Stack, params ...interface{}) (getter, error) {
	var (
		p = parse.NewStack(0)
	)
	t, err := w.Pop()
	if err != nil {
		return func([]Getter) interface{} {
			return true
		}, nil
	}
	for {
		if isOperator(t.Type()) {
			r, e := p.Pop()
			if e != nil {
				return nil, fmt.Errorf("invalid operand")
			}
			rg := getGetter(r)
			l, e := p.Pop()
			if e != nil {
				return nil, fmt.Errorf("invalid operand")
			}
			lg := getGetter(l)
			g := getOpGetter(t, lg, rg)
			w.Push(g)
		} else if t.Type() == parse.ItemNot {
			// pop the operand
			ts, e := p.Pop()
			if e != nil {
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
