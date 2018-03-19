package executor

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/fzerorubigd/goql/internal/parse"
	"github.com/fzerorubigd/goql/structures"
)

type fieldType int

const (
	fieldTypeColumn fieldType = iota
	fieldTypeCopy
	fieldTypeStaticNumber
	fieldTypeStaticString
	fieldTypeFunction
)

type field struct {
	order     int
	name      string // TODO : support for alias
	show      bool
	typ       fieldType
	staticStr string // for static column only
	staticNum float64
	copy      int   // for duplicated field, copy from another field
	argsOrder []int // for function column only, the order of arguments in the fields list
}

type context struct {
	pkg interface{}
	q   *parse.Query

	table      string
	definition map[string]structures.ColumnDef
	flds       []field
	selected   map[string]int

	where parse.Stack

	order int
}

const itemColumn parse.ItemType = -999

// A hack to handle column, I don't like this kind of hacks but I'm too bored :)
type col struct {
	field string
	index int
}

func (col) Type() parse.ItemType {
	return itemColumn
}

// Pos is the position of the column in the requested index
func (c col) Pos() int {
	return c.index
}

func (c col) Value() string {
	return c.field
}

func (c col) String() string {
	return ""
}

// Execute the query
func Execute(c interface{}, src *parse.Query) ([]string, [][]structures.Valuer, error) {
	var err error
	ctx := &context{pkg: c, q: src}

	err = selectColumn2(ctx)
	if err != nil {
		return nil, nil, err
	}

	return doQuery(ctx)
}

func getStaticColumn(ctx *context, fl parse.Field, show bool) field {
	defer func() {
		ctx.order++
	}()
	name := fmt.Sprintf("COL_%d", ctx.order+1)
	if fl.String != "" {
		return field{
			order:     ctx.order,
			name:      name,
			show:      show,
			typ:       fieldTypeStaticString,
			staticStr: fl.String,
		}
	}
	if fl.Number != "" {
		f, _ := strconv.ParseFloat(fl.Number, 64)
		return field{
			order:     ctx.order,
			name:      fmt.Sprintf("COL_%d", ctx.order+1),
			show:      show,
			typ:       fieldTypeStaticNumber,
			staticNum: f,
		}
	}

	panic("runtime error")
}

func getFieldColumn(ctx *context, fl parse.Field, show bool) (field, error) {
	if fl.Table != "" && fl.Table != ctx.table {
		return field{}, fmt.Errorf("table %s is not in select, join is not supported", fl.Table)
	}
	_, ok := ctx.definition[fl.Column]
	if !ok {
		return field{}, fmt.Errorf("field %s is not available in table %s", fl.Column, ctx.table)
	}

	defer func() {
		ctx.order++
	}()
	if o, sel := ctx.selected[fl.Column]; sel {
		// this is already selected, simply use the copy type for it
		return field{
			order: ctx.order,
			name:  fl.Column,
			show:  show,
			typ:   fieldTypeCopy,
			copy:  o,
		}, nil
	}

	ctx.selected[fl.Column] = ctx.order
	return field{
		order: ctx.order,
		name:  fl.Column,
		show:  show,
		typ:   fieldTypeColumn,
	}, nil
}

func getFieldStar(ctx *context) []field {
	res := make([]field, len(ctx.definition))
	for i := range ctx.definition {
		res[ctx.definition[i].Order()], _ = getFieldColumn(ctx, parse.Field{Column: i}, true)
	}
	return res
}

func getFieldFunction(ctx *context, fl parse.Field, show bool) ([]field, error) {
	if fl.Function == nil {
		panic("runtime error")
	}
	if !structures.HasFunction(fl.Function.Name) {
		return nil, fmt.Errorf("function '%s' is not registered", fl.Function.Name)
	}
	f := []field{
		field{
			order: ctx.order,
			name:  fl.Function.Name,
			show:  show,
			typ:   fieldTypeFunction,
		},
	}
	var params []field
	ctx.order++
	var err error
	params, f[0].argsOrder, err = getFields(ctx, false, fl.Function.Parameters...)
	if err != nil {
		return nil, err
	}
	return append(f, params...), nil
}

func getFields(ctx *context, show bool, fls ...parse.Field) ([]field, []int, error) {
	var direct = make([]int, len(fls))
	var res []field
	for i := range fls {
		direct = append(direct, ctx.order)
		switch {
		case fls[i].WildCard:
			if !show {
				return nil, nil, fmt.Errorf("invalid * position")
			}
			res = append(res, getFieldStar(ctx)...)
		case fls[i].String != "" || fls[i].Number != "":
			res = append(res, getStaticColumn(ctx, fls[i], show))
		case fls[i].Function != nil:
			fs, err := getFieldFunction(ctx, fls[i], show)
			if err != nil {
				return nil, nil, err
			}
			res = append(res, fs...)
		default:
			f, err := getFieldColumn(ctx, fls[i], show)
			if err != nil {
				return nil, nil, err
			}
			res = append(res, f)
		}
	}

	return res, direct, nil
}

