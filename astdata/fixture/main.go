// Package fixture is a fixture used for testing astdata package
// do not remove/change the function unless the tests are fixed
package fixture

import (
	// comment ctx
	ctx "context"
	// comment http
	"net/http"
)

const (
	// X Docs
	X = iota
	// Y Docs
	Y
	// Z Docs
	Z
)

const testConst = 10

// all
var (
	// hi
	hi  string
	ok  = true
	bye alpha
)

// aaa
type (
	// alpha comment
	alpha string

	beta struct {
		A int
		B string
	}
)

// Multi line
// comment
// here
func test(c ctx.Context, r *http.Request, w http.ResponseWriter) {

}

func (beta) assert(err error) {
	if err != nil {
		panic(err)
	}
}

func (a *alpha) testing() {
	panic("hi")
}
