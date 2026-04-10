package ast

import "github.com/hilthontt/lotus/token"

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (e *IntegerLiteral) expressionNode() {}

// TokenLiteral returns the literal token.
func (e *IntegerLiteral) TokenLiteral() string {
	return e.Token.Literal
}

func (e *IntegerLiteral) String() string {
	return e.Token.Literal
}
