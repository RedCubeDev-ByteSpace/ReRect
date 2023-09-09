package syntaxnodes

import (
	"bytespace.network/rerect/span"
)

type AssignmentExpressionNode struct {
    ExpressionNode

    Expression ExpressionNode
    Value ExpressionNode
}

func NewAssignmentExpressionNode(expr ExpressionNode, val ExpressionNode) *AssignmentExpressionNode {
    return &AssignmentExpressionNode{
        Expression: expr,
        Value: val,
    }
}

func (n *AssignmentExpressionNode) Position() span.Span {
    return n.Expression.Position().SpanBetween(n.Expression.Position())
}

func (n *AssignmentExpressionNode) Type() SyntaxNodeType {
    return NT_AssignmentExpr
}
