package ast

import "github.com/hilthontt/lotus/token"

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (e *BooleanLiteral) expressionNode() {}

// TokenLiteral returns the literal token.
func (e *BooleanLiteral) TokenLiteral() string {
	return e.Token.Literal
}

func (e *BooleanLiteral) String() string {
	return e.Token.Literal
}
