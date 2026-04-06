package ast

import (
	"bytes"

	"github.com/hilthontt/lotus/token"
)

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (e *PrefixExpression) expressionNode() {}

// TokenLiteral returns the literal token.
func (e *PrefixExpression) TokenLiteral() string {
	return e.Token.Literal
}

func (e *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(e.Operator)
	out.WriteString(e.Right.String())
	out.WriteString(")")
	return out.String()
}
