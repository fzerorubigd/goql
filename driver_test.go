package goql

import (
	"database/sql"
	"testing"

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
		var name, pkg, def string

		assert.NoError(t, s.Scan(&name, &pkg, &def))
		assert.NotEmpty(t, name)
		assert.NotEmpty(t, pkg)
		assert.NotEmpty(t, def)
	}
}
