package ast

import (
	"bytes"

	"github.com/hilthontt/lotus/token"
)

type MapLiteral struct {
	Token token.Token // the { token
	Pairs map[Expression]Expression
	Keys  []Expression // preserve order
}

func (e *MapLiteral) expressionNode() {}

// TokenLiteral returns the literal token.
func (e *MapLiteral) TokenLiteral() string {
	return e.Token.Literal
}

func (e *MapLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("{")
	for i, key := range e.Keys {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(key.String())
		out.WriteString(": ")
		out.WriteString(e.Pairs[key].String())
	}
	out.WriteString("}")
	return out.String()
}
