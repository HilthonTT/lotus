package object

import "fmt"

// Module holds the exported values of a compiled .lotus fule.
type Module struct {
	Path    string
	Exports map[string]Object
}

func (m *Module) Type() ObjectType {
	return MODULE_OBJ
}

func (m *Module) Inspect() string {
	return fmt.Sprintf("<module %q>", m.Path)
}
