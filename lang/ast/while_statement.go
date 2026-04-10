package ast

import (
	"bytes"

	"github.com/hilthontt/lotus/token"
)

type WhileStatement struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (s *WhileStatement) statementNode() {}

// TokenLiteral returns the literal token.
func (s *WhileStatement) TokenLiteral() string {
	return s.Token.Literal
}

func (s *WhileStatement) String() string {
	var out bytes.Buffer
	out.WriteString("while ")
	out.WriteString(s.Condition.String())
	out.WriteString(" { ")
	out.WriteString(s.Body.String())
	out.WriteString(" }")
	return out.String()
}
