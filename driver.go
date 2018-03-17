package goql

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/fzerorubigd/goql/executor"
	// runtime is the go type runtimes
	_ "github.com/fzerorubigd/goql/internal/runtime"
	"github.com/fzerorubigd/goql/structures"
)

type goqlDriver struct{}

type goqlConn struct {
	pkg *astdata.Package
}

type goqlStmt struct {
	pkg      *astdata.Package
	query    string
	executed bool
}

type row struct {
	cursor int
	rows   []string
	data   [][]structures.Valuer
}

func (goqlDriver) Open(name string) (driver.Conn, error) {
	p, err := astdata.ParsePackage(name)
	if err != nil {
		return nil, err
	}

	return &goqlConn{pkg: p}, nil
}

func (gc *goqlConn) Prepare(query string) (driver.Stmt, error) {
	// TODO : move parser out of internal and change execute to accept ast tree
	return &goqlStmt{pkg: gc.pkg, query: query}, nil
}

func (gc *goqlConn) Close() error {
	return nil
}

func (gc *goqlConn) Begin() (driver.Tx, error) {
	return nil, fmt.Errorf("not supported")
}

func (gs *goqlStmt) Close() error {
	return nil
}

func (gs *goqlStmt) NumInput() int {
	return 0
}

func (gs *goqlStmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, fmt.Errorf("currently only select via query is supported")
}

func (gs *goqlStmt) Query(args []driver.Value) (driver.Rows, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("args not supported yet, but %d args is provided", len(args))
	}
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

func (r *row) Next(dest []driver.Value) error {
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
	sql.Register("goql", &goqlDriver{})
}
