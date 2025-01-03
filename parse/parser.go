package parser

import(
	"Monkey/ast"
	"Monkey/lexer"
	"Monkey/token"
)

type Parse struct{
	l *lexer.lexer
	// curToken 指向当前的token
	curToken token.Token
	// peekToken 指向下一个token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parse{
	p := &Parse{l: l}
	// 读取两个词法单元来设置当前token和下一个token
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parse) nextToken(){
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parse) ParseProgram() *ast.Program {
	return nil
}
