package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParesr(t *testing.T) {
	q := "SELECT * FROM TEST WHERE INVALID STRING HERE"
	p := &parser{
		l: lex(q),
	}

	ahead := p.scan()
	assert.Equal(t, ahead.Type(), ItemSelect)
	assert.Equal(t, ahead.Value(), "SELECT")

	p.reject()
	ahead = p.scan()
	assert.Equal(t, ahead.Type(), ItemSelect)
	assert.Equal(t, ahead.Value(), "SELECT")

	ahead = p.scanIgnoreWhiteSpace()
	assert.Equal(t, ahead.Type(), ItemWildCard)
	assert.Equal(t, ahead.Value(), "*")

}

func TestAST(t *testing.T) {
	q := "SELECT * FROM TEST"
	s, err := AST(q)
	assert.NoError(t, err)
	assert.IsType(t, &SelectStmt{}, s.Statement)

	q = "SELECT * TEST FROM HI"
	_, err = AST(q)
	assert.Error(t, err)

	q = "UPDATE TEST SET x=1"
	_, err = AST(q)
	assert.Error(t, err)

}
