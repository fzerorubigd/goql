package parse

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type stateFn func(*lexer) stateFn

// ItemType is each lexeme type, also error and eof. any new keyword should be here first
type ItemType int

const (
	eof = 0
)

const (
	// ItemEOF eof is zero, so any data from closed channel (with zero value) is eof
	ItemEOF ItemType = iota
	// ItemError when there is an error
	ItemError
	// ItemWhiteSpace any whitespace sequence
	ItemWhiteSpace
	// ItemSelect sql select stmt
	ItemSelect
	// ItemFrom sql from stmt
	ItemFrom
	// ItemWhere sql where stmt
	ItemWhere
	// ItemOrder sql order stmt
	ItemOrder
	// ItemBy sql by stmt
	ItemBy
	// ItemOr sql or stmt
	ItemOr
	// ItemAnd sql and stmt
	ItemAnd
	// ItemIs sql is stmt
	ItemIs
	// ItemNull sql null stmt
	ItemNull
	// ItemNot sql not stmt
	ItemNot
	// ItemLimit sql limit stmt
	ItemLimit
	// ItemAsc sql asc stmt
	ItemAsc
	// ItemDesc sql desc stmt
	ItemDesc
	// ItemLike sql like stmt
	ItemLike
	// ItemAlpha is a string
	ItemAlpha
	// ItemNumber is a number
	ItemNumber
	// ItemFalse is the bool expression (false)
	ItemFalse
	// ItemTrue is the true
	ItemTrue
	// ItemEqual is =
	ItemEqual
	// ItemGreater is >
	ItemGreater
	// ItemLesser is <
	ItemLesser
	// ItemGreaterEqual is >=
	ItemGreaterEqual
	// ItemLesserEqual is <=
	ItemLesserEqual
	// ItemNotEqual is <>
	ItemNotEqual
	// ItemParenOpen is (
	ItemParenOpen
	// ItemParenClose is )
	ItemParenClose
	// ItemComma is ,
	ItemComma
	// ItemWildCard is *
	ItemWildCard
	// ItemLiteral1 is 'string in single quote'
	ItemLiteral1
	// ItemLiteral2 is "string in double quote"
	ItemLiteral2
	// ItemSemicolon is ;
	ItemSemicolon
	// ItemDot is .
	ItemDot
	// ItemDollarSign is the $ followed by the number
	ItemQuestionMark
	// ItemFunc is the function
	ItemFunc
)

var (
	// alphaItem which are an keyword, all lower case
	keywords = map[string]ItemType{
		"select": ItemSelect,
		"from":   ItemFrom,
		"where":  ItemWhere,
		"order":  ItemOrder,
		"by":     ItemBy,
		"or":     ItemOr,
		"and":    ItemAnd,
		"not":    ItemNot,
		"limit":  ItemLimit,
		"asc":    ItemAsc,
		"desc":   ItemDesc,
		"like":   ItemLike,
		"is":     ItemIs,
		"null":   ItemNull,
		"false":  ItemFalse,
		"true":   ItemTrue,
	}
)

// Item is an interface to handle the item, any item in the query
type Item interface {
	fmt.Stringer
	Type() ItemType
	Pos() int
	Value() string
	Data() int
}

type item struct {
	typ   ItemType
	pos   int
	value string
	data  int
}

func (i item) Type() ItemType {
	return i.typ
}

func (i item) Pos() int {
	return i.pos
}

func (i item) Value() string {
	return i.value
}

func (i item) Data() int {
	return i.data
}

func (i item) String() string {
	return fmt.Sprintf("pos %d, token %s", i.pos, i.value)
}

func assertTrue(in bool, msg string) {
	if !in {
		panic(msg)
	}
}

func assertType(i Item, t ItemType) {
	assertTrue(i.Type() == t, fmt.Sprintf("assertion failed, type is %d want %d", i.Type(), t))
}

type lexer struct {
	input string // input string
	start int    // start position for the current lexeme
	pos   int    // current position
	width int    // last rune width
	items chan item

	parenDepth int
	qIndex     int
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *lexer) emit(t ItemType) {
	data := 0
	if t == ItemQuestionMark {
		l.qIndex++
		data = l.qIndex
	}
	l.items <- item{t, l.start, l.input[l.start:l.pos], data}
	// Some items contain text internally. If so, count their newlines.
	l.start = l.pos
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{ItemError, l.start, fmt.Sprintf(format, args...), 0}
	return nil
}

// nextItem returns the next item from the input.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) nextItem() item {
	return <-l.items
}

