package structures

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type tablet int

type row int

type c1 struct {
}

func (c c1) Value(in interface{}) Number {
	r := in.(row)
	return Number{Number: float64(r) * 2.0}
}

type c2 struct {
}

func (c c2) Value(in interface{}) String {
	r := in.(row)
	return String{String: fmt.Sprintf("%dth row", r)}
}

type c3 struct {
}

func (c c3) Value(in interface{}) Bool {
	r := in.(row)
	return Bool{Bool: r%2 == 0}
}

func TestTables(t *testing.T) {
	RegisterTable("test", func(in interface{}) []interface{} {
		tbl := in.(tablet)
		ln := int(tbl) * 10
		res := make([]interface{}, ln)
		for i := 0; i < ln; i++ {
			res[i] = row(i)
		}

		return res
	})

	RegisterField("test", "c1", c1{})
	RegisterField("test", "c2", c2{})
	RegisterField("test", "c3", c3{})

	res := make(chan []interface{}, 3)

	err := GetFields(tablet(1), "test", res, "c1", "c2", "c3")
	assert.NoError(t, err)

	var cnt int64
	for i := range res {
		assert.Equal(t, 3, len(i))
		assert.Equal(t, float64(cnt*2), i[0].(Number).Number)
		assert.Equal(t, fmt.Sprintf("%dth row", cnt), i[1].(String).String)
		assert.Equal(t, cnt%2 == 0, i[2].(Bool).Bool)
		cnt++
	}

	res = make(chan []interface{}, 3)

	err = GetFields(tablet(1), "test", res, "c2", "c3")
	assert.NoError(t, err)

	cnt = 0
	for i := range res {
		assert.Equal(t, 2, len(i))
		assert.Equal(t, fmt.Sprintf("%dth row", cnt), i[0].(String).String)
		assert.Equal(t, cnt%2 == 0, i[1].(Bool).Bool)
		cnt++
	}

	assert.Panics(t, func() { RegisterTable("test", func(in interface{}) []interface{} { return nil }) })
	assert.Panics(t, func() { RegisterField("not-exist", "test", c1{}) })
	assert.Panics(t, func() { RegisterField("test", "c1", c1{}) })
	assert.Panics(t, func() { RegisterField("test", "c11", 10) })

	assert.Error(t, GetFields(1, "not-exist", res, "col"))
	assert.Error(t, GetFields(1, "test", res, "col"))
	assert.Error(t, GetFields(1, "test", res))
}
