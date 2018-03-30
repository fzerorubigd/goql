package goql

import (
	"github.com/fzerorubigd/goql/astdata"
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

func (genericFileName) Value(in interface{}) String {
	f := in.(filer)
	return String{String: f.File().FileName()}
}

type genericPackageName struct {
}

func (genericPackageName) Value(in interface{}) String {
	p := in.(packager)
	return String{String: p.Package().Name()}
}

type genericPackagePath struct {
}

func (genericPackagePath) Value(in interface{}) String {
	p := in.(packager)
	return String{String: p.Package().Path()}
}

type genericName struct {
}

func (genericName) Value(in interface{}) String {
	p := in.(namer)
	return String{String: p.Name()}
}

type genericIsExported struct {
}

func (genericIsExported) Value(in interface{}) Bool {
	p := in.(namer).Name()
	t := p[0] <= 'Z' && p[0] >= 'A'

	return Bool{Bool: t}
}

type genericDoc struct{}

func (genericDoc) Value(in interface{}) String {
	p := in.(docer).Docs()
	if len(p) == 0 {
		return String{Null: true}
	}
	return String{String: p.String()}
}
