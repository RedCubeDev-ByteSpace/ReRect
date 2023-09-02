package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type BinaryExpressionNode struct {
    ExpressionNode

    Left ExpressionNode
    Right ExpressionNode
    Operator lexer.Token
}

func NewBinaryExpressionNode(left ExpressionNode, right ExpressionNode, op lexer.Token) *BinaryExpressionNode {
    return &BinaryExpressionNode{
        Left: left,
        Right: right,
        Operator: op,
    }
}

func (n *BinaryExpressionNode) Position() span.Span {
    return n.Left.Position().SpanBetween(n.Right.Position())
}

func (n *BinaryExpressionNode) Type() SyntaxNodeType {
    return NT_BinaryExpr
}
