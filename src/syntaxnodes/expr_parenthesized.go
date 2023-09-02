package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type ParenthesizedExpressionNode struct {
    ExpressionNode

    StartParenth lexer.Token
    Expression ExpressionNode
    EndParenth lexer.Token
}

func NewParenthesizedExpressionNode(strt lexer.Token, expr ExpressionNode, end lexer.Token) *ParenthesizedExpressionNode {
    return &ParenthesizedExpressionNode{
        StartParenth: strt,
        Expression: expr,
        EndParenth: end,
    }
}

func (n *ParenthesizedExpressionNode) Position() span.Span {
    return n.StartParenth.Position.SpanBetween(n.EndParenth.Position)
}

func (n *ParenthesizedExpressionNode) Type() SyntaxNodeType {
    return NT_ParenthesizedExpr
}
