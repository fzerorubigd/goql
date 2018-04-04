package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectSimple(t *testing.T) {
	q := "SELECT a,b,test.c FROM test"
	stmt, err := AST(q)
	assert.NoError(t, err)
	assert.IsType(t, &SelectStmt{}, stmt.Statement)
	ss := stmt.Statement.(*SelectStmt)
	assert.Equal(t, "test", ss.Table)

	assert.Equal(t, 3, len(ss.Fields))
	assert.Equal(t, "a", ss.Fields[0].Item.Value())
	assert.Equal(t, "b", ss.Fields[1].Item.Value())
	assert.Equal(t, "c", ss.Fields[2].Item.Value())
	assert.Equal(t, "test", ss.Fields[2].Table)

	q = "SELECT func(a,'b',10), c, 10, 'string' FROM test"
	stmt, err = AST(q)
	assert.NoError(t, err)
	assert.IsType(t, &SelectStmt{}, stmt.Statement)
	ss = stmt.Statement.(*SelectStmt)
	assert.Equal(t, "test", ss.Table)

	assert.Equal(t, 4, len(ss.Fields))
	fn := ss.Fields[0]
	assert.Equal(t, "func", fn.Item.Value())
	assert.Equal(t, 3, len(fn.Parameters))
	assert.Equal(t, "a", fn.Parameters[0].Item.Value())
	assert.Equal(t, ItemAlpha, fn.Parameters[0].Item.Type())
	assert.Equal(t, "'b'", fn.Parameters[1].Item.Value())
	assert.Equal(t, ItemLiteral1, fn.Parameters[1].Item.Type())
	assert.Equal(t, "10", fn.Parameters[2].Item.Value())
	assert.Equal(t, ItemNumber, fn.Parameters[2].Item.Type())

	assert.Equal(t, "c", ss.Fields[1].Item.Value())
	assert.Equal(t, ItemAlpha, ss.Fields[1].Item.Type())
	assert.Equal(t, "10", ss.Fields[2].Item.Value())
	assert.Equal(t, ItemNumber, ss.Fields[2].Item.Type())
	assert.Equal(t, "'string'", ss.Fields[3].Item.Value())
	assert.Equal(t, ItemLiteral1, ss.Fields[3].Item.Type())

	q = "SELECT FN1(FN2(FN3(), x)) FROM test"
	stmt, err = AST(q)
	assert.NoError(t, err)
	assert.IsType(t, &SelectStmt{}, stmt.Statement)
	ss = stmt.Statement.(*SelectStmt)
	assert.Equal(t, "test", ss.Table)

	assert.Equal(t, 1, len(ss.Fields))
	fn1 := ss.Fields[0]
	assert.Equal(t, "FN1", fn1.Item.Value())
	assert.Equal(t, 1, len(fn1.Parameters))
	assert.Equal(t, ItemFunc, fn1.Parameters[0].Item.Type())
	fn2 := fn1.Parameters[0]
	assert.Equal(t, "FN2", fn2.Item.Value())
	assert.Equal(t, 2, len(fn2.Parameters))
	assert.Equal(t, "x", fn2.Parameters[1].Item.Value())
	assert.Equal(t, ItemAlpha, fn2.Parameters[1].Item.Type())

	assert.Equal(t, ItemFunc, fn2.Parameters[0].Item.Type())
	fn3 := fn2.Parameters[0]
	assert.Equal(t, "FN3", fn3.Item.Value())
	assert.Equal(t, 0, len(fn3.Parameters))

	q = "SELECT * FROM x where  x (  )"
	_, err = AST(q)
	assert.NoError(t, err)

	q = "SELECT a,, FROM test"
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)

	q = "SELECT * FROM ,"
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)

	q = "SELECT * FROM test hahaha"
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)

	q = "SELECT func(invalid,) FROM test "
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)

	q = "SELECT func(invalid,  test  |)  "
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)

	q = "SELECT test. From test "
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)

}

func pop(t *testing.T, stack Stack) Item {
	p, err := stack.Pop()
	assert.NoError(t, err)
	return p
}

