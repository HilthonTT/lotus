package object

import (
	"fmt"
	"strings"
)

// InterfaceMethodSpec describes one method required by an interface.
type InterfaceMethodSpec struct {
	Name       string
	ParamCount int    // number of params excluding self
	ReturnType string // "" means any / unspecified
}

// Interface is the runtime object for an interface definition.
type Interface struct {
	Name    string
	Methods []InterfaceMethodSpec
}

func (i *Interface) Type() ObjectType {
	return INTERFACE_OBJ
}

func (i *Interface) Inspect() string {
	names := make([]string, len(i.Methods))
	for j, m := range i.Methods {
		names[j] = m.Name
	}
	return fmt.Sprintf("<interface %s [%s]>", i.Name, strings.Join(names, ", "))
}

// Implements checks whether a Lotus Instance (or any object with methods)
// satisfies this interface structurally.
func (iface *Interface) Implements(obj Object) bool {
	inst, ok := obj.(*Instance)
	if !ok {
		return false
	}
	for _, spec := range iface.Methods {
		method, found := inst.Class.LookupMethod(spec.Name)
		if !found {
			return false
		}
		// Check param count (NumParams includes self)
		if method.Fn.NumParams != spec.ParamCount+1 {
			return false
		}
	}
	return true
}
