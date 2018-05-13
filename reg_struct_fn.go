package goql

import (
	"reflect"

	"github.com/fzerorubigd/goql/astdata"
)

type isStructFn int

func (isStructFn) Execute(args ...Getter) (Getter, error) {
	def, err := getSingleDef(args...)
	if err != nil {
		return nil, err
	}
	if def == nil {
		return Bool{Null: true}, nil
	}

	_, ok := def.(*astdata.StructType)
	return Bool{Bool: ok}, nil
}

type structFieldCountFn int

func (sfc structFieldCountFn) Execute(args ...Getter) (Getter, error) {
	def, err := getSingleDef(args...)
	if err != nil {
		return nil, err
	}
	if def == nil {
		return Number{Null: true}, nil
	}

	st, ok := def.(*astdata.StructType)
	if !ok {
		return Number{Null: true}, nil
	}
	return Number{Number: float64(len(st.Fields()))}, nil
}

type structFieldDefFn int

func (structFieldDefFn) Execute(args ...Getter) (Getter, error) {
	if err := required(2, 2, args...); err != nil {
		return nil, err
	}

	def := toDefinition(args[0].Get())
	st, ok := def.(*astdata.StructType)
	if !ok {
		return Definition{}, nil
	}

	var f astdata.Definition
	fl := st.Fields()
bigSwitch:
	switch t := args[1].Get().(type) {
	case float64:
		nm := int(t)
		if len(fl) < nm || nm < 1 {
			return Definition{}, nil
		}
		f = fl[nm-1].Definition()
	case string:
		for i := range fl {
			if fl[i].Name() == t {
				f = fl[i].Definition()
				break bigSwitch
			}
		}
		return Definition{}, nil
	default:
		return Definition{}, nil
	}
	return Definition{Definition: f}, nil
}

type structFieldNameFn int

func (structFieldNameFn) Execute(args ...Getter) (Getter, error) {
	if err := required(2, 2, args...); err != nil {
		return nil, err
	}

	def := toDefinition(args[0].Get())
	nm := int(toNumber(args[1].Get()))
	st, ok := def.(*astdata.StructType)
	if !ok {
		return String{Null: true}, nil
	}
	fl := st.Fields()
	if len(fl) < nm || nm < 1 {
		return String{Null: true}, nil
	}
	return String{String: fl[nm-1].Name()}, nil
}

type structFieldTagFn int

func (sft structFieldTagFn) Execute(args ...Getter) (Getter, error) {
	if err := required(2, 3, args...); err != nil {
		return nil, err
	}

	def := toDefinition(args[0].Get())
	st, ok := def.(*astdata.StructType)
	if !ok {
		return String{Null: true}, nil
	}
	var f reflect.StructTag
	if sft == 0 {
		fl := st.Fields()
	bigSwitch:
		switch t := args[1].Get().(type) {
		case float64:
			nm := int(t)
			if len(fl) < nm || nm < 1 {
				return String{Null: true}, nil
			}
			f = fl[nm-1].Tags()
		case string:
			for i := range fl {
				if fl[i].Name() == t {
					f = fl[i].Tags()
					break bigSwitch
				}
			}
			return String{Null: true}, nil
		default:
			return String{Null: true}, nil
		}
	} else {
		fl := st.Embeds()
		switch t := args[1].Get().(type) {
		case float64:
			nm := int(t)
			if len(fl) < nm || nm < 1 {
				return String{Null: true}, nil
			}
			f = fl[nm-1].Tags()
		default:
			return String{Null: true}, nil
		}
	}
	if len(args) == 3 {
		tag := toString(args[2].Get())
		return String{String: f.Get(tag)}, nil
	}
	return String{String: string(f)}, nil
}

func registerStructFunc() {
	// Struct funcs
	RegisterFunction("is_struct", isStructFn(0))
	RegisterFunction("field_def", structFieldDefFn(0))
	RegisterFunction("field_name", structFieldNameFn(0))
	RegisterFunction("field_count", structFieldCountFn(0))
	RegisterFunction("field_tag", structFieldTagFn(0))
	RegisterFunction("embed_tag", structFieldTagFn(1))
}

func init() {
	registerStructFunc()
}
