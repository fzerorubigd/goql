package goql

import (
	"database/sql"
	drv "database/sql/driver"
	"fmt"
	"io"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/fzerorubigd/goql/parse"
)

type driver struct{}

type conn struct {
	pkg *astdata.Package
}

type stmt struct {
	pkg   *astdata.Package
	query *parse.Query
}

type rows struct {
	cursor int
	rows   []string
	data   [][]Getter
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
	r := &rows{}
	r.rows, r.data, err = execute(gs.pkg, gs.query)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *rows) Columns() []string {
	return r.rows
}

func (r *rows) Close() error {
	r.data = nil
	r.rows = nil
	return nil
}

func (r *rows) Next(dest []drv.Value) error {
	if r.cursor >= len(r.data) {
		return io.EOF
	}

	for i := range dest {
		in := r.data[r.cursor][i].Get()
		dest[i] = in
	}
	r.cursor++
	return nil
}

func init() {
	sql.Register("goql", &driver{})
}
