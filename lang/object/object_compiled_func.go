package object

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hilthontt/lotus/code"
)

// CompiledFunction holds the bytecode of a compiled function.
type CompiledFunction struct {
	Instructions code.Instructions
	NumLocals    int
	NumParams    int
	Name         string
	IsVariadic   bool
	Lines        code.LineTable // maps bytecode offsets -> source lines
}

func (cf *CompiledFunction) Type() ObjectType {
	return "COMPILED_FUNCTION"
}

func (cf *CompiledFunction) Inspect() string {
	name := cf.Name
	if name == "" {
		name = "<anonymous>"
	}
	return fmt.Sprintf("<fn %s>", name)
}

func (cf *CompiledFunction) InvokeMethod(method string, env Environment, args ...Object) Object {
	if method == "methods" {
		static := []string{"methods"}
		dynamic := env.Names("function.")

		var names []string
		names = append(names, static...)
		for _, e := range dynamic {
			bits := strings.Split(e, ".")
			names = append(names, bits[1])
		}
		sort.Strings(names)

		result := make([]Object, len(names))
		for i, txt := range names {
			result[i] = &String{Value: txt}
		}
		return &Array{Elements: result}
	}
	return nil
}

func (cf *CompiledFunction) ToInterface() any {
	return "<COMPILED_FUNCTION>"
}
