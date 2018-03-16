package runtime

import (
	"github.com/fzerorubigd/goql/astdata"
	"github.com/fzerorubigd/goql/structures"
)

type filer interface {
	File() *astdata.File
}

type packager interface {
	Package() *astdata.Package
}

type namer interface {
	Name() string
}

type docer interface {
	Docs() astdata.Docs
}

type genericFileName struct {
}

func (genericFileName) Value(in interface{}) structures.String {
	f := in.(filer)
	return structures.String{String: f.File().FileName()}
}

type genericPackageName struct {
}

func (genericPackageName) Value(in interface{}) structures.String {
	p := in.(packager)
	return structures.String{String: p.Package().Name()}
}

type genericPackagePath struct {
}

func (genericPackagePath) Value(in interface{}) structures.String {
	p := in.(packager)
	return structures.String{String: p.Package().Path()}
}

type genericName struct {
}

func (genericName) Value(in interface{}) structures.String {
	p := in.(namer)
	return structures.String{String: p.Name()}
}

type genericIsExported struct {
}

func (genericIsExported) Value(in interface{}) structures.Bool {
	p := in.(namer).Name()
	t := p[0] <= 'Z' && p[0] >= 'A'

	return structures.Bool{Bool: t}
}

type genericDoc struct{}

func (genericDoc) Value(in interface{}) structures.String {
	p := in.(docer).Docs()
	if len(p) == 0 {
		return structures.String{Null: true}
	}
	return structures.String{String: p.String()}
}
