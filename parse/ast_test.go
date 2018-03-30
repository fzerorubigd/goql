package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStringToken(t *testing.T) {
	i := item{
		typ:   ItemLiteral1,
		value: `'this " is \'string\''`,
	}
	assert.Equal(t, `this " is 'string'`, GetTokenString(i))

	i = item{
		typ:   ItemLiteral2,
		value: `"this ' is \"string\""`,
	}
	assert.Equal(t, `this ' is "string"`, GetTokenString(i))

	i = item{
		typ:   ItemLiteral2,
		value: `this ' is \"string\"'`,
	}
	assert.Panics(t, func() { GetTokenString(i) })

	j := itemFn{
		item:      i,
		parameter: Fields{Field{}},
	}

	assert.Equal(t, Fields{Field{}}, j.Parameters())

}
