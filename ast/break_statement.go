package ast

import (
	"bytes"

	"github.com/hilthontt/lotus/token"
)

type BreakStatement struct {
	Token token.Token
}

func (s *BreakStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	out.WriteString(";")
	return out.String()
}

func (s *BreakStatement) statementNode() {}

// TokenLiteral returns the literal token.
func (s *BreakStatement) TokenLiteral() string {
	return s.Token.Literal
}
