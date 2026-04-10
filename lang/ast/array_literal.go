package ast

import (
	"bytes"

	"github.com/hilthontt/lotus/token"
)

type ArrayLiteral struct {
	Token    token.Token // the [ token
	Elements []Expression
}

func (e *ArrayLiteral) expressionNode() {}

// TokenLiteral returns the literal token.
func (e *ArrayLiteral) TokenLiteral() string {
	return e.Token.Literal
}

func (e *ArrayLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("[")
	for i, el := range e.Elements {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(el.String())
	}
	out.WriteString("]")
	return out.String()
}
