package ast

import (
	"bytes"

	"github.com/hilthontt/lotus/token"
)

type IndexAssignStatement struct {
	Token token.Token
	Left  Expression // the array/map expression
	Index Expression // the index expression
	Value Expression
}

func (s *IndexAssignStatement) statementNode() {}

// TokenLiteral returns the literal token.
func (s *IndexAssignStatement) TokenLiteral() string {
	return s.Token.Literal
}

func (s *IndexAssignStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.Left.String())
	out.WriteString("[")
	out.WriteString(s.Index.String())
	out.WriteString("] = ")
	out.WriteString(s.Value.String())
	out.WriteString(";")
	return out.String()
}