func TestSelectWhere(t *testing.T) {
	q := "SELECT * FROM test WHERE id = 2 "
	stmt, err := AST(q)
	assert.NoError(t, err)
	assert.IsType(t, &SelectStmt{}, stmt.Statement)
	ss := stmt.Statement.(*SelectStmt)
	assert.Equal(t, "test", ss.Table)

	assert.Equal(t, 1, len(ss.Fields))
	assert.Equal(t, ItemWildCard, ss.Fields[0].Item.Type())

	ip := pop(t, ss.Where)
	assert.Equal(t, ItemEqual, ip.Type())

	ip = pop(t, ss.Where)
	assert.Equal(t, ItemNumber, ip.Type())
	assert.Equal(t, "2", ip.Value())

	assert.NotNil(t, ss.Where)
	ip = pop(t, ss.Where)
	assert.Equal(t, ItemAlpha, ip.Type())
	assert.Equal(t, "id", ip.Value())

	_, err = ss.Where.Pop()
	assert.Error(t, err)

	q = "SELECT * FROM test WHERE a like '%ss%' and x or not s or (x is not null)"
	stmt, err = AST(q)
	assert.NoError(t, err)
	assert.IsType(t, &SelectStmt{}, stmt.Statement)
	ss = stmt.Statement.(*SelectStmt)
	assert.Equal(t, "test", ss.Table)
	/*
		TODO : Pop and check the values
	*/
	q = "SELECT * FROM x where "
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)

	q = "SELECT * FROM x where x x "
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)

	q = "SELECT * FROM x where and or  "
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)

	q = "SELECT * FROM x where not  "
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)

	q = "SELECT * FROM x where not not "
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)

	q = "SELECT * FROM x where ( x = )"
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)

	q = "SELECT * FROM x where ( and 2 )"
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)

	q = "SELECT * FROM x where (  )"
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)

	q = "SELECT * FROM x where fn () ()"
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)

	q = "SELECT * FROM x where fn (x,) "
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)

	q = "SELECT * FROM x where fn (x|) "
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)

	q = "SELECT * FROM x where  x ; x"
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)
}

func TestSelectOrder(t *testing.T) {
	q := "SELECT * FROM test WHERE id = 2 ORDER BY aa asc, bb desc"
	stmt, err := AST(q)
	assert.NoError(t, err)
	assert.IsType(t, &SelectStmt{}, stmt.Statement)
	ss := stmt.Statement.(*SelectStmt)
	assert.Equal(t, "test", ss.Table)
	assert.Equal(t, 2, len(ss.Order))

	assert.Equal(t, Order{Field: "aa"}, ss.Order[0])
	assert.Equal(t, Order{Field: "bb", DESC: true}, ss.Order[1])

	q = "SELECT * FROM x order "
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)

	q = "SELECT * FROM x order by , "
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)

}

func TestSelectLimit(t *testing.T) {
	q := "SELECT * FROM test limit 1, 10 "
	stmt, err := AST(q)
	assert.NoError(t, err)
	assert.IsType(t, &SelectStmt{}, stmt.Statement)
	ss := stmt.Statement.(*SelectStmt)
	assert.Equal(t, "test", ss.Table)

	assert.Equal(t, 1, ss.Start)
	assert.Equal(t, 10, ss.Count)

	q = "SELECT * FROM test limit 1 "
	stmt, err = AST(q)
	assert.NoError(t, err)
	assert.IsType(t, &SelectStmt{}, stmt.Statement)
	ss = stmt.Statement.(*SelectStmt)
	assert.Equal(t, "test", ss.Table)

	assert.Equal(t, 0, ss.Start)
	assert.Equal(t, 1, ss.Count)

	q = "SELECT * FROM test "
	stmt, err = AST(q)
	assert.NoError(t, err)
	assert.IsType(t, &SelectStmt{}, stmt.Statement)
	ss = stmt.Statement.(*SelectStmt)
	assert.Equal(t, "test", ss.Table)

	assert.Equal(t, -1, ss.Start)
	assert.Equal(t, -1, ss.Count)

	q = "SELECT * FROM x limit "
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)

	q = "SELECT * FROM x limit 1.88"
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)

	q = "SELECT * FROM x limit 1,"
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)

	q = "SELECT * FROM x limit 1,1.99"
	stmt, err = AST(q)
	assert.Error(t, err)
	assert.Nil(t, stmt)

}
