package structures

import (
	"fmt"
	"sync"
)

// Function is the functions in the system
type Function interface {
	Execute(...Valuer) (Valuer, error)
}

var (
	functions = make(map[string]Function)
	fnLock    = &sync.RWMutex{}
)

// RegisterFunction is entry point for registering a function into system
func RegisterFunction(name string, fn Function) {
	fnLock.Lock()
	defer fnLock.Unlock()

	if _, ok := functions[name]; ok {
		panic(fmt.Sprintf("function with name '%s' is already registered", name))
	}

	functions[name] = fn
}

// HasFunction return if the function is available
func HasFunction(name string) bool {
	fnLock.RLock()
	defer fnLock.RUnlock()

	_, ok := functions[name]
	return ok
}

// ExecuteFunction is a helper to execute function by its name
func ExecuteFunction(name string, value ...Valuer) (Valuer, error) {
	fnLock.RLock()
	defer fnLock.RUnlock()

	fn, ok := functions[name]
	if !ok {
		return nil, fmt.Errorf("function with name '%s' is not registered", name)
	}

	return fn.Execute(value...)
}
