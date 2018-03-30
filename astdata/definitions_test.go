package astdata

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefinition(t *testing.T) {
	require.Nil(t, newType(nil, nil, nil))
}
