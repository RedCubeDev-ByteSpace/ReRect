package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type UnaryExpressionNode struct {
    ExpressionNode

    Operator lexer.Token
    Operand ExpressionNode
}

func NewUnaryExpressionNode(operand ExpressionNode, op lexer.Token) *UnaryExpressionNode {
    return &UnaryExpressionNode{
        Operator: op,
        Operand: operand,
    }
}

func (n *UnaryExpressionNode) Position() span.Span {
    return n.Operand.Position().SpanBetween(n.Operator.Position)
}

func (n *UnaryExpressionNode) Type() SyntaxNodeType {
    return NT_UnaryExpr
}
