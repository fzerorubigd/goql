package executor

import (
	"fmt"
	"testing"

	"github.com/fzerorubigd/goql/structures"
	"github.com/stretchr/testify/assert"
)

type tablet int

type row int

type c1 struct {
}

func (c c1) Value(in interface{}) structures.Number {
	r := in.(row)
	return structures.Number{Number: float64(r) * 2.0}
}

type c2 struct {
}

func (c c2) Value(in interface{}) structures.String {
	r := in.(row)
	return structures.String{String: fmt.Sprintf("%dth row", r)}
}

type c3 struct {
}

func (c c3) Value(in interface{}) structures.Bool {
	r := in.(row)
	return structures.Bool{Bool: r%2 == 0}
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

func TestContext(t *testing.T) {
	structures.RegisterTable("test", provider{})

	structures.RegisterField("test", "c1", c1{})
	structures.RegisterField("test", "c2", c2{})
	structures.RegisterField("test", "c3", c3{})

	q := "SELECT * FROM test"
	row, data, err := Execute(tablet(1), q)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(row))
	assert.Equal(t, []string{"c1", "c2", "c3"}, row)
	assert.Equal(t, 1, len(data))

	row, data, err = Execute(tablet(10), q)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(row))
	assert.Equal(t, []string{"c1", "c2", "c3"}, row)
	assert.Equal(t, 10, len(data))

	q = "SELECT c1, c2 FROM test WHERE c3"
	row, data, err = Execute(tablet(10), q)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(row))
	assert.Equal(t, []string{"c1", "c2"}, row)
	assert.Equal(t, 5, len(data))

	q = "SELECT c1, c2, c2 FROM test LIMIT 10"
	row, data, err = Execute(tablet(100), q)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(row))
	assert.Equal(t, []string{"c1", "c2", "c2"}, row)
	assert.Equal(t, 10, len(data))

	q = `SELECT c1, c2 FROM test WHERE "c2" like '%t_h%' OR "c3" > 0  LIMIT 10`
	row, data, err = Execute(tablet(100), q)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(row))
	assert.Equal(t, []string{"c1", "c2"}, row)
	assert.Equal(t, 10, len(data))

	q = "SELECT c1 FROM test WHERE c3 ORDER by c2 desc"
	row, data, err = Execute(tablet(10), q)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(row))
	assert.Equal(t, []string{"c1"}, row)
	assert.Equal(t, 5, len(data))

	q = "SELECT * FROM test limit 5,50"
	row, data, err = Execute(tablet(10), q)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(row))
	assert.Equal(t, []string{"c1", "c2", "c3"}, row)
	assert.Equal(t, 5, len(data))

	q = "SELECT * FROM test limit 15,50"
	row, data, err = Execute(tablet(10), q)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(row))
	assert.Equal(t, []string{"c1", "c2", "c3"}, row)
	assert.Equal(t, 0, len(data))

}

func TestContextErr(t *testing.T) {
	// Err
	q := "SELECT c1, c2 FROM test  10"
	row, data, err := Execute(tablet(100), q)
	assert.Error(t, err)
	assert.Nil(t, row)
	assert.Nil(t, data)

	// Err
	q = "SELECT c1, c2 FROM notexists"
	row, data, err = Execute(tablet(100), q)
	assert.Error(t, err)
	assert.Nil(t, row)
	assert.Nil(t, data)

	// Err
	q = "SELECT c1, no FROM test"
	row, data, err = Execute(tablet(100), q)
	assert.Error(t, err)
	assert.Nil(t, row)
	assert.Nil(t, data)

	// Err
	q = `SELECT c1 FROM test WHERE "no"`
	row, data, err = Execute(tablet(100), q)
	assert.Error(t, err)
	assert.Nil(t, row)
	assert.Nil(t, data)

	// Err
	q = "SELECT c1, noooo.no FROM test"
	row, data, err = Execute(tablet(100), q)
	assert.Error(t, err)
	assert.Nil(t, row)
	assert.Nil(t, data)

	// Err
	q = "SELECT c1, c2 FROM test WHERE noo is null"
	row, data, err = Execute(tablet(100), q)
	assert.Error(t, err)
	assert.Nil(t, row)
	assert.Nil(t, data)

	// Err
	q = "SELECT c1, c2 FROM test ORDER BY no"
	row, data, err = Execute(tablet(100), q)
	assert.Error(t, err)
	assert.Nil(t, row)
	assert.Nil(t, data)

	g := func([]structures.Valuer) interface{} {
		panic("err")
	}

	b, err := callWhere(g, nil)
	assert.False(t, b)
	assert.Error(t, err)
}
