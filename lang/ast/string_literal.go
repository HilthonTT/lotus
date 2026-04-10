package ast

import "github.com/hilthontt/lotus/token"

type StringLiteral struct {
	Token token.Token
	Value string
}

func (e *StringLiteral) expressionNode() {}

// TokenLiteral returns the literal token.
func (e *StringLiteral) TokenLiteral() string {
	return e.Token.Literal
}

func (e *StringLiteral) String() string {
	return e.Value
}
