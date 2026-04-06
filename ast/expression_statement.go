package ast

import "github.com/hilthontt/lotus/token"

// ExpressionStatement is an expression.
type ExpressionStatement struct {
	// Token is the literal token.
	Token token.Token

	// Expression holds the expression.
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// TokenLiteral returns the literal token.
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

// String returns this object as a string.
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
