package parser

import(
	"Monkey/ast"
	"Monkey/lexer"
	"Monkey/token"
	"fmt"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS  //==
	LESSGREATER //< or >
	SUM //+
	PRODUCT //*
	PREFIX //-X or !X
	CALL //myFunction(X)
)

type (
	prefixParserFn func() ast.Expression
	infixParserFn func(ast.Expression) ast.Expression
)

type Parser struct{
	l *lexer.Lexer
	errors []string
	// curToken 指向当前的token
	curToken token.Token
	// peekToken 指向下一个token
	peekToken token.Token

	prefixParserFns map[token.TokenType]prefixParserFn
	infixParserFns map[token.TokenType]infixParserFn
}
func (p *Parser) nextToken(){
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []string{
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) ParseProgram() *ast.Program {
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
func (p *Parser) parseStatement() ast.Statement{
	switch p.curToken.Type{
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement{
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

func (p *Parser) parseReturnStatement() *ast.ReturnStatement{
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	// Todo: 跳过对表达式的解析，直到遇见分号
	for !p.curTokenIs(token.SEMICOLON){
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement{
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	// Todo: 跳过对表达式的解析，直到遇见分号
	for !p.curTokenIs(token.SEMICOLON){
		p.nextToken()
	}
	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool{
	return p.curToken.Type == t
}
func (p *Parser) peekTokenIs(t token.TokenType) bool{
	return p.peekToken.Type == t
}
func (p *Parser) expectPeek(t token.TokenType) bool{
	if p.peekTokenIs(t){
		p.nextToken()
		return true
	}else{
		p.peekError(t)
		return false
	}
}
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParserFn){
	p.prefixParserFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParserFn){
	p.infixParserFns[tokenType] = fn
}
func New(l *lexer.Lexer) *Parser{
	p := &Parser{
		l: l,
		errors: []string{},
	}
	p.prefixParserFns = make(map[token.TokenType]prefixParserFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	// 读取两个词法单元来设置当前token和下一个token
	p.nextToken()
	p.nextToken()
	return p
}
func (p *Parser) parseIdentifier() ast.Expression{
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

// 前缀为空的情况下报的错误
func (p *Parser) noPrefixParseFnError(t token.TokenType){
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}


func (p *Parser) parseExpression(precedence int) ast.Expression{
	prefix := p.prefixParserFns[p.curToken.Type]
	if prefix == nil{
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	return leftExp
}

func (p *Parser) parseIntegerLiteral() ast.Expression{
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil{
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression{
	expression := &ast.PrefixExpression{
		Token: p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}