// drain drains the output so the lexing goroutine will exit.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) drain() {
	for range l.items {
	}
}

// lex creates a new scanner for the input string.
func lex(input string) *lexer {
	l := &lexer{
		input: input,
		items: make(chan item, 2),
	}
	go l.run()
	return l
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for state := lexStart; state != nil; {
		state = state(l)
	}
	close(l.items)
}

func (l *lexer) parenCheck() stateFn {
	if l.parenDepth > 0 {
		l.errorf("paren not closed")
	}
	return nil
}

func lexStart(l *lexer) stateFn {
	switch r := l.peek(); {
	case r == eof:
		return l.parenCheck()
	case isSQLOperator(r):
		return lexOp
	case isAlpha(r):
		return lexAlpha
	case isNumeric(r):
		return lexNumber
	case isSpace(r):
		return lexWhiteSpace
	case r == '(':
		return lexParenOpen
	case r == ')':
		return lexParenClose
	case r == '"':
		return createLiteralFunc('"', ItemLiteral2)
	case r == '\'':
		return createLiteralFunc('\'', ItemLiteral1)
	case r == ';':
		return lexSemicolon
	case r == ',':
		return lexComma
	case r == '*':
		return lexWildCard
	case r == '.':
		return lexDot
	case r == '?':
		return lexParameter
	default:
		return l.errorf("invalid character %c", r)
	}
}

func lexWhiteSpace(l *lexer) stateFn {
	l.acceptRun(" \n\t")
	l.emit(ItemWhiteSpace)
	return lexStart
}

func lexOp(l *lexer) stateFn {
	r := l.next()
	rn := l.peek()
	var t ItemType
	switch r {
	case '>':
		t = ItemGreater
		if rn == '=' {
			l.next()
			t = ItemGreaterEqual
		}
	case '<':
		t = ItemLesser
		if rn == '=' {
			l.next()
			t = ItemLesserEqual
		} else if rn == '>' {
			l.next()
			t = ItemNotEqual
		}

	case '=':
		t = ItemEqual
	}
	l.emit(t)
	return lexStart
}

func lexParameter(l *lexer) stateFn {
	l.next()
	l.emit(ItemQuestionMark)
	return lexStart
}

func lexAlpha(l *lexer) stateFn {
	for isAlphaNumeric(l.next()) {
	}
	l.backup()
	t := strings.ToLower(l.input[l.start:l.pos])
	item := ItemAlpha
	if n, ok := keywords[t]; ok {
		item = n
	}
	l.emit(item)
	return lexStart
}

func lexParenOpen(l *lexer) stateFn {
	l.next()
	l.parenDepth++
	l.emit(ItemParenOpen)
	return lexStart
}

func lexParenClose(l *lexer) stateFn {
	l.next()
	l.parenDepth--
	if l.parenDepth < 0 {
		l.errorf("invalid ) ")
		return nil
	}
	l.emit(ItemParenClose)
	return lexStart
}

func lexSemicolon(l *lexer) stateFn {
	l.next()
	l.emit(ItemSemicolon)
	return lexStart
}

func lexComma(l *lexer) stateFn {
	l.next()
	l.emit(ItemComma)
	return lexStart
}

func createLiteralFunc(c rune, it ItemType) stateFn {
	return func(l *lexer) stateFn {
		l.next()
		var escape bool
		for {
			r := l.next()
			if escape && r != c && r != '\\' {
				l.backup() // get better error position
				return l.errorf("invalid escape character")
			}
			if r == c && !escape {
				break
			}
			if r == '\\' && !escape {
				escape = true
			} else {
				escape = false
			}
			if r == eof {
				return l.errorf("string is not terminated")
			}
		}
		l.emit(it)
		return lexStart
	}
}

func lexWildCard(l *lexer) stateFn {
	l.next()
	l.emit(ItemWildCard)
	return lexStart
}

func lexNumber(l *lexer) stateFn {
	var dot bool
	for {
		r := l.next()
		if r == '.' {
			if dot {
				return l.errorf("two dot in one number")
			}
			dot = true
			continue
		}
		if !isNumeric(r) {
			break
		}
	}
	l.backup()
	l.emit(ItemNumber)
	return lexStart
}

func lexDot(l *lexer) stateFn {
	l.next()
	l.emit(ItemDot)
	return lexStart
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}

func isAlpha(r rune) bool {
	return unicode.IsLetter(r)
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func isNumeric(r rune) bool {
	return unicode.IsDigit(r)
}

func isSQLOperator(r rune) bool {
	return r == '>' || r == '<' || r == '='
}
