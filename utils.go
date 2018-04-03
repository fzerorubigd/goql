package goql

import (
	"fmt"

	"github.com/fzerorubigd/goql/astdata"
)

func getSingleDef(args ...Getter) (astdata.Definition, error) {
	// using reflection to handle all types in one fn is possible, but not my style :)
	if err := required(1, 1, args...); err != nil {
		return nil, err
	}

	def := toDefinition(args[0].Get())
	return def, nil
}

func required(min, max int, args ...Getter) error {
	if len(args) < min || len(args) > max {
		return fmt.Errorf("argument count is wrong, got %d, needs %d to %d", len(args), min, max)
	}
	return nil
}
