package goql

import (
	"fmt"
	"testing"

	"github.com/fzerorubigd/goql/parse"
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

type provider struct {
}

func (provider) Provide(in interface{}) []interface{} {
	ln := int(in.(tablet))
	res := make([]interface{}, ln)
	for i := 0; i < ln; i++ {
		res[i] = row(i)
	}

	return res
}

type concat struct {
}

func (concat) Execute(in ...Getter) (Getter, error) {
	s := ""
	for i := range in {
		s += fmt.Sprint(in[i].Get())
	}
	return String{String: s}, nil
}

type wrong struct {
}

func (wrong) Execute(in ...Getter) (Getter, error) {
	return nil, fmt.Errorf("hi, i am error")
}

func ast(q string) *parse.Query {
	ast, err := parse.AST(q)
	if err != nil {
		panic(err)
	}
	return ast
}

func TestContext(t *testing.T) {
	RegisterFunction("concat", concat{})
	RegisterTable("test", provider{})

	RegisterField("test", "c1", c1{})
	RegisterField("test", "c2", c2{})
	RegisterField("test", "c3", c3{})

	q := "SELECT * FROM test"
	row, data, err := execute(tablet(1), ast(q))
	assert.NoError(t, err)
	assert.Equal(t, 3, len(row))
	assert.Equal(t, []string{"c1", "c2", "c3"}, row)
	assert.Equal(t, 1, len(data))

	row, data, err = execute(tablet(10), ast(q))
	assert.NoError(t, err)
	assert.Equal(t, 3, len(row))
	assert.Equal(t, []string{"c1", "c2", "c3"}, row)
	assert.Equal(t, 10, len(data))

	q = "SELECT c1, c2 FROM test WHERE c3"
	row, data, err = execute(tablet(10), ast(q))
	assert.NoError(t, err)
	assert.Equal(t, 2, len(row))
	assert.Equal(t, []string{"c1", "c2"}, row)
	assert.Equal(t, 5, len(data))

	q = "SELECT c1, c2, c2 FROM test LIMIT 10"
	row, data, err = execute(tablet(100), ast(q))
	assert.NoError(t, err)
	assert.Equal(t, 3, len(row))
	assert.Equal(t, []string{"c1", "c2", "c2"}, row)
	assert.Equal(t, 10, len(data))

	q = `SELECT c1, c2, false FROM test WHERE "c2" like '%t_h%' OR "c3" > 0 or true and concat('t', 's') = 'ts' and c3 is not null LIMIT 10`
	row, data, err = execute(tablet(100), ast(q))
	assert.NoError(t, err)
	assert.Equal(t, 3, len(row))
	assert.Equal(t, []string{"c1", "c2", "static"}, row)
	assert.Equal(t, 10, len(data))

	q = "SELECT c1 FROM test WHERE c3 ORDER by c2 desc"
	row, data, err = execute(tablet(10), ast(q))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(row))
	assert.Equal(t, []string{"c1"}, row)
	assert.Equal(t, 5, len(data))

	q = "SELECT * FROM test limit 5,50"
	row, data, err = execute(tablet(10), ast(q))
	assert.NoError(t, err)
	assert.Equal(t, 3, len(row))
	assert.Equal(t, []string{"c1", "c2", "c3"}, row)
	assert.Equal(t, 5, len(data))

	q = "SELECT * FROM test limit 15,50"
	row, data, err = execute(tablet(10), ast(q))
	assert.NoError(t, err)
	assert.Equal(t, 3, len(row))
	assert.Equal(t, []string{"c1", "c2", "c3"}, row)
	assert.Equal(t, 0, len(data))

	q = "SELECT concat(c2, 'string'), 10 FROM test where true limit 1"
	row, data, err = execute(tablet(10), ast(q))
	assert.NoError(t, err)
	assert.Equal(t, 2, len(row))
	assert.Equal(t, []string{"concat", "static"}, row)
	assert.Equal(t, 1, len(data))
	assert.Equal(t, "0th rowstring", data[0][0].Get())
	assert.Equal(t, 10.0, data[0][1].Get())
}

func TestContextErr(t *testing.T) {
	RegisterFunction("wrong", wrong{})
	// Err
	q := "SELECT c1, c2 FROM notexists"
	row, data, err := execute(tablet(100), ast(q))
	assert.Error(t, err)
	assert.Nil(t, row)
	assert.Nil(t, data)

	// Err
	q = "SELECT c1, no FROM test"
	row, data, err = execute(tablet(100), ast(q))
	assert.Error(t, err)
	assert.Nil(t, row)
	assert.Nil(t, data)

	// Err
	q = `SELECT c1 FROM test WHERE "no"`
	row, data, err = execute(tablet(100), ast(q))
	assert.Error(t, err)
	assert.Nil(t, row)
	assert.Nil(t, data)

	// Err
	q = "SELECT c1, noooo.no FROM test"
	row, data, err = execute(tablet(100), ast(q))
	assert.Error(t, err)
	assert.Nil(t, row)
	assert.Nil(t, data)

	// Err
	q = "SELECT c1, c2 FROM test WHERE noo is null"
	row, data, err = execute(tablet(100), ast(q))
	assert.Error(t, err)
	assert.Nil(t, row)
	assert.Nil(t, data)

	// Err
	q = "SELECT c1, c2 FROM test ORDER BY no"
	row, data, err = execute(tablet(100), ast(q))
	assert.Error(t, err)
	assert.Nil(t, row)
	assert.Nil(t, data)

	q = "SELECT c1, fuuunccc(c2) FROM test"
	row, data, err = execute(tablet(100), ast(q))
	assert.Error(t, err)
	assert.Nil(t, row)
	assert.Nil(t, data)

	q = "SELECT c1, concat(c22) FROM test"
	row, data, err = execute(tablet(100), ast(q))
	assert.Error(t, err)
	assert.Nil(t, row)
	assert.Nil(t, data)

	// Not supporting func(*)
	q = "SELECT c1, concat(*) FROM test"
	row, data, err = execute(tablet(100), ast(q))
	assert.Error(t, err)
	assert.Nil(t, row)
	assert.Nil(t, data)

	q = "SELECT c1, wrong(c2) FROM test"
	row, data, err = execute(tablet(100), ast(q))
	assert.Error(t, err)
	assert.Nil(t, row)
	assert.Nil(t, data)

	q = "SELECT c1, c2 FROM test WHERE wrong()"
	row, data, err = execute(tablet(100), ast(q))
	assert.Error(t, err)
	assert.Nil(t, row)
	assert.Nil(t, data)

	g := func([]Getter) interface{} {
		panic("err")
	}

	b, err := callWhere(g, nil)
	assert.False(t, b)
	assert.Error(t, err)
}
