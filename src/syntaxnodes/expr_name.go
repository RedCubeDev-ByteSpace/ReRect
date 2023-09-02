package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type NameExpressionNode struct {
    ExpressionNode

    Identifier lexer.Token
}

func NewNameExpressionNode(id lexer.Token) *NameExpressionNode {
    return &NameExpressionNode{
        Identifier: id,
    }
}

func (n *NameExpressionNode) Position() span.Span {
    return n.Identifier.Position
}

func (n *NameExpressionNode) Type() SyntaxNodeType {
    return NT_NameExpr
}
