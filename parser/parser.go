package parser

import (
	"fmt"
	"pythia/ast"
	"pythia/lexer"
	"pythia/token"
)

const (
	_ int = iota
	LOWEST
	ASSIGN        // =
	LOGICAL_OR    // ||
	LOGICAL_AND   // &&
	BITWISE_OR    // |
	BITWISE_XOR   // ^
	BITWISE_AND   // &
	EQUALS        // ==
	LESSGREATER   // > or <
	BITWISE_SHIFT // >> or <<
	SUM           // +
	PRODUCT       // *
	PREFIX        // -X or !X
	CALL          // myFunction(X)
	INDEX         // array[index]
)

var precedences = map[token.TokenType]int{
	token.ASSIGN:             ASSIGN,
	token.PLUS_ASSIGN:        ASSIGN,
	token.MINUS_ASSIGN:       ASSIGN,
	token.ASTERISK_ASSIGN:    ASSIGN,
	token.SLASH_ASSIGN:       ASSIGN,
	token.PERCENT_ASSIGN:     ASSIGN,
	token.LOGICAL_AND:        LOGICAL_AND,
	token.LOGICAL_OR:         LOGICAL_OR,
	token.BINARY_OR:          BITWISE_OR,
	token.BINARY_XOR:         BITWISE_XOR,
	token.BINARY_AND:         BITWISE_AND,
	token.BINARY_LEFT_SHIFT:  BITWISE_SHIFT,
	token.BINARY_RIGHT_SHIFT: BITWISE_SHIFT,
	token.EQ:                 EQUALS,
	token.NOT_EQ:             EQUALS,
	token.LT:                 LESSGREATER,
	token.GT:                 LESSGREATER,
	token.LT_OR_EQ:           LESSGREATER,
	token.GT_OR_EQ:           LESSGREATER,
	token.PLUS:               SUM,
	token.MINUS:              SUM,
	token.SLASH:              PRODUCT,
	token.ASTERISK:           PRODUCT,
	token.PERCENT:            PRODUCT,
	token.LPAREN:             CALL,
	token.LBRACKET:           INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(expression ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.nextToken() // [nil, 0번째 토큰]
	p.nextToken() // [0번째 토큰, 1번째 토큰]

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)
	p.registerPrefix(token.NULL, p.parseNullLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.PERCENT, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LOGICAL_AND, p.parseInfixExpression)
	p.registerInfix(token.LOGICAL_OR, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LT_OR_EQ, p.parseInfixExpression)
	p.registerInfix(token.GT_OR_EQ, p.parseInfixExpression)
	p.registerInfix(token.BINARY_XOR, p.parseInfixExpression)
	p.registerInfix(token.BINARY_AND, p.parseInfixExpression)
	p.registerInfix(token.BINARY_OR, p.parseInfixExpression)
	p.registerInfix(token.BINARY_RIGHT_SHIFT, p.parseInfixExpression)
	p.registerInfix(token.BINARY_LEFT_SHIFT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(token.PLUS_ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(token.MINUS_ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(token.ASTERISK_ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(token.SLASH_ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(token.PERCENT_ASSIGN, p.parseAssignmentExpression)

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead, %s", t, p.peekToken.Type, p.l.GetErrorInfo())
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParserError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found, %s", t, p.l.GetErrorInfo())
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
