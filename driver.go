package goql

import (
	"database/sql"
	drv "database/sql/driver"
	"fmt"
	"io"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/fzerorubigd/goql/executor"
	"github.com/fzerorubigd/goql/internal/parse"
	"github.com/fzerorubigd/goql/structures"
)

type driver struct{}

type conn struct {
	pkg *astdata.Package
}

type stmt struct {
	pkg      *astdata.Package
	query    *parse.Query
	executed bool
}

type row struct {
	cursor int
	rows   []string
	data   [][]structures.Valuer
}

func (driver) Open(name string) (drv.Conn, error) {
	p, err := astdata.ParsePackage(name)
	if err != nil {
		return nil, err
	}

	return &conn{pkg: p}, nil
}

func (gc *conn) Prepare(query string) (drv.Stmt, error) {
	st, err := parse.AST(query)
	if err != nil {
		return nil, err
	}
	return &stmt{pkg: gc.pkg, query: st}, nil
}

func (gc *conn) Close() error {
	return nil
}

func (gc *conn) Begin() (drv.Tx, error) {
	return nil, fmt.Errorf("not supported")
}

func (gs *stmt) Close() error {
	return nil
}

func (gs *stmt) NumInput() int {
	return 0
}

func (gs *stmt) Exec(args []drv.Value) (drv.Result, error) {
	return nil, fmt.Errorf("currently only select via query is supported")
}

func (gs *stmt) Query(args []drv.Value) (drv.Rows, error) {
	var err error
	r := &row{}
	r.rows, r.data, err = executor.Execute(gs.pkg, gs.query)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *row) Columns() []string {
	return r.rows
}

func (r *row) Close() error {
	r.data = nil
	r.rows = nil
	return nil
}

func (r *row) Next(dest []drv.Value) error {
	if r.cursor >= len(r.data) {
		return io.EOF
	}

	for i := range dest {
		dest[i] = r.data[r.cursor][i].Value()
	}
	r.cursor++
	return nil
}

func init() {
	sql.Register("goql", &driver{})
}
