package ast

import (
	"bytes"

	"github.com/hilthontt/lotus/token"
)

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (e *InfixExpression) expressionNode() {}

// TokenLiteral returns the literal token.
func (e *InfixExpression) TokenLiteral() string {
	return e.Token.Literal
}

func (e *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(e.Left.String())
	out.WriteString(" " + e.Operator + " ")
	out.WriteString(e.Right.String())
	out.WriteString(")")
	return out.String()
}
