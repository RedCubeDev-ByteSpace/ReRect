package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type CallExpressionNode struct {
    ExpressionNode

    Package lexer.Token
    HasPackage bool

    Identifier lexer.Token
    Parameters []ExpressionNode
    CloseParam lexer.Token
}

func NewCallExpressionNode(id lexer.Token, pack lexer.Token, hasPack bool, param []ExpressionNode, cprm lexer.Token) *CallExpressionNode {
    return &CallExpressionNode{
        Package: pack,
        HasPackage: hasPack,

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
