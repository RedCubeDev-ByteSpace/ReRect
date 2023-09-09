package syntaxnodes

import (
	"bytespace.network/rerect/span"
)

type ArrayIndexExpressionNode struct {
    ExpressionNode

    Expression ExpressionNode
    Index ExpressionNode
}

func NewArrayIndexExpressionNode(expr ExpressionNode, idx ExpressionNode) *ArrayIndexExpressionNode {
    return &ArrayIndexExpressionNode{
        Expression: expr,
        Index: idx,
    }
}

func (n *ArrayIndexExpressionNode) Position() span.Span {
    return n.Expression.Position().SpanBetween(n.Index.Position())
}

func (n *ArrayIndexExpressionNode) Type() SyntaxNodeType {
    return NT_ArrayIndexExpr
}
