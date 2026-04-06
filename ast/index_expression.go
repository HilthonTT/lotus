package ast

import (
	"bytes"

	"github.com/hilthontt/lotus/token"
)

type IndexExpression struct {
	Token token.Token // the [ token
	Left  Expression
	Index Expression
}

func (e *IndexExpression) expressionNode() {}

// TokenLiteral returns the literal token.
func (e *IndexExpression) TokenLiteral() string {
	return e.Token.Literal
}

// String returns this object as a string.
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}
