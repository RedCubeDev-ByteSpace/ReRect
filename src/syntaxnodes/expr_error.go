package syntaxnodes

import (
	"bytespace.network/rerect/span"
)

type ErrorExpressionNode struct {
    ExpressionNode
    
    ErrPosition span.Span
}

func NewErrorExpressionNode(pos span.Span) *ErrorExpressionNode {
    return &ErrorExpressionNode{
        ErrPosition: pos,
    }
}

func (n *ErrorExpressionNode) Position() span.Span {
    return n.ErrPosition
}

func (n *ErrorExpressionNode) Type() SyntaxNodeType {
    return NT_ErrorExpr
}
