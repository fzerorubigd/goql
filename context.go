package goql

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/fzerorubigd/goql/parse"
)

type fieldType int

const (
	fieldTypeColumn fieldType = iota
	fieldTypeCopy
	fieldTypeStaticNumber
	fieldTypeStaticString
	fieldTypeStaticBool
	fieldTypeFunction
	fieldTypeParameter // means this is a ? and must fill it from the query parameters
)

type field struct {
	order      int
	name       string // TODO : support for alias
	typ        fieldType
	staticStr  string // for static column only
	staticNum  float64
	copy       int   // for duplicated field, copy from another field
	argsOrder  []int // for function column only, the order of arguments in the fields list
	show       bool
	staticBool bool
	paramIdx   int // for parameters from query input
}

type context struct {
	pkg interface{}
	q   *parse.Query

	table      string
	definition map[string]columnDef
	flds       []field
	selected   map[string]int

	where parse.Stack

	order int
}

const (
	itemColumn parse.ItemType = -999
)

// A hack to handle column, I don't like this kind of hacks but I'm too bored :)
type dummy struct {
	typ   parse.ItemType
	pos   int
	value string
	data  int
}

func (d dummy) Type() parse.ItemType {
	return d.typ
}

func (d dummy) Pos() int {
	panic("here")
}

func (d dummy) Value() string {
	return d.value
}

func (d dummy) Data() int {
	return d.data
}

func (d dummy) String() string {
	return "dummy"
}

func newItem(t parse.ItemType, v string, p int, d int) parse.Item {
	return dummy{
		typ:   t,
		pos:   p,
		value: v,
		data:  d,
	}
}

// execute the query
func execute(c interface{}, src *parse.Query, params ...interface{}) ([]string, [][]Getter, error) {
	var err error
	ctx := &context{pkg: c, q: src}

	err = selectColumn(ctx)
	if err != nil {
		return nil, nil, err
	}

	return doQuery(ctx, params...)
}

func getStaticColumn(ctx *context, fl parse.Field, show bool) field {
	assertType(fl.Item, parse.ItemNumber, parse.ItemTrue, parse.ItemFalse, parse.ItemLiteral1, parse.ItemQuestionMark)
	defer func() {
		ctx.order++
	}()
	name := "static"
	var t field
	if fl.Item.Type() == parse.ItemLiteral1 {
		t = field{
			order:     ctx.order,
			name:      name,
			show:      show,
			typ:       fieldTypeStaticString,
			staticStr: parse.GetTokenString(fl.Item),
		}
	}
	if fl.Item.Type() == parse.ItemNumber {
		f, _ := strconv.ParseFloat(fl.Item.Value(), 64)
		t = field{
			order:     ctx.order,
			name:      name,
			show:      show,
			typ:       fieldTypeStaticNumber,
			staticNum: f,
		}
	}

	if fl.Item.Type() == parse.ItemTrue || fl.Item.Type() == parse.ItemFalse {
		f := fl.Item.Type() == parse.ItemTrue
		t = field{
			order:      ctx.order,
			name:       name,
			show:       show,
			typ:        fieldTypeStaticBool,
			staticBool: f,
		}
	}

	if fl.Item.Type() == parse.ItemQuestionMark {
		t = field{
			order:    ctx.order,
			name:     name,
			show:     show,
			typ:      fieldTypeParameter,
			paramIdx: fl.Item.Data(),
		}
	}
	return t
}

func getFieldColumn(ctx *context, fl parse.Field, show bool) (field, error) {
	if fl.Table != "" && fl.Table != ctx.table {
		return field{}, fmt.Errorf("table %s is not in select, join is not supported", fl.Table)
	}
	_, ok := ctx.definition[fl.Item.Value()]
	if !ok {
		return field{}, fmt.Errorf("field %s is not available in table %s", fl.Item.Value(), ctx.table)
	}

	defer func() {
		ctx.order++
	}()
	if o, sel := ctx.selected[fl.Item.Value()]; sel {
		// this is already selected, simply use the copy type for it
		return field{
			order: ctx.order,
			name:  fl.Item.Value(),
			show:  show,
			typ:   fieldTypeCopy,
			copy:  o,
		}, nil
	}

	ctx.selected[fl.Item.Value()] = ctx.order
	return field{
		order: ctx.order,
		name:  fl.Item.Value(),
		show:  show,
		typ:   fieldTypeColumn,
	}, nil
}

func getFieldStar(ctx *context) []field {
	res := make([]field, len(ctx.definition))
	for i := range ctx.definition {
		// the order in * is based on order in definition, so fix it
		f, _ := getFieldColumn(ctx, parse.Field{Item: newItem(parse.ItemAlpha, i, 0, 0)}, true)
		f.order = ctx.definition[i].Order()
		ctx.selected[f.name] = f.order
		res[f.order] = f
	}

	return res
}

func getFieldFunction(ctx *context, fl parse.Field, show bool) ([]field, error) {
	assertType(fl.Item, parse.ItemFunc)
	if !hasFunction(fl.Item.Value()) {
		return nil, fmt.Errorf("function '%s' is not registered", fl.Item.Value())
	}
	f := []field{
		{
			order: ctx.order,
			name:  fl.Item.Value(),
			show:  show,
			typ:   fieldTypeFunction,
		},
	}
	var params []field
	ctx.order++
	var err error
	params, f[0].argsOrder, err = getFields(ctx, false, fl.Parameters...)
	if err != nil {
		return nil, err
	}
	return append(f, params...), nil
}

func isStatic(t parse.ItemType) bool {
	return t == parse.ItemNumber ||
		t == parse.ItemTrue ||
		t == parse.ItemFalse ||
		t == parse.ItemNull ||
		t == parse.ItemLiteral1 ||
		t == parse.ItemQuestionMark
}

