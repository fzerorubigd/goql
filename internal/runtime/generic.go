package runtime

import "github.com/fzerorubigd/goql/astdata"

type filer interface {
	File() *astdata.File
}

type packager interface {
	Package() *astdata.Package
}

type namer interface {
	Name() string
}

type genericFileName struct {
}

func (genericFileName) Value(in interface{}) string {
	f := in.(filer)
	return f.File().FileName()
}

type genericPackageName struct {
}

func (genericPackageName) Value(in interface{}) string {
	p := in.(packager)
	return p.Package().Name()
}

type genericPackagePath struct {
}

func (genericPackagePath) Value(in interface{}) string {
	p := in.(packager)
	return p.Package().Path()
}

type genericName struct {
}

func (genericName) Value(in interface{}) string {
	p := in.(namer)
	return p.Name()
}
