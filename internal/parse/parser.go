package parse

import "fmt"

type parser struct {
	l        *lexer
	last     item
	rejected bool
}

func (p *parser) scan() item {
	// there is one rejected item
	if p.rejected {
		p.rejected = false
		return p.last
	}

	p.last = p.l.nextItem()
	p.rejected = false
	return p.last
}

func (p *parser) scanIgnoreWhiteSpace() item {
	t := p.scan()
	if t.typ == itemWhiteSpace {
		t = p.scan()
	}
	return t
}

func (p *parser) reject() {
	p.rejected = true
}

// AST return the abstract source tree for given query
func AST(q string) (Query, error) {
	fmt.Println(q)
	p := &parser{
		l: lex(q),
	}
	s, err := newStatement(p)
	if err != nil {
		p.l.drain() // make sure the lexer is terminated
		return Query{}, err
	}

	return Query{
		Statement: s,
	}, nil
}
