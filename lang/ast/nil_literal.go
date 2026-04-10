package ast

import "github.com/hilthontt/lotus/token"

type NilLiteral struct {
	Token token.Token
}

func (e *NilLiteral) expressionNode() {}

// TokenLiteral returns the literal token.
func (e *NilLiteral) TokenLiteral() string {
	return e.Token.Literal
}

func (e *NilLiteral) String() string {
	return "nil"
}
