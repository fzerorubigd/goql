package goql

import "github.com/fzerorubigd/goql/astdata"

type isMapFn int

func (isMapFn) Execute(args ...Getter) (Getter, error) {
	def, err := getSingleDef(args...)
	if err != nil {
		return nil, err
	}
	if def == nil {
		return Bool{Null: true}, nil
	}

	_, ok := def.(*astdata.MapType)
	return Bool{Bool: ok}, nil
}

type mapKeyFn int

func (mapKeyFn) Execute(args ...Getter) (Getter, error) {
	def, err := getSingleDef(args...)
	if err != nil {
		return nil, err
	}
	if def == nil {
		return Definition{}, nil
	}

	mp, ok := def.(*astdata.MapType)
	if !ok {
		return Definition{}, nil
	}
	return Definition{Definition: mp.Key()}, nil
}

type mapValFn int

func (mapValFn) Execute(args ...Getter) (Getter, error) {
	def, err := getSingleDef(args...)
	if err != nil {
		return nil, err
	}
	if def == nil {
		return Definition{}, nil
	}

	mp, ok := def.(*astdata.MapType)
	if !ok {
		return Definition{}, nil
	}
	return Definition{Definition: mp.Val()}, nil
}

func registerMapFunc() {

	// Map functions
	RegisterFunction("is_map", isMapFn(0))
	RegisterFunction("map_key", mapKeyFn(0))
	RegisterFunction("map_val", mapValFn(0))
}

func init() {
	registerMapFunc()
}
