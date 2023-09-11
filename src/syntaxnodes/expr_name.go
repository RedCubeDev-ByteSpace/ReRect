package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type NameExpressionNode struct {
    ExpressionNode

    PackageName lexer.Token
    HasPackage bool

    Identifier lexer.Token
}

func NewNameExpressionNode(pack lexer.Token, haspack bool, id lexer.Token) *NameExpressionNode {
    return &NameExpressionNode{
        PackageName: pack,
        HasPackage: haspack,
        Identifier: id,
    }
}

func (n *NameExpressionNode) Position() span.Span {
    return n.Identifier.Position
}

func (n *NameExpressionNode) Type() SyntaxNodeType {
    return NT_NameExpr
}
