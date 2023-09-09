package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type AccessExpressionNode struct {
    ExpressionNode

    Expression ExpressionNode
    Identifier lexer.Token

    Arguments []ExpressionNode
    Closing lexer.Token
    IsCall bool
}

func NewAccessExpressionNode(expr ExpressionNode, id lexer.Token, args []ExpressionNode, cls lexer.Token, iscall bool) *AccessExpressionNode {
    return &AccessExpressionNode{
        Expression: expr,
        Identifier: id,
        Arguments: args,
        Closing: cls,
        IsCall: iscall,
    }
}

func (n *AccessExpressionNode) Position() span.Span {
    if !n.IsCall {
        return n.Expression.Position().SpanBetween(n.Identifier.Position)
    } else {
        return n.Expression.Position().SpanBetween(n.Closing.Position)
    }
}

func (n *AccessExpressionNode) Type() SyntaxNodeType {
    return NT_AccessExpr
}
