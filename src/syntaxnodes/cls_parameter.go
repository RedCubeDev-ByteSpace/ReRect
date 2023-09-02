package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type ParameterClauseNode struct {
    SyntaxNode

    ParameterName lexer.Token
    ParameterType TypeClauseNode
}

func NewParameterClauseNode(prmname lexer.Token, typ TypeClauseNode) ParameterClauseNode {
    return ParameterClauseNode{
        ParameterName: prmname,
        ParameterType: typ,
    }
}

func (n *ParameterClauseNode) Position() span.Span {
    return n.ParameterName.Position.SpanBetween(n.ParameterType.Position())
}

func (n *ParameterClauseNode) Type() SyntaxNodeType {
    return NT_ParameterCls
}

