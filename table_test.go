package goql

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nilProvider struct {
}

func (nilProvider) Provide(in interface{}) []interface{} {
	return nil
}

func TestTables(t *testing.T) {
	RegisterTable("test1", provider{})

	RegisterField("test1", "c1", c1{})
	RegisterField("test1", "c2", c2{})
	RegisterField("test1", "c3", c3{})

	tt, err := getTable("not-exists")
	assert.Error(t, err)
	assert.Nil(t, tt)

	tbl, err := getTable("test")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(tbl))
	assert.Equal(t, 0, tbl["c1"].Order())
	assert.Equal(t, 1, tbl["c2"].Order())
	assert.Equal(t, 2, tbl["c3"].Order())

	assert.Equal(t, ValueTypeNumber, tbl["c1"].Type())
	assert.Equal(t, ValueTypeString, tbl["c2"].Type())
	assert.Equal(t, ValueTypeBool, tbl["c3"].Type())

	res := make(chan []Valuer, 3)

	err = getTableFields(tablet(1), "test1", res, "c1", "c2", "c3")
	assert.NoError(t, err)

	var cnt int64
	for i := range res {
		assert.Equal(t, 3, len(i))
		assert.Equal(t, float64(cnt*2), i[0].(Number).Number)
		assert.Equal(t, fmt.Sprintf("%dth row", cnt), i[1].(String).String)
		assert.Equal(t, cnt%2 == 0, i[2].(Bool).Bool)
		cnt++
	}

	res = make(chan []Valuer, 3)

	err = getTableFields(tablet(1), "test1", res, "c2", "c3")
	assert.NoError(t, err)

	cnt = 0
	for i := range res {
		assert.Equal(t, 2, len(i))
		assert.Equal(t, fmt.Sprintf("%dth row", cnt), i[0].(String).String)
		assert.Equal(t, cnt%2 == 0, i[1].(Bool).Bool)
		cnt++
	}

	res = make(chan []Valuer, 3)
	err = getTableFields(tablet(1), "test1", res, "c2", "", "c3")
	assert.NoError(t, err)

	cnt = 0
	for i := range res {
		assert.Equal(t, 3, len(i))
		assert.Equal(t, fmt.Sprintf("%dth row", cnt), i[0].(String).String)
		assert.Equal(t, cnt%2 == 0, i[2].(Bool).Bool)
		assert.Equal(t, nil, i[1])
		cnt++
	}

	assert.Panics(t, func() { RegisterTable("test1", nilProvider{}) })
	assert.Panics(t, func() { RegisterField("not-exist", "test", c1{}) })
	assert.Panics(t, func() { RegisterField("test1", "c1", c1{}) })
	assert.Panics(t, func() { RegisterField("test1", "c11", 10) })

	assert.Error(t, getTableFields(1, "not-exist", res, "col"))
	assert.Error(t, getTableFields(1, "test1", res, "col"))
	assert.Error(t, getTableFields(1, "test1", res))
}

func TestTypes(t *testing.T) {
	b := Bool{}
	assert.Equal(t, false, b.Value())
	b.Bool = true
	assert.Equal(t, true, b.Value())
	b.Null = true
	assert.Nil(t, b.Value())

	n := Number{}
	assert.Equal(t, 0.0, n.Value())
	n.Number = 10.0
	assert.Equal(t, 10.0, n.Value())
	n.Null = true
	assert.Nil(t, n.Value())

	s := String{}
	assert.Equal(t, "", s.Value())
	s.String = "test"
	assert.Equal(t, "test", s.Value())
	s.Null = true
	assert.Nil(t, s.Value())

	cd := columnDef{
		typ: 10,
	}
	assert.Panics(t, func() { cd.Type() })
}
