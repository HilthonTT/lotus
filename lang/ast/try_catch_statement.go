package ast

import (
	"bytes"

	"github.com/hilthontt/lotus/token"
)

// try { ... } catch err { ... }
type TryCatchStatement struct {
	Token    token.Token // the 'try' token
	Try      *BlockStatement
	CatchVar *Identifier // may be nil if catch has no binding: catch { ... }
	Catch    *BlockStatement
}

func (tc *TryCatchStatement) statementNode() {}

func (tc *TryCatchStatement) TokenLiteral() string {
	return tc.Token.Literal
}

func (tc *TryCatchStatement) String() string {
	var out bytes.Buffer
	out.WriteString("try { ")
	out.WriteString(tc.Try.String())
	out.WriteString(" } catch ")
	if tc.CatchVar != nil {
		out.WriteString(tc.CatchVar.Value + " ")
	}
	out.WriteString("{ ")
	out.WriteString(tc.Catch.String())
	out.WriteString(" }")
	return out.String()
}
