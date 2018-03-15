package structures

import (
	"fmt"
	"strings"
	"sync"
)

// ValueType is the value type of query
type ValueType int

// String is the string type in our system
type String struct {
	String string
	Null   bool
}

// Number is the number
type Number struct {
	Number float64
	Null   bool
}

// Bool is the boolean type
type Bool struct {
	Bool bool
	Null bool
}

const (
	// ValueTypeString is the string type
	ValueTypeString ValueType = iota
	// ValueTypeNumber is the number type
	ValueTypeNumber
	// ValueTypeBool is the bool type
	ValueTypeBool
)

var (
	tables = make(map[string]*table)
	lock   = &sync.Mutex{}
)

// StringValuer is provider for a value for a table
type StringValuer interface {
	Value(interface{}) String
}

// IntValuer is the integer valuer
type IntValuer interface {
	Value(interface{}) Number
}

// BoolValuer is the Boolean valuer
type BoolValuer interface {
	Value(interface{}) Bool
}

// Table is the single table in system
type table struct {
	name   string
	fields map[string]interface{} // interface is one of the Valuer interface and not anything else
	data   TableData
	lock   *sync.Mutex
}

// TableData is a callback to get table data from a package
type TableData func(interface{}) []interface{}

// RegisterTable is the function to handle registration of a table
func RegisterTable(name string, data TableData) {
	lock.Lock()
	defer lock.Unlock()

	if _, ok := tables[name]; ok {
		panic(fmt.Sprintf("table with name %s is already registered", name))
	}
	tables[name] = &table{
		name:   name,
		data:   data,
		fields: make(map[string]interface{}),
		lock:   &sync.Mutex{},
	}
}

// GetTable return the table definition
func GetTable(t string) (map[string]ValueType, error) {
	tbl, ok := tables[t]
	if !ok {
		panic(fmt.Sprintf("table %s is not available", t))
	}
	res := make(map[string]ValueType)
	for i := range tbl.fields {
		switch tbl.fields[i].(type) {
		case BoolValuer:
			res[i] = ValueTypeBool
		case IntValuer:
			res[i] = ValueTypeNumber
		case StringValuer:
			res[i] = ValueTypeString
		}
	}
	return res, nil
}

// RegisterField is the field registration
func RegisterField(t string, name string, valuer interface{}) {
	tbl, ok := tables[t]
	if !ok {
		panic(fmt.Sprintf("table %s is not available", t))
	}

	if _, ok := tbl.fields[name]; ok {
		panic(fmt.Sprintf("table %s is already have field %s", t, name))
	}

	switch valuer.(type) {
	case BoolValuer:
	case IntValuer:
	case StringValuer:
	default:
		panic(fmt.Sprintf("valuer is not a valid valuer, its is %T", valuer))
	}

	tbl.fields[name] = valuer
}

// GetFields is the get field fro a table
func GetFields(p interface{}, t string, res chan<- []interface{}, fields ...string) error {
	lock.Lock()
	defer lock.Unlock()
	tbl, ok := tables[t]
	if !ok {
		return fmt.Errorf("invalid table name %s", t)
	}

	if len(fields) == 0 {
		return fmt.Errorf("no field selected")
	}

	var invalid []string
	for i := range fields {
		if _, ok := tbl.fields[fields[i]]; !ok {
			invalid = append(invalid, fields[i])
		}
	}
	if len(invalid) > 0 {
		return fmt.Errorf("invalid field(s) : %s", strings.Join(invalid, ", "))
	}

	// do concurrently
	go func() {
		defer close(res)
		cache := tbl.data(p)
		for i := range cache {
			n := make([]interface{}, len(fields))
			for f := range fields {
				switch t := tbl.fields[fields[f]].(type) {
				case StringValuer:
					n[f] = t.Value(cache[i])
				case IntValuer:
					n[f] = t.Value(cache[i])
				case BoolValuer:
					n[f] = t.Value(cache[i])
				}
			}
			res <- n
		}
	}()
	return nil
}
