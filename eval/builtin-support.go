// builtin-support.go - Helpers for registering "built-in" functions

package eval

import (
	"sync"

	"github.com/skx/gobasic/object"
)

// BuiltinSig is the signature of a builtin-function.
//
// Each built-in will receive an array of objects, and will return a
// single object back to the caller.
//
// In the case of an error then the object will be an error-object.
type BuiltinSig func(env Interpreter, args []object.Object) object.Object

// Builtins holds our state.
type Builtins struct {
	// lock holds a mutex to prevent corruption.
	lock sync.Mutex

	// argRegistry holds the number of arguments the given name requires.
	argRegistry map[string]int

	// fnRegistry holds a reference to the golang function which
	// implements the builtin.
	fnRegistry map[string]BuiltinSig
}

// NewBuiltins returns a new helper/holder for builtin functions.
func NewBuiltins() *Builtins {
	t := &Builtins{}
	t.argRegistry = make(map[string]int)
	t.fnRegistry = make(map[string]BuiltinSig)

	return t
}

// Register records a built-in function.
// The three arguments are:
//  NAME  - The thing that the BASIC program will call
//  nARGS - The number of arguments the built-in requires.
//          NOTE: Arguments are comma-separated in the BASIC program,
//          but commas are stripped out.
//  FT    - The function which provides the implementation.
func (b *Builtins) Register(name string, nArgs int, ft BuiltinSig) {
	b.lock.Lock()
	defer b.lock.Unlock()

	// Record the details.
	b.argRegistry[name] = nArgs
	b.fnRegistry[name] = ft
}

// Get the values associated with the given built-in.
func (b *Builtins) Get(name string) (int, BuiltinSig) {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.argRegistry[name], b.fnRegistry[name]
}
