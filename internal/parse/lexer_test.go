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

type lexTypeTest struct {
	Sql   string
	Items []itemType
}

var lexType = []lexTypeTest{
	lexTypeTest{
		Sql:   "SELECT STRING STRING FROM TABLE",
		Items: []itemType{itemSelect, itemAlpha, itemAlpha, itemFrom, itemAlpha},
	},
	lexTypeTest{
		Sql:   "< > <= >= <> =",
		Items: []itemType{itemLesser, itemGreater, itemLesserEqual, itemGreaterEqual, itemNotEqual, itemEqual},
	},
	lexTypeTest{
		Sql:   ". , ; ( )",
		Items: []itemType{itemDot, itemComma, itemSemicolon, itemParenOpen, itemParenClose},
	},
	lexTypeTest{
		Sql:   "( ( ) error ", // space at the end is for test sake. in loop we need every thing with an space in it
		Items: []itemType{itemParenOpen, itemParenOpen, itemParenClose, itemAlpha, itemError},
	},
	lexTypeTest{
		Sql:   ")",
		Items: []itemType{itemError},
	},
	lexTypeTest{
		Sql:   "11 2.22 1233 23434343 33.99",
		Items: []itemType{itemNumber, itemNumber, itemNumber, itemNumber, itemNumber},
	},
	lexTypeTest{
		Sql:   "11.00.2",
		Items: []itemType{itemError},
	},
	lexTypeTest{
		Sql:   `"str" 'another' "''ssss''" '""sss"' 'ss\'' "\\" '\\' "\""`,
		Items: []itemType{itemLiteral2, itemLiteral1, itemLiteral2, itemLiteral1, itemLiteral1, itemLiteral2, itemLiteral1, itemLiteral2},
	},
	lexTypeTest{
		Sql:   `"str`,
		Items: []itemType{itemError},
	},
	lexTypeTest{
		Sql:   `'str`,
		Items: []itemType{itemError},
	},
	lexTypeTest{
		Sql:   `"str\c"`,
		Items: []itemType{itemError},
	},
	lexTypeTest{
		Sql:   `'s\tr'`,
		Items: []itemType{itemError},
	},
	lexTypeTest{
		Sql:   `&`,
		Items: []itemType{itemError},
	},
}

func TestAlpha(t *testing.T) {
	for i := range lexType {
		l := lex(lexType[i].Sql)
		var ll []item
		for j := range l.items {
			ll = append(ll, j)
		}
		assert.Equal(t, len(lexType[i].Items)*2-1, len(ll))
		for j := range lexType[i].Items {
			assert.Equal(t, lexType[i].Items[j], ll[j*2].typ)
			if j*2+1 < len(ll) {
				assert.Equal(t, itemWhiteSpace, ll[j*2+1].typ)
			}
		}
	}
}
