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
					typ:   ItemSelect,
					pos:   0,
					value: "SELECT",
				},
				item{
					typ:   ItemWhiteSpace,
					pos:   6,
					value: " ",
				},
				item{
					typ:   ItemWildCard,
					pos:   7,
					value: "*",
				},
				item{
					typ:   ItemWhiteSpace,
					pos:   8,
					value: " ",
				},
				item{
					typ:   ItemFrom,
					pos:   9,
					value: "FROM",
				},
				item{
					typ:   ItemWhiteSpace,
					pos:   13,
					value: " ",
				},
				item{
					typ:   ItemAlpha,
					pos:   14,
					value: "test",
				},
				item{
					typ:   ItemEOF,
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
			it := l.nextItem()
			assert.Equal(t, exp.value, it.Value())
			assert.Equal(t, exp.pos, it.Pos())
			assert.Equal(t, exp.typ, it.Type())
		}
	}
}

type lexTypeTest struct {
	Sql   string
	Items []ItemType
}

var lexType = []lexTypeTest{
	lexTypeTest{
		Sql:   "SELECT STRING STRING FROM TABLE",
		Items: []ItemType{ItemSelect, ItemAlpha, ItemAlpha, ItemFrom, ItemAlpha},
	},
	lexTypeTest{
		Sql:   "< > <= >= <> =",
		Items: []ItemType{ItemLesser, ItemGreater, ItemLesserEqual, ItemGreaterEqual, ItemNotEqual, ItemEqual},
	},
	lexTypeTest{
		Sql:   ". , ; ( )",
		Items: []ItemType{ItemDot, ItemComma, ItemSemicolon, ItemParenOpen, ItemParenClose},
	},
	lexTypeTest{
		Sql:   "( ( ) error ", // space at the end is for test sake. in loop we need every thing with an space in it
		Items: []ItemType{ItemParenOpen, ItemParenOpen, ItemParenClose, ItemAlpha, ItemError},
	},
	lexTypeTest{
		Sql:   ")",
		Items: []ItemType{ItemError},
	},
	lexTypeTest{
		Sql:   "11 2.22 1233 23434343 33.99",
		Items: []ItemType{ItemNumber, ItemNumber, ItemNumber, ItemNumber, ItemNumber},
	},
	lexTypeTest{
		Sql:   "11.00.2",
		Items: []ItemType{ItemError},
	},
	lexTypeTest{
		Sql:   `"str" 'another' "''ssss''" '""sss"' 'ss\'' "\\" '\\' "\""`,
		Items: []ItemType{ItemLiteral2, ItemLiteral1, ItemLiteral2, ItemLiteral1, ItemLiteral1, ItemLiteral2, ItemLiteral1, ItemLiteral2},
	},
	lexTypeTest{
		Sql:   `"str`,
		Items: []ItemType{ItemError},
	},
	lexTypeTest{
		Sql:   `'str`,
		Items: []ItemType{ItemError},
	},
	lexTypeTest{
		Sql:   `"str\c"`,
		Items: []ItemType{ItemError},
	},
	lexTypeTest{
		Sql:   `'s\tr'`,
		Items: []ItemType{ItemError},
	},
	lexTypeTest{
		Sql:   `&`,
		Items: []ItemType{ItemError},
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
				assert.Equal(t, ItemWhiteSpace, ll[j*2+1].typ)
			}
		}
	}
}

func TestMisc(t *testing.T) {
	w := item{
		typ: ItemAlpha,
	}
	assert.Panics(t, func() { assertType(w, ItemWhiteSpace) })
	assert.NotPanics(t, func() { assertType(w, ItemAlpha) })
}
