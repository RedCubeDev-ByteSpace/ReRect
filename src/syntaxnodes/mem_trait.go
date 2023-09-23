package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type TraitNode struct {
    MemberNode

    TraitKw lexer.Token
    TraitName lexer.Token

    Fields []*FieldClauseNode
    Methods []*FunctionNode

    Closing lexer.Token
}

func NewTraitNode(kw lexer.Token, name lexer.Token, fields []*FieldClauseNode, meth []*FunctionNode, cls lexer.Token) *TraitNode {
    return &TraitNode{
        TraitKw: kw,
        TraitName: name,
        Fields: fields,
        Methods: meth,
        Closing: cls,
    }
}

func (n *TraitNode) Position() span.Span {
    return n.TraitKw.Position.SpanBetween(n.Closing.Position)
}

func (n *TraitNode) Type() SyntaxNodeType {
    return NT_Trait
}
