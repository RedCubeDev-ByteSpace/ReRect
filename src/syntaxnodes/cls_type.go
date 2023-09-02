package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type TypeClauseNode struct {
    SyntaxNode

    TypeName lexer.Token
    SubTypes []*TypeClauseNode
}

func NewTypeClauseNode(typname lexer.Token, subtypes []*TypeClauseNode) TypeClauseNode {
    return TypeClauseNode{
        TypeName: typname,
        SubTypes: subtypes,
    }
}

func (n *TypeClauseNode) Position() span.Span {
    spn := n.TypeName.Position

    for _, v := range n.SubTypes {
        spn = spn.SpanBetween(v.Position())
    }

    return spn
}

func (n *TypeClauseNode) Type() SyntaxNodeType {
    return NT_TypeCls
}
