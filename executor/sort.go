package executor

import (
	"strings"

	"github.com/fzerorubigd/goql/internal/parse"
	"github.com/fzerorubigd/goql/structures"
)

type sortMe struct {
	order parse.Orders
	data  [][]structures.Valuer
}

func (s *sortMe) Len() int {
	return len(s.data)
}

func interfaceLess(i, j interface{}, desc bool) int {
	lesser := 1
	greater := -1
	equal := 0
	if desc {
		lesser = -1
		greater = 1
	}
	if i == nil && j != nil {
		return lesser
	}
	if i != nil && j == nil {
		return greater
	}
	if i == nil && j == nil {
		return equal
	}
	// i, j are not nil
	switch i.(type) {
	case bool:
		if i.(bool) == j.(bool) {
			return equal
		}
		if i.(bool) {
			return lesser
		}
		return greater
	case float64:
		if i.(float64) == j.(float64) {
			return equal
		}
		if i.(float64) < j.(float64) {
			return lesser
		}
		return greater
	case string:
		if i.(string) == j.(string) {
			return equal
		}
		i := strings.Compare(strings.ToLower(i.(string)), strings.ToLower(j.(string)))
		if i < 0 {
			return lesser
		}
		return greater
	}

	return equal
}

func (s *sortMe) Less(i int, j int) bool {
	ii := s.data[i]
	ij := s.data[j]
	for o := range s.order {
		res := interfaceLess(ii[s.order[o].Index].Value(), ij[s.order[o].Index].Value(), s.order[o].DESC)
		if res < 0 {
			return false
		} else if res > 0 {
			return true
		}
	}

	return false
}

func (s *sortMe) Swap(i int, j int) {
	s.data[i], s.data[j] = s.data[j], s.data[i]
}
