package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type LiteralExpressionNode struct {
    ExpressionNode

    Literal lexer.Token
}

func NewLiteralExpressionNode(lit lexer.Token) *LiteralExpressionNode {
    return &LiteralExpressionNode{
        Literal: lit,
    }
}

func (n *LiteralExpressionNode) Position() span.Span {
    return n.Literal.Position
}

func (n *LiteralExpressionNode) Type() SyntaxNodeType {
    return NT_LiteralExpr
}
