package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type lexerTestHelper struct {
	Str   string
	Items []item
}

var (
	lexerTest = []lexerTestHelper{
		{
			Str: "SELECT * FROM test",
			Items: []item{
				item{
					typ:   itemSelect,
					pos:   0,
					value: "SELECT",
				},
				item{
					typ:   itemWhiteSpace,
					pos:   6,
					value: " ",
				},
				item{
					typ:   itemWildCard,
					pos:   7,
					value: "*",
				},
				item{
					typ:   itemWhiteSpace,
					pos:   8,
					value: " ",
				},
				item{
					typ:   itemFrom,
					pos:   9,
					value: "FROM",
				},
				item{
					typ:   itemWhiteSpace,
					pos:   13,
					value: " ",
				},
				item{
					typ:   itemAlpha,
					pos:   14,
					value: "test",
				},
				item{
					typ:   itemEOF,
					pos:   0,
					value: "",
				},
			},
		},
	}
)

func TestGeneralLexer(t *testing.T) {
	for i := range lexerTest {
		l := lex(lexerTest[i].Str)
		for _, exp := range lexerTest[i].Items {
			item := l.nextItem()
			assert.Equal(t, exp.value, item.value)
			assert.Equal(t, exp.pos, item.pos)
			assert.Equal(t, exp.typ, item.typ)
		}
	}
}
