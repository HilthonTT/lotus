package parser

import "github.com/hilthontt/lotus/token"

// synchronize skips tokens until it finds a safe statement boundary.
// Called after a parse error to allow collection of further errors.
func (p *Parser) synchronize() {
	p.nextToken()

	for !p.curTokenIs(token.EOF) {
		if p.prevToken.Type == token.SEMICOLON {
			return
		}
		switch p.curToken.Type {
		case token.FN,
			token.LET,
			token.MUT,
			token.FOR,
			token.WHILE,
			token.IF,
			token.RETURN,
			token.CLASS,
			token.ENUM,
			token.INTERFACE,
			token.IMPORT,
			token.EXPORT,
			token.DEFER,
			token.TRY,
			token.THROW:
			return
		}
		p.nextToken()
	}
}
