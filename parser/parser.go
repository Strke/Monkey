package parser

import(
	"Monkey/ast"
	"Monkey/lexer"
	"Monkey/token"
	"fmt"
)

type Parse struct{
	l *lexer.Lexer
	// curToken 指向当前的token
	curToken token.Token
	// peekToken 指向下一个token
	peekToken token.Token

	errors []string
}

func New(l *lexer.Lexer) *Parse{
	p := &Parse{
		l: l,
		errors: []string{},
	}
	// 读取两个词法单元来设置当前token和下一个token
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parse) nextToken(){
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parse) Errors() []string{
	return p.errors
}

func (p *Parse) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parse) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF{
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parse) parseStatement() ast.Statement{
	switch p.curToken.Type{
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parse) parseLetStatement() *ast.LetStatement{
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT){
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN){
		return nil
	}
	// Todo: 跳过对表达式的解析，直到遇见分号
	for !p.curTokenIs(token.SEMICOLON){
		p.nextToken()
	}
	return stmt
}
func (p *Parse) curTokenIs(t token.TokenType) bool{
	return p.curToken.Type == t
}
func (p *Parse) peekTokenIs(t token.TokenType) bool{
	return p.peekToken.Type == t
}
func (p *Parse) expectPeek(t token.TokenType) bool{
	if p.peekTokenIs(t){
		p.nextToken()
		return true
	}else{
		p.peekError(t)
		return false
	}
}