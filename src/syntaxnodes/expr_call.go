package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type CallExpressionNode struct {
    ExpressionNode

    Identifier lexer.Token
    Parameters []ExpressionNode
    CloseParam lexer.Token
}

func NewCallExpressionNode(id lexer.Token, param []ExpressionNode, cprm lexer.Token) CallExpressionNode {
    return CallExpressionNode{
        Identifier: id,
        Parameters: param,
        CloseParam: cprm,
    }
}

func (n *CallExpressionNode) Position() span.Span {
    return n.Identifier.Position.SpanBetween(n.CloseParam.Position)
}

func (n *CallExpressionNode) Type() SyntaxNodeType {
    return NT_CallExpr
}
