package goql

import (
	"fmt"
	"sync"
)

// Function is the functions in the system, it should check for arguments count and return error on
// wrong arguments count, but for type, it should try to cast
// TODO : check for function arguments on prepare
type Function interface {
	// Execute is called on each row
	Execute(...Valuer) (Valuer, error)
}

var (
	functions = make(map[string]Function)
	fnLock    = &sync.RWMutex{}
)

// RegisterFunction is entry point for registering a function into system, the name must be unique
func RegisterFunction(name string, fn Function) {
	fnLock.Lock()
	defer fnLock.Unlock()

	if _, ok := functions[name]; ok {
		panic(fmt.Sprintf("function with name '%s' is already registered", name))
	}

	functions[name] = fn
}

// hasFunction return if the function is available
func hasFunction(name string) bool {
	fnLock.RLock()
	defer fnLock.RUnlock()

	_, ok := functions[name]
	return ok
}

// executeFunction is a helper to execute function by its name
func executeFunction(name string, value ...Valuer) (Valuer, error) {
	fnLock.RLock()
	defer fnLock.RUnlock()

	fn, ok := functions[name]
	if !ok {
		return nil, fmt.Errorf("function with name '%s' is not registered", name)
	}

	return fn.Execute(value...)
}
