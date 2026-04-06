package ast

import (
	"bytes"

	"github.com/hilthontt/lotus/token"
)

type ContinueStatement struct {
	Token token.Token
}

func (s *ContinueStatement) statementNode() {}

// TokenLiteral returns the literal token.
func (s *ContinueStatement) TokenLiteral() string {
	return s.Token.Literal
}

func (s *ContinueStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	out.WriteString(";")
	return out.String()
}
