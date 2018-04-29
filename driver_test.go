package goql

import (
	"database/sql"
	"testing"
	"time"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDriver(t *testing.T) {
	o, err := sql.Open("goql", "github.com/not/exist/path")
	assert.NoError(t, err)
	assert.Error(t, o.Ping()) // actual connection :)

	c, err := sql.Open("goql", "github.com/fzerorubigd/fixture")
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, c.Close())
	}()
	assert.NoError(t, c.Ping())

	tx, err := c.Begin()
	assert.Error(t, err, "not supported")
	assert.Nil(t, tx)

	ss, err := c.Prepare("select * from")
	assert.Error(t, err)
	assert.Nil(t, ss)

	ss, err = c.Prepare("select * from notable")
	require.NoError(t, err)
	_, err = ss.Query()
	assert.Error(t, err)
	assert.NoError(t, ss.Close())

	ss, err = c.Prepare("select name, pkg_name,def from funcs")
	defer func() {
		assert.NoError(t, ss.Close())
	}()
	assert.NoError(t, err)
	q, err := ss.Exec()
	assert.Error(t, err)
	assert.Nil(t, q)

	s, err := ss.Query("args")
	assert.Error(t, err)
	assert.Nil(t, s)

	s, err = ss.Query()
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, s.Close())
	}()

	for s.Next() {
		var name, pkg string
		var def astdata.Definition

		assert.NoError(t, s.Scan(&name, &pkg, &def))
		assert.NotEmpty(t, name)
		assert.NotEmpty(t, pkg)
		assert.NotEmpty(t, def.String())
	}
}

func TestQueryParams(t *testing.T) {
	c, err := sql.Open("goql", "github.com/fzerorubigd/fixture")
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, c.Close())
	}()

	ss, err := c.Prepare("select name, def from funcs where def=?")
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, ss.Close())
	}()
	s, err := ss.Query("func(error)")
	require.NoError(t, err)

	defer func() {
		assert.NoError(t, s.Close())
	}()

	for s.Next() {
		var name string
		var def astdata.Definition

		assert.NoError(t, s.Scan(&name, &def))
		assert.NotEmpty(t, name)
		assert.IsType(t, &astdata.FuncType{}, def)
		assert.Equal(t, "func (error)", def.String())
	}

	// I don't know if having this kind of bind is normal
	ss, err = c.Prepare("select ?, ?, ?, ?, ?, ?, ?, ? from funcs")
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, ss.Close())
	}()
	s, err = ss.Query(true, "hi", []byte("bye"), int(10), int64(100), float32(1000), float64(10000), time.Now())
	require.NoError(t, err)

	defer func() {
		assert.NoError(t, s.Close())
	}()

	for s.Next() {
		var (
			b     bool
			hi    string
			bye   []byte
			in10  int
			in100 int64
			in32  float32
			in64  float64
			def   string
		)
		assert.NoError(t, s.Scan(&b, &hi, &bye, &in10, &in100, &in32, &in64, &def))
		assert.True(t, b)
		assert.Equal(t, "hi", hi)
		assert.Equal(t, []byte("bye"), bye)
		assert.Equal(t, int(10), in10)
		assert.Equal(t, int64(100), in100)
		assert.Equal(t, float32(1000), in32)
		assert.Equal(t, float64(10000), in64)
		assert.Equal(t, "time.Time", def)
	}

}
