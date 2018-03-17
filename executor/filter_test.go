package executor

import (
	"testing"

	"github.com/fzerorubigd/goql/internal/parse"
	"github.com/fzerorubigd/goql/structures"
	"github.com/stretchr/testify/assert"
)

type dummy struct {
	typ   parse.ItemType
	pos   int
	value string
}

func (d dummy) Type() parse.ItemType {
	return d.typ
}

func (d dummy) Pos() int {
	return d.pos
}

func (d dummy) Value() string {
	return d.value
}

func (d dummy) String() string {
	return "dummy"
}

var (
	rowTest = []structures.Valuer{
		structures.Number{Number: 1},
		structures.Bool{Bool: true},
		structures.String{String: "a"},
		structures.Number{Null: true},
		structures.Bool{Null: true},
		structures.String{Null: true},
		unknown{},
	}
)

func newItem(t parse.ItemType, v string, p int) parse.Item {
	return dummy{
		typ:   t,
		pos:   p,
		value: v,
	}
}

func newGetter(in interface{}) getter {
	return func([]structures.Valuer) interface{} {
		return in
	}
}

func TestWhereOp(t *testing.T) {
	//null
	assert.Panics(t, func() { nullGetterGenerator(newItem(parse.ItemEOF, "", 0)) })
	assert.Equal(t, nullValue, nullGetterGenerator(newItem(parse.ItemNull, "", 0))(nil))

	// alpha
	assert.Panics(t, func() { alphaGetterGenerator(newItem(parse.ItemEOF, "", 0)) })
	assert.Equal(t, true, alphaGetterGenerator(newItem(parse.ItemAlpha, "true", 0))(nil))
	assert.Equal(t, false, alphaGetterGenerator(newItem(parse.ItemAlpha, "false", 0))(nil))
	assert.Panics(t, func() { alphaGetterGenerator(newItem(parse.ItemAlpha, "anything", 0))(nil) })

	// field
	assert.Panics(t, func() { fieldGetterGenerator(newItem(parse.ItemEOF, "", 0)) })
	assert.Equal(t, 1.0, fieldGetterGenerator(newItem(itemColumn, "", 0))(rowTest))
	assert.Equal(t, true, fieldGetterGenerator(newItem(itemColumn, "", 1))(rowTest))
	assert.Equal(t, "a", fieldGetterGenerator(newItem(itemColumn, "", 2))(rowTest))
	assert.Equal(t, nullValue, fieldGetterGenerator(newItem(itemColumn, "", 3))(rowTest))
	assert.Equal(t, nullValue, fieldGetterGenerator(newItem(itemColumn, "", 4))(rowTest))
	assert.Equal(t, nullValue, fieldGetterGenerator(newItem(itemColumn, "", 5))(rowTest))
	assert.Equal(t, unknown{}, fieldGetterGenerator(newItem(itemColumn, "", 6))(rowTest))

	// litteral
	assert.Panics(t, func() { literal1GetterGenerator(newItem(parse.ItemEOF, "", 0)) })
	assert.Equal(t, "test", literal1GetterGenerator(newItem(parse.ItemLiteral1, "'test'", 0))(nil))

	// Number
	assert.Panics(t, func() { numberGetterGenerator(newItem(parse.ItemEOF, "", 0)) })
	assert.Equal(t, 42.42, numberGetterGenerator(newItem(parse.ItemNumber, "42.42", 0))(nil))

	// equal
	assert.Equal(t, true, equal(newGetter(10.0), newGetter(10.0))(nil))
	assert.Equal(t, false, equal(newGetter(11.0), newGetter(10.0))(nil))
	assert.Equal(t, true, equal(newGetter(true), newGetter(true))(nil))
	assert.Equal(t, false, equal(newGetter(true), newGetter(false))(nil))
	assert.Equal(t, true, equal(newGetter("aaa"), newGetter("aaa"))(nil))
	assert.Equal(t, false, equal(newGetter("bbb"), newGetter("aaa"))(nil))

	assert.Equal(t, false, notEqual(newGetter(10.0), newGetter(10.0))(nil))
	assert.Equal(t, true, notEqual(newGetter(11.0), newGetter(10.0))(nil))
	assert.Equal(t, false, notEqual(newGetter(true), newGetter(true))(nil))
	assert.Equal(t, true, notEqual(newGetter(true), newGetter(false))(nil))
	assert.Equal(t, false, notEqual(newGetter("aaa"), newGetter("aaa"))(nil))
	assert.Equal(t, true, notEqual(newGetter("bbb"), newGetter("aaa"))(nil))

	// or / and
	assert.Equal(t, true, operOr(newGetter(true), newGetter(false))(nil))
	assert.Equal(t, false, operAnd(newGetter(true), newGetter(false))(nil))
	assert.Equal(t, true, operOr(newGetter(true), newGetter(true))(nil))
	assert.Equal(t, true, operAnd(newGetter(true), newGetter(true))(nil))

	// < > <= >=
	assert.Equal(t, true, operGreater(newGetter(true), newGetter(false))(nil))
	assert.Equal(t, false, operGreater(newGetter(true), newGetter(true))(nil))
	assert.Equal(t, true, operGreater(newGetter("reza"), newGetter("ali"))(nil))
	assert.Equal(t, true, operGreater(newGetter(15.0), newGetter(10.0))(nil))
	assert.Equal(t, true, operGreater(newGetter(notNullValue), newGetter(nullValue))(nil))
	assert.Panics(t, func() { operGreater(newGetter(struct{}{}), newGetter(struct{}{}))(nil) })

	assert.Equal(t, true, operGreaterEqual(newGetter(true), newGetter(false))(nil))
	assert.Equal(t, true, operGreaterEqual(newGetter(true), newGetter(true))(nil))
	assert.Equal(t, true, operGreaterEqual(newGetter("reza"), newGetter("ali"))(nil))
	assert.Equal(t, true, operGreaterEqual(newGetter(15.0), newGetter(10.0))(nil))
	assert.Equal(t, true, operGreaterEqual(newGetter(notNullValue), newGetter(nullValue))(nil))
	assert.Panics(t, func() { operGreaterEqual(newGetter(struct{}{}), newGetter(struct{}{}))(nil) })

	assert.Equal(t, false, operLesser(newGetter(true), newGetter(false))(nil))
	assert.Equal(t, false, operLesser(newGetter(true), newGetter(true))(nil))
	assert.Equal(t, false, operLesser(newGetter("reza"), newGetter("ali"))(nil))
	assert.Equal(t, false, operLesser(newGetter(15.0), newGetter(10.0))(nil))
	assert.Equal(t, false, operLesser(newGetter(notNullValue), newGetter(nullValue))(nil))
	assert.Panics(t, func() { operLesser(newGetter(struct{}{}), newGetter(struct{}{}))(nil) })

	assert.Equal(t, false, operLesserEqual(newGetter(true), newGetter(false))(nil))
	assert.Equal(t, true, operLesserEqual(newGetter(true), newGetter(true))(nil))
	assert.Equal(t, false, operLesserEqual(newGetter("reza"), newGetter("ali"))(nil))
	assert.Equal(t, false, operLesserEqual(newGetter(15.0), newGetter(10.0))(nil))
	assert.Equal(t, false, operLesserEqual(newGetter(notNullValue), newGetter(nullValue))(nil))
	assert.Panics(t, func() { operLesserEqual(newGetter(struct{}{}), newGetter(struct{}{}))(nil) })

	// is
	assert.Equal(t, true, operIs(newGetter(nullValue), newGetter(nullValue))(nil))
	assert.Equal(t, false, operIs(newGetter(nullValue), newGetter(notNullValue))(nil))
	assert.Panics(t, func() { operIs(newGetter(nullValue), newGetter(99.99))(nil) })

	// not
	assert.Equal(t, false, operNot(newGetter(true))(nil))
	assert.Equal(t, notNullValue, operNot(newGetter(nullValue))(nil))
	assert.Equal(t, nullValue, operNot(newGetter(notNullValue))(nil))
	assert.Panics(t, func() { operNot(newGetter("string"))(nil) })

	// casts
	assert.IsType(t, true, toBool(nullValue))
	assert.False(t, toBool(nullValue))
	assert.True(t, toBool(notNullValue))
	assert.True(t, toBool("true"))
	assert.Panics(t, func() { toBool(struct{}{}) })

	assert.Equal(t, 0.0, toNumber(nullValue))
	assert.Equal(t, 1.0, toNumber(notNullValue))
	assert.Equal(t, 1.0, toNumber(true))
	assert.Equal(t, 0.0, toNumber(false))
	assert.Equal(t, 10.0, toNumber(10.0))
	assert.Equal(t, 10.0, toNumber("10.0"))
	assert.Panics(t, func() { toNumber(struct{}{}) })

	assert.Equal(t, "1", toString(1.0))
	assert.Equal(t, "true", toString(true))
	assert.Equal(t, "false", toString(false))
	assert.Equal(t, "hi", toString("hi"))
	assert.Panics(t, func() { toString(struct{}{}) })

	assert.Equal(t, notNullValue, toNull(1.0))
	assert.Equal(t, notNullValue, toNull("true"))
	assert.Equal(t, notNullValue, toNull(false))
	assert.Equal(t, nullValue, toNull(nullValue))
	assert.Panics(t, func() { toNull(struct{}{}) })

	// getter
	assert.Panics(t, func() { getGetter(newItem(parse.ItemEOF, "", 0)) })
	assert.Panics(t, func() { getOpGetter(newItem(parse.ItemEOF, "", 0), newGetter(1.0), newGetter(1.0)) })

}
