package ast

import (
	"bytes"

	"github.com/hilthontt/lotus/token"
)

type CallExpression struct {
	Token     token.Token // the ( token
	Function  Expression
	Arguments []Expression
}

func (e *CallExpression) expressionNode() {}

// TokenLiteral returns the literal token.
func (e *CallExpression) TokenLiteral() string {
	return e.Token.Literal
}
func (e *CallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(e.Function.String())
	out.WriteString("(")
	for i, a := range e.Arguments {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(a.String())
	}
	out.WriteString(")")
	return out.String()
}
