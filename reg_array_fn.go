package goql

import "github.com/fzerorubigd/goql/astdata"

type isArrayFn int

func (isArrayFn) Execute(args ...Getter) (Getter, error) {
	def, err := getSingleDef(args...)
	if err != nil {
		return nil, err
	}
	if def == nil {
		return Bool{Null: true}, nil
	}

	_, ok1 := def.(*astdata.ArrayType)
	_, ok2 := def.(*astdata.EllipsisType)

	return Bool{Bool: ok1 || ok2}, nil
}

type arrayItemFn int

func (arrayItemFn) Execute(args ...Getter) (Getter, error) {
	def, err := getSingleDef(args...)
	if err != nil {
		return nil, err
	}
	if def == nil {
		return Definition{}, nil
	}

	a, ok := def.(*astdata.ArrayType)
	if ok {
		return Definition{Definition: a.ValueDefinition()}, nil
	}
	e, ok := def.(*astdata.EllipsisType)
	if ok {
		return Definition{Definition: e.ValueDefinition()}, nil
	}

	return Definition{}, nil
}

type arrayTypeFn int

func (arrayTypeFn) Execute(args ...Getter) (Getter, error) {
	def, err := getSingleDef(args...)
	if err != nil {
		return nil, err
	}
	if def == nil {
		return String{Null: true}, nil
	}

	if arr, ok := def.(*astdata.ArrayType); ok {
		if arr.Slice() {
			return String{String: "slice"}, nil // enum?
		}
		return String{String: "array"}, nil // enum?
	}
	if _, ok := def.(*astdata.EllipsisType); ok {
		return String{String: "ellipsis"}, nil // enum?
	}

	return String{Null: true}, nil

}

func registerArrayFunc() {

	// Map functions
	RegisterFunction("is_array", isArrayFn(0))
	RegisterFunction("array_type", arrayTypeFn(0))
	RegisterFunction("array_def", arrayItemFn(0))
}

func init() {
	registerArrayFunc()
}
