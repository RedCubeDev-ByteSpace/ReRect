package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type TraitClauseNode struct {
    SyntaxNode

    Package lexer.Token
    HasPackage bool

    TraitName lexer.Token
}

func NewTraitClauseNode(pck lexer.Token, haspack bool, name lexer.Token) *TraitClauseNode {
    return &TraitClauseNode{
        Package: pck,
        HasPackage: haspack,
        TraitName: name,
    }
}

func (n *TraitClauseNode) Position() span.Span {
    if n.HasPackage {
        return n.Package.Position.SpanBetween(n.TraitName.Position)
    } else {
        return n.TraitName.Position
    }
}

func (n *TraitClauseNode) Type() SyntaxNodeType {
    return NT_TraitCls
}
