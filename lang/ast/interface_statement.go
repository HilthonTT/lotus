package ast

import (
	"bytes"

	"github.com/hilthontt/lotus/token"
)

// InterfaceStatement: interface Drawable { fn draw(self) -> string }
type InterfaceStatement struct {
	Token   token.Token // the 'interface' token
	Name    *Identifier
	Methods []*InterfaceMethodSig
}

// InterfaceMethodSig describes one method signature inside an interface.
type InterfaceMethodSig struct {
	Name       string
	ParamTypes []*TypeAnnotation // parallel to params, may have nil entries
	ReturnType *TypeAnnotation   // may be nil
	ParamCount int
}

func (is *InterfaceStatement) statementNode() {}

func (is *InterfaceStatement) TokenLiteral() string {
	return is.Token.Literal
}

func (is *InterfaceStatement) String() string {
	var out bytes.Buffer
	out.WriteString("interface ")
	out.WriteString(is.Name.Value)
	out.WriteString(" { ")
	for _, m := range is.Methods {
		out.WriteString("fn " + m.Name + "(...)")
		if m.ReturnType != nil {
			out.WriteString(" -> " + m.ReturnType.Name)
		}
		out.WriteString("; ")
	}
	out.WriteString("}")
	return out.String()
}
