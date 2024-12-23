package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const(
	// 未知的词法字符
	ILLEGAL = "ILLEGAL"
	// 文件结尾
	EOF = "EOF"

	// 标识符
	IDENT = "IDENT"
	// 字面量
	INT = "INT"

	// 运算符
	ASSIGN = "="
	PLUS = "+"

	// 分隔符
	COMMA = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// 关键字
	FUNCTION = "FUNCTION"
	LET = "LET"
)