func getFields(ctx *context, show bool, fls ...parse.Field) ([]field, []int, error) {
	var direct []int
	var res []field
	for i := range fls {
		direct = append(direct, ctx.order)
		switch {
		case fls[i].Item.Type() == parse.ItemWildCard:
			if !show {
				return nil, nil, fmt.Errorf("invalid * position")
			}
			res = append(res, getFieldStar(ctx)...)
		case isStatic(fls[i].Item.Type()):
			res = append(res, getStaticColumn(ctx, fls[i], show))
		case fls[i].Item.Type() == parse.ItemFunc:
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
		fs, err := getFieldColumn(ctx, parse.Field{Item: newItem(parse.ItemAlpha, o[i].Field, 0, 0)}, false)
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
			case parse.ItemAlpha, parse.ItemLiteral2:
				v := parse.GetTokenString(p)
				f, err := getFieldColumn(ctx, parse.Field{Item: newItem(parse.ItemAlpha, v, 0, 0)}, false)
				if err != nil {
					return nil, nil, err
				}
				res = append(res, f)
				ts = newItem(itemColumn, v, 0, f.order)
			case parse.ItemQuestionMark:
				res = append(res, field{
					order:    ctx.order,
					name:     "static",
					show:     false,
					typ:      fieldTypeParameter,
					paramIdx: p.Data(),
				})
				ts = newItem(itemColumn, "static", 0, ctx.order)
				ctx.order++
			case parse.ItemFunc:
				fn := p.(parse.FuncItem)
				fs, err := getFieldFunction(ctx, parse.Field{Item: p, Parameters: fn.Parameters()}, false)
				if err != nil {
					return nil, nil, err
				}
				res = append(res, fs...)
				// function is calculated on fill gaps, then simply we can use the result at the end a static
				ts = newItem(itemColumn, fn.Value(), 0, fs[0].order)
			}
			s.Push(ts)
		}
	}
	return s, res, nil

}

func selectColumn(ctx *context) error {
	ss := ctx.q.Statement.(*parse.SelectStmt)
	tbl, err := getTable(ss.Table)
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

func filterColumn(ctx *context, items ...Getter) []Getter {
	fl := ctx.flds
	res := make([]Getter, 0, len(items))
	for i := range fl {
		if fl[i].show {
			res = append(res, items[i])
		}
	}

	return res
}

func fillGaps(ctx *context, res []Getter) error {
	fl := ctx.flds
	for i := range fl {
		switch fl[i].typ {
		case fieldTypeCopy:
			res[fl[i].order] = res[fl[i].copy]
		case fieldTypeStaticNumber:
			res[i] = Number{Number: fl[i].staticNum}
		case fieldTypeStaticString:
			res[i] = String{String: fl[i].staticStr}
		case fieldTypeStaticBool:
			res[i] = Bool{Bool: fl[i].staticBool}
		}
	}

	var err error
	// once more for functions :/ if there is a way to fill it in one loop :/
	// TODO : exec this at the getTableFields not after that
	// since functions added before its parameter, doing it backward to make sure if a function is argument to another function
	// is field before
	for i := len(fl) - 1; i > -1; i-- {
		if fl[i].typ == fieldTypeFunction {
			args := make([]Getter, len(fl[i].argsOrder))
			for j, v := range fl[i].argsOrder {
				args[j] = res[v]
			}
			res[fl[i].order], err = executeFunction(fl[i].name, args...)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func callWhere(where getter, i []Getter) (ok bool, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("error : %v", e)
			ok = false
		}
	}()
	return toBool(where(i)), nil
}

func paramToStatic(f *field, p interface{}) {

	switch t := p.(type) {
	case bool:
		f.typ = fieldTypeStaticBool
		f.staticBool = t
	case string:
		f.typ = fieldTypeStaticString
		f.staticStr = t
	case fmt.Stringer:
		f.typ = fieldTypeStaticString
		f.staticStr = t.String()
	case int:
		f.typ = fieldTypeStaticNumber
		f.staticNum = float64(t)
	case int64:
		f.typ = fieldTypeStaticNumber
		f.staticNum = float64(t)
	case float32:
		f.typ = fieldTypeStaticNumber
		f.staticNum = float64(t)
	case float64:
		f.typ = fieldTypeStaticNumber
		f.staticNum = t
	default:
		f.typ = fieldTypeStaticString
		f.staticStr = fmt.Sprintf("%T", t)
	}
}

func doQuery(ctx *context, params ...interface{}) ([]string, [][]Getter, error) {
	res := make(chan []Getter, 3)
	ss := ctx.q.Statement.(*parse.SelectStmt)
	var all = make([]string, len(ctx.flds))

	for i := range ctx.flds {
		// only fields are allowed
		if ctx.flds[i].typ == fieldTypeColumn {
			all[i] = ctx.flds[i].name
		}
		if ctx.flds[i].typ == fieldTypeParameter {
			paramToStatic(&ctx.flds[i], params[ctx.flds[i].paramIdx-1])
		}
	}
	quit := make(chan struct{}, 1)
	defer func() {
		quit <- struct{}{}
		close(quit) // prevent the channel leak.
	}()
	err := getTableFields(ctx.pkg, ss.Table, res, quit, all...)
	if err != nil {
		return nil, nil, err
	}
	where, err := buildFilter(ctx.where)
	if err != nil {
		return nil, nil, err
	}
	a := make([][]Getter, 0)
	for i := range res {
		if err = fillGaps(ctx, i); err != nil {
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
			a = [][]Getter{}
		} else if ss.Start+ss.Count >= l {
			a = a[ss.Start:] // to the end
		} else {
			a = a[ss.Start : ss.Start+ss.Count]
		}
	}

	return column, a, nil
}
