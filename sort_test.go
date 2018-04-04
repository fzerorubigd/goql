package goql

import (
	"sort"
	"testing"

	"github.com/fzerorubigd/goql/parse"
	"github.com/stretchr/testify/assert"
)

type unknown struct{}

func (u unknown) Get() interface{} {
	return u // invalid type
}

var sortData = [][]Getter{
	{
		Number{Number: 0},
		Bool{Bool: true},
		String{String: "a"},
		Number{Number: 4},
		Number{Number: 4},
		unknown{},
	},
	{
		Number{Number: 1},
		Bool{Bool: false},
		String{String: "b"},
		Number{Number: 3},
		Number{Null: true},
		unknown{},
	},
	{
		Number{Number: 2},
		Bool{Bool: true},
		String{String: "c"},
		Number{Number: 3},
		Number{Number: 3},
		unknown{},
	},
	{
		Number{Number: 3},
		Bool{Bool: false},
		String{String: "d"},
		Number{Number: 1},
		Number{Null: true},
		unknown{},
	},
	{
		Number{Number: 4},
		Bool{Bool: true},
		String{String: "a"},
		Number{Number: 10},
		Number{Number: 5},
		unknown{},
	},
}

func order(in [][]Getter) []float64 {
	res := make([]float64, len(in))
	for i := range in {
		res[i] = in[i][0].Get().(float64)
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
