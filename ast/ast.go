package ast

type Node interface {
	TokenLiteral() string
	String() string // 디버깅용도
}
