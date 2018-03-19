package structures

import "fmt"

// Function is the functions in the system
type Function interface {
	Execute(...Valuer) (Valuer, error)
}

var (
	functions = make(map[string]Function)
)

// RegisterFunc is entry point for registering a function into system
func RegisterFunc(name string, fn Function) {
	lock.Lock()
	defer lock.Unlock()

	if _, ok := functions[name]; ok {
		panic(fmt.Sprintf("function with name '%s' is already registered", name))
	}

	functions[name] = fn
}

// ExecuteFunction is a helper to execute function by its name
func ExecuteFunction(name string, value ...Valuer) (Valuer, error) {
	lock.Lock()
	defer lock.Unlock()
	fn, ok := functions[name]
	if !ok {
		return nil, fmt.Errorf("function with name '%s' is not registered", name)
	}

	return fn.Execute(value...)
}
