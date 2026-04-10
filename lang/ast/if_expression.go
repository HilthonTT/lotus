package ast

import (
	"bytes"

	"github.com/hilthontt/lotus/token"
)

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (e *IfExpression) expressionNode() {}

// TokenLiteral returns the literal token.
func (e *IfExpression) TokenLiteral() string {
	return e.Token.Literal
}

func (e *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if ")
	out.WriteString(e.Condition.String())
	out.WriteString(" { ")
	out.WriteString(e.Consequence.String())
	out.WriteString(" }")
	if e.Alternative != nil {
		out.WriteString(" else { ")
		out.WriteString(e.Alternative.String())
		out.WriteString(" }")
	}
	return out.String()
}
