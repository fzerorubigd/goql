package goql

import (
	"github.com/fzerorubigd/goql/astdata"
)

type embeder interface {
	Embeds() astdata.Embeds
}

type embedCountFn int

func (embedCountFn) Execute(args ...Getter) (Getter, error) {
	def, err := getSingleDef(args...)
	if err != nil {
		return nil, err
	}
	if def == nil {
		return Number{Null: true}, nil
	}

	st, ok := def.(embeder)
	if !ok {
		return Number{Null: true}, nil
	}
	return Number{Number: float64(len(st.Embeds()))}, nil
}

type embedDefFn int

func (embedDefFn) Execute(args ...Getter) (Getter, error) {
	if err := required(2, 2, args...); err != nil {
		return nil, err
	}

	def := toDefinition(args[0].Get())
	if def == nil {
		return Definition{}, nil
	}

	st, ok := def.(embeder)
	if !ok {
		return Definition{}, nil
	}
	em := st.Embeds()
	t := toNumber(args[1].Get())
	nm := int(t)
	if len(em) < nm || nm < 1 {
		return Definition{}, nil
	}
	return Definition{Definition: em[nm-1].Definition()}, nil
}

func registerCommonFunc() {
	RegisterFunction("embed_count", embedCountFn(0))
	RegisterFunction("embed_def", embedDefFn(0))
}

func init() {
	registerCommonFunc()
}