func getOrderFileds(ctx *context, o ...parse.Order) ([]field, error) {
	var res []field
	for i := range o {
		fs, err := getFieldColumn(ctx, parse.Field{Column: o[i].Field}, false)
		if err != nil {
			return nil, err
		}
		res = append(res, fs)
	}
	return res, nil
}

func getWhereField(ctx *context) (parse.Stack, []field, error) {
	ss := ctx.q.Statement.(*parse.SelectStmt)
	var res []field
	s := parse.NewStack(0)
	// which column are needed in where?
	if st := ss.Where; st != nil {
		for {
			p, err := st.Pop()
			if err != nil {
				break
			}
			ts := p
			switch p.Type() {
			case parse.ItemAlpha:
				// this must be a column name
				v := strings.ToLower(p.Value())
				if v == "null" || v == "true" || v == "false" {
					break
				}
				fallthrough
			case parse.ItemLiteral2:
				v := parse.GetTokenString(p)
				f, err := getFieldColumn(ctx, parse.Field{Column: v}, false)
				if err != nil {
					return nil, nil, err
				}
				res = append(res, f)
				ts = col{
					index: f.order,
					field: v,
				}
			}
			s.Push(ts)
		}
	}
	return s, res, nil

}

func selectColumn2(ctx *context) error {
	ss := ctx.q.Statement.(*parse.SelectStmt)
	tbl, err := structures.GetTable(ss.Table)
	if err != nil {
		return err
	}

	ctx.table = ss.Table
	ctx.definition = tbl
	ctx.order = 0
	ctx.selected = make(map[string]int)
	ctx.where = parse.NewStack(0)

	// fields after select
	fl, _, err := getFields(ctx, true, ss.Fields...)
	if err != nil {
		return err
	}

	ctx.flds = fl

	// order fields
	fl, err = getOrderFileds(ctx, ss.Order...)
	if err != nil {
		return err
	}

	ctx.flds = append(ctx.flds, fl...)

	// where fields

	ctx.where, fl, err = getWhereField(ctx)
	if err != nil {
		return err
	}

	ctx.flds = append(ctx.flds, fl...)
	return nil
}

func filterColumn(ctx *context, items ...structures.Valuer) []structures.Valuer {
	fl := ctx.flds
	res := make([]structures.Valuer, 0, len(items))
	for i := range fl {
		if fl[i].show {
			res = append(res, items[i])
		}
	}

	return res
}

func fillGaps(ctx *context, res []structures.Valuer) error {
	fl := ctx.flds
	for i := range fl {
		switch fl[i].typ {
		case fieldTypeCopy:
			res[i] = res[fl[i].copy]
		case fieldTypeStaticNumber:
			res[i] = structures.Number{Number: fl[i].staticNum}
		case fieldTypeStaticString:
			res[i] = structures.String{String: fl[i].staticStr}
		}
	}

	var err error
	// once more for functions :/ if there is a way to fill it in one loop :/
	// TODO : exec this at the GetFields not after that
	for i := range fl {
		if fl[i].typ == fieldTypeFunction {
			args := make([]structures.Valuer, len(fl[i].argsOrder))
			for j := range fl[i].argsOrder {
				args[j] = res[fl[i].argsOrder[j]]
			}

			res[i], err = structures.ExecuteFunction(fl[i].name, args...)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func callWhere(where getter, i []structures.Valuer) (ok bool, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("error : %v", e)
			ok = false
		}
	}()
	return toBool(where(i)), nil
}

func doQuery(ctx *context) ([]string, [][]structures.Valuer, error) {
	res := make(chan []structures.Valuer, 3)
	ss := ctx.q.Statement.(*parse.SelectStmt)
	var all = make([]string, len(ctx.flds))
	for i := range ctx.flds {
		// only fields are allowed
		if ctx.flds[i].typ == fieldTypeColumn {
			all[i] = ctx.flds[i].name
		}
	}

	err := structures.GetFields(ctx.pkg, ss.Table, res, all...)
	if err != nil {
		return nil, nil, err
	}
	where, err := buildFilter(ctx.where)
	if err != nil {
		return nil, nil, err
	}
	a := make([][]structures.Valuer, 0)
	for i := range res {
		if err = fillGaps(ctx, i); err != nil {
			close(res) // prevent the channel leak. TODO : better way
			return nil, nil, err
		}
		ok, err := callWhere(where, i)
		if err != nil {
			return nil, nil, err
		}
		if !ok {
			continue
		}

		a = append(a, filterColumn(ctx, i...))
	}

	column := make([]string, 0, len(ctx.flds))
	for i := range ctx.flds {
		if ctx.flds[i].show {
			column = append(column, ctx.flds[i].name)
		}
	}

	// sort
	s := &sortMe{
		data:  a,
		order: ss.Order,
	}
	sort.Sort(s)

	a = s.data
	if ss.Count >= 0 && ss.Start >= 0 {
		l := len(a)
		if ss.Start >= l {
			a = [][]structures.Valuer{}
		} else if ss.Start+ss.Count >= l {
			a = a[ss.Start:] // to the end
		} else {
			a = a[ss.Start : ss.Start+ss.Count]
		}
	}

	return column, a, nil
}
