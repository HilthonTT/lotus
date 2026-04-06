package ast

import "github.com/hilthontt/lotus/token"

// Identifier holds a single identifier.
type Identifier struct {
	// Token is the literal token
	Token token.Token

	// Value is the name of the identifier
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral returns the literal token.
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// String returns this object as a string.
func (i *Identifier) String() string {
	return i.Value
}
