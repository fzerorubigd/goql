package goql

import "github.com/fzerorubigd/goql/astdata"

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

func registerInterfaceFunc() {
	RegisterFunction("is_interface", isInterfaceFn(0))
}

func init() {
	registerInterfaceFunc()
}
