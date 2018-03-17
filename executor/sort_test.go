package executor

import (
	"sort"
	"testing"

	"github.com/fzerorubigd/goql/internal/parse"
	"github.com/fzerorubigd/goql/structures"
	"github.com/stretchr/testify/assert"
)

type unknown struct{}

func (u unknown) Value() interface{} {
	return u // invalid type
}

var sortData = [][]structures.Valuer{
	{
		structures.Number{Number: 0},
		structures.Bool{Bool: true},
		structures.String{String: "a"},
		structures.Number{Number: 4},
		structures.Number{Number: 4},
		unknown{},
	},
	{
		structures.Number{Number: 1},
		structures.Bool{Bool: false},
		structures.String{String: "b"},
		structures.Number{Number: 3},
		structures.Number{Null: true},
		unknown{},
	},
	{
		structures.Number{Number: 2},
		structures.Bool{Bool: true},
		structures.String{String: "c"},
		structures.Number{Number: 3},
		structures.Number{Number: 3},
		unknown{},
	},
	{
		structures.Number{Number: 3},
		structures.Bool{Bool: false},
		structures.String{String: "d"},
		structures.Number{Number: 1},
		structures.Number{Null: true},
		unknown{},
	},
	{
		structures.Number{Number: 4},
		structures.Bool{Bool: true},
		structures.String{String: "a"},
		structures.Number{Number: 10},
		structures.Number{Number: 5},
		unknown{},
	},
}

func order(in [][]structures.Valuer) []float64 {
	res := make([]float64, len(in))
	for i := range in {
		res[i] = in[i][0].Value().(float64)
	}

	return res
}

func TestSort(t *testing.T) {
	s := &sortMe{
		order: parse.Orders{
			parse.Order{
				Index: 2,
			},
		},
		data: sortData,
	}

	sort.Sort(s)
	assert.Equal(t, []float64{0, 4, 1, 2, 3}, order(s.data))

	s = &sortMe{
		order: parse.Orders{
			parse.Order{
				Index: 2,
				DESC:  true,
			},
		},
		data: sortData,
	}
	sort.Sort(s)
	assert.Equal(t, []float64{3, 2, 1, 0, 4}, order(s.data))

	s = &sortMe{
		order: parse.Orders{
			parse.Order{
				Index: 1,
			},
		},
		data: sortData,
	}

	sort.Sort(s)
	assert.Equal(t, []float64{2, 0, 4, 3, 1}, order(s.data))

	s = &sortMe{
		order: parse.Orders{
			parse.Order{
				Index: 1,
				DESC:  true,
			},
		},
		data: sortData,
	}
	sort.Sort(s)
	assert.Equal(t, []float64{3, 1, 2, 0, 4}, order(s.data))

	s = &sortMe{
		order: parse.Orders{
			parse.Order{
				Index: 0,
			},
		},
		data: sortData,
	}

	sort.Sort(s)
	assert.Equal(t, []float64{0, 1, 2, 3, 4}, order(s.data))

	s = &sortMe{
		order: parse.Orders{
			parse.Order{
				Index: 0,
				DESC:  true,
			},
		},
		data: sortData,
	}
	sort.Sort(s)
	assert.Equal(t, []float64{4, 3, 2, 1, 0}, order(s.data))
	//////
	s = &sortMe{
		order: parse.Orders{
			parse.Order{
				Index: 3,
			},
			parse.Order{
				Index: 0,
			},
		},
		data: sortData,
	}

	sort.Sort(s)
	assert.Equal(t, []float64{3, 1, 2, 0, 4}, order(s.data))

	s = &sortMe{
		order: parse.Orders{
			parse.Order{
				Index: 3,
				DESC:  true,
			},
			parse.Order{
				Index: 0,
			},
		},
		data: sortData,
	}
	sort.Sort(s)
	assert.Equal(t, []float64{4, 0, 1, 2, 3}, order(s.data))
	//////
	s = &sortMe{
		order: parse.Orders{
			parse.Order{
				Index: 4,
			},
		},
		data: sortData,
	}

	sort.Sort(s)
	assert.Equal(t, []float64{1, 3, 2, 0, 4}, order(s.data))

	s = &sortMe{
		order: parse.Orders{
			parse.Order{
				Index: 4,
				DESC:  true,
			},
		},
		data: sortData,
	}
	sort.Sort(s)
	assert.Equal(t, []float64{4, 0, 2, 1, 3}, order(s.data))
	////
	s = &sortMe{
		order: parse.Orders{
			parse.Order{
				Index: 5,
			},
			parse.Order{
				Index: 0,
			},
		},
		data: sortData,
	}

	sort.Sort(s)
	assert.Equal(t, []float64{0, 1, 2, 3, 4}, order(s.data))

}
