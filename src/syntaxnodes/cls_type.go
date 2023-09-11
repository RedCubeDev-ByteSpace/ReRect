package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type TypeClauseNode struct {
    SyntaxNode

    PackageName lexer.Token
    HasPackageName bool

    TypeName lexer.Token
    SubTypes []*TypeClauseNode
}

func NewTypeClauseNode(packname lexer.Token, haspack bool, typname lexer.Token, subtypes []*TypeClauseNode) *TypeClauseNode {
    return &TypeClauseNode{
        PackageName: packname,
        HasPackageName: haspack,
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
