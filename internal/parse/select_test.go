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
	assert.Equal(t, Field{Column: "a"}, ss.Fields[0])
	assert.Equal(t, Field{Column: "b"}, ss.Fields[1])
	assert.Equal(t, Field{Column: "c", Table: "test"}, ss.Fields[2])

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
	assert.Equal(t, Field{WildCard: true}, ss.Fields[0])

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

	q = "SELECT * FROM x where  x (  )"
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
