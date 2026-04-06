package ast

import "github.com/hilthontt/lotus/token"

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (e *FloatLiteral) expressionNode() {}

// TokenLiteral returns the literal token.
func (e *FloatLiteral) TokenLiteral() string {
	return e.Token.Literal
}

func (e *FloatLiteral) String() string {
	return e.Token.Literal
}
