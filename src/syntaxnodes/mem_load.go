package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type LoadNode struct {
    MemberNode

    LoadKw lexer.Token
    Library lexer.Token
    IncludeKw lexer.Token

    Included bool
}

func NewLoadNode(loadkw lexer.Token, lib lexer.Token, includekw lexer.Token, included bool) *LoadNode {
    return &LoadNode{
        LoadKw: loadkw,
        Library: lib,
        IncludeKw: includekw,
        Included: included,
    }
}

func (n *LoadNode) Position() span.Span {
    if n.Included {
        return n.LoadKw.Position.SpanBetween(n.IncludeKw.Position)
    } else {
        return n.LoadKw.Position.SpanBetween(n.Library.Position)
    }
}

func (n *LoadNode) Type() SyntaxNodeType {
    return NT_Load
}
