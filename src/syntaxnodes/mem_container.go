package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type ContainerNode struct {
    MemberNode

    ContainerKw lexer.Token
    ContainerName lexer.Token

    Fields []*FieldClauseNode
    Methods []*FunctionNode

    Closing lexer.Token
}

func NewContainerNode(kw lexer.Token, name lexer.Token, fields []*FieldClauseNode, meth []*FunctionNode, cls lexer.Token) *ContainerNode {
    return &ContainerNode{
        ContainerKw: kw,
        ContainerName: name,
        Fields: fields,
        Methods: meth,
        Closing: cls,
    }
}

func (n *ContainerNode) Position() span.Span {
    return n.ContainerKw.Position.SpanBetween(n.Closing.Position)
}

func (n *ContainerNode) Type() SyntaxNodeType {
    return NT_Container
}
