package astdata

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var testPkg = `
package example

type HHH int

func (h HHH) Test() {

}
`

func TestUtilityFuncs(t *testing.T) {
	p := &Package{}
	f, err := ParseFile(testPkg, p)
	require.NoError(t, err)
	p.files = append(p.files, f)

	_, err = p.FindImport("_")
	require.Error(t, err)

	_, err = p.FindImport("context")
	require.Error(t, err)

	_, err = p.FindType("NotExist")
	require.Error(t, err)

	_, err = p.FindConstant("NotExist")
	require.Error(t, err)

	_, err = p.FindFunction("NotExist")
	require.Error(t, err)

	_, err = p.FindMethod("n", "NotExist")
	require.Error(t, err)

	_, err = p.FindVariable("NotExist")
	require.Error(t, err)

}
