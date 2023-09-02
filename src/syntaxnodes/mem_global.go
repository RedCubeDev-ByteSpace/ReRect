package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type GlobalNode struct {
    MemberNode

    VarKw lexer.Token
    GlobalName lexer.Token
    VarType *TypeClauseNode
}

func NewGlobalNode(varkw lexer.Token, glbname lexer.Token, typ *TypeClauseNode) *GlobalNode {
    return &GlobalNode{
        VarKw: varkw,
        GlobalName: glbname,
        VarType: typ,
    }
}

func (n *GlobalNode) Position() span.Span {
    return n.VarKw.Position.SpanBetween(n.VarType.Position())
}

func (n *GlobalNode) Type() SyntaxNodeType {
    return NT_Global
}
