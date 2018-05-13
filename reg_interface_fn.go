package goql

import (
	"github.com/fzerorubigd/goql/astdata"
)

type isInterfaceFn int

func (isInterfaceFn) Execute(args ...Getter) (Getter, error) {
	def, err := getSingleDef(args...)
	if err != nil {
		return nil, err
	}
	if def == nil {
		return Bool{Null: true}, nil
	}

	_, ok := def.(*astdata.InterfaceType)
	return Bool{Bool: ok}, nil

}

type interfaceFieldCountFn int

func (ifc interfaceFieldCountFn) Execute(args ...Getter) (Getter, error) {
	def, err := getSingleDef(args...)
	if err != nil {
		return nil, err
	}
	if def == nil {
		return Number{Null: true}, nil
	}

	st, ok := def.(*astdata.InterfaceType)
	if !ok {
		return Number{Null: true}, nil
	}

	return Number{Number: float64(len(st.Functions()))}, nil
}

type interfaceFuncDef int

func (ifc interfaceFuncDef) Execute(args ...Getter) (Getter, error) {
	if err := required(2, 2, args...); err != nil {
		return nil, err
	}

	def := toDefinition(args[0].Get())
	st, ok := def.(*astdata.InterfaceType)
	if !ok {
		return Definition{}, nil
	}

	switch t := args[1].Get().(type) {
	case float64:
		i := int(t)
		fn := st.Functions()
		if i > len(fn) || i < 1 {
			return Definition{}, nil
		}

		return Definition{Definition: fn[i-1].Definition()}, nil
	case string:
		fn := st.Functions()
		for i := range fn {
			if fn[i].Name() == t {
				return Definition{Definition: fn[i].Definition()}, nil
			}
		}
	}
	return Definition{}, nil
}

type interfaceFuncName int

func (ifc interfaceFuncName) Execute(args ...Getter) (Getter, error) {
	if err := required(2, 2, args...); err != nil {
		return nil, err
	}

	def := toDefinition(args[0].Get())
	st, ok := def.(*astdata.InterfaceType)
	if !ok {
		return String{Null: true}, nil
	}

	i := int(toNumber(args[1].Get()))
	fn := st.Functions()
	if i > len(fn) || i < 1 {
		return String{Null: true}, nil
	}

	return String{String: fn[i-1].Name()}, nil

}

func registerInterfaceFunc() {
	RegisterFunction("is_interface", isInterfaceFn(0))
	RegisterFunction("func_count", interfaceFieldCountFn(0))
}

func init() {
	registerInterfaceFunc()
}
