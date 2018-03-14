package astdata

import (
	"go/ast"
	"strings"
)

// Docs is the code documents
type Docs []string

// String convert object to go docs again
func (d Docs) String() string {
	return strings.Join(d, "\n")
}

func docsFromNodeDoc(cgs ...*ast.CommentGroup) Docs {
	var res = Docs{} // not nil, an empty array
	for _, cg := range cgs {
		if cg != nil {
			for i := range cg.List {
				res = append(res, cg.List[i].Text)
			}
		}
	}
	return res
}
