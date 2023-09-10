package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type FieldClauseNode struct {
    SyntaxNode

    FieldName lexer.Token
    FieldType *TypeClauseNode
}

func NewFieldClauseNode(prmname lexer.Token, typ *TypeClauseNode) *FieldClauseNode {
    return &FieldClauseNode{
        FieldName: prmname,
        FieldType: typ,
    }
}

func (n *FieldClauseNode) Position() span.Span {
    return n.FieldName.Position.SpanBetween(n.FieldType.Position())
}

func (n *FieldClauseNode) Type() SyntaxNodeType {
    return NT_FieldCls
